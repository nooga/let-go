//go:build unix

package main

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func testPool(t *testing.T) string {
	t.Helper()
	t.Setenv("BENCH_BATON_HOME", t.TempDir())
	return "testpool"
}

// Two shared leases overlap: the shared lane is concurrent up to maxShared.
func TestSharedLeasesOverlap(t *testing.T) {
	pool := testPool(t)
	a, err := acquire("a", modeShared, pool, "t", 5*time.Second, 4, false)
	if err != nil {
		t.Fatal(err)
	}
	defer a.release()
	b, err := acquire("b", modeShared, pool, "t", 500*time.Millisecond, 4, false)
	if err != nil {
		t.Fatalf("second shared lease did not overlap the first: %v", err)
	}
	b.release()
	st := status(pool)
	if len(st.Holders) != 1 || st.Holders[0].Mode != modeShared {
		t.Fatalf("expected one shared holder after releasing b, got %+v", st.Holders)
	}
}

// The shared lane is capped at maxShared slots.
func TestSharedLaneIsCapped(t *testing.T) {
	pool := testPool(t)
	a, err := acquire("a", modeShared, pool, "t", time.Second, 1, false)
	if err != nil {
		t.Fatal(err)
	}
	defer a.release()
	if _, err := acquire("b", modeShared, pool, "t", 300*time.Millisecond, 1, false); err == nil {
		t.Fatal("second shared lease was granted with maxShared=1")
	}
}

// An exclusive lease waits for every shared holder to finish.
func TestExclusiveWaitsForShared(t *testing.T) {
	pool := testPool(t)
	sh, err := acquire("build", modeShared, pool, "t", time.Second, 4, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := acquire("bench", modeExclusive, pool, "t", 300*time.Millisecond, 4, false); err == nil {
		t.Fatal("exclusive lease was granted while a shared lease was held")
	}
	sh.release()
	ex, err := acquire("bench", modeExclusive, pool, "t", time.Second, 4, false)
	if err != nil {
		t.Fatalf("exclusive lease not granted after shared release: %v", err)
	}
	ex.release()
}

// A waiting exclusive closes the gate: shared work arriving behind it is
// deferred until the exclusive has run, so builds cannot starve a benchmark.
func TestWaitingExclusiveClosesGate(t *testing.T) {
	pool := testPool(t)
	first, err := acquire("build-1", modeShared, pool, "t", time.Second, 4, false)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	granted := make(chan time.Time, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ex, err := acquire("bench", modeExclusive, pool, "t", 5*time.Second, 4, false)
		if err != nil {
			t.Error(err)
			return
		}
		granted <- time.Now()
		time.Sleep(200 * time.Millisecond)
		ex.release()
	}()
	// Give the exclusive time to queue and close the gate.
	time.Sleep(150 * time.Millisecond)
	if _, err := acquire("build-2", modeShared, pool, "t", 200*time.Millisecond, 4, false); err == nil {
		t.Fatal("shared lease was granted while an exclusive was waiting at the gate")
	}
	first.release()
	wg.Wait()
	select {
	case <-granted:
	default:
		t.Fatal("exclusive never ran after the shared holder released")
	}
	late, err := acquire("build-3", modeShared, pool, "t", time.Second, 4, false)
	if err != nil {
		t.Fatalf("shared lease not granted after the exclusive finished: %v", err)
	}
	late.release()
}

// Bookkeeping for a holder whose process is gone is dropped by reap; the
// flock itself was already released by the OS.
func TestReapDropsDeadHolders(t *testing.T) {
	pool := testPool(t)
	if err := mutateState(pool, func(st *state) {
		st.Holders = append(st.Holders, token{ID: "dead", Owner: "x", Mode: modeShared, PID: 1 << 30})
	}); err != nil {
		t.Fatal(err)
	}
	dropped, err := reap(pool)
	if err != nil {
		t.Fatal(err)
	}
	if len(dropped) != 1 || dropped[0].ID != "dead" {
		t.Fatalf("reap dropped %+v, want the dead holder", dropped)
	}
	if st := status(pool); len(st.Holders) != 0 {
		t.Fatalf("dead holder still recorded: %+v", st.Holders)
	}
}

// run executes under the lease, writes the full output to a log, appends a
// ledger row, and reports the command's exit code.
func TestRunLogsAndLedgers(t *testing.T) {
	pool := testPool(t)
	res, err := run(runOptions{
		argv: []string{"sh", "-c", "echo hello; exit 3"}, cwd: t.TempDir(),
		intent: "probe", mode: modeShared, pool: pool, owner: "t",
		leaseTimeout: time.Second, maxShared: 4, tailLines: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.RC != 3 {
		t.Fatalf("rc = %d, want 3", res.RC)
	}
	if res.Tail != "hello" {
		t.Fatalf("tail = %q, want %q", res.Tail, "hello")
	}
	if _, err := os.Stat(res.Log); err != nil {
		t.Fatalf("log file missing: %v", err)
	}
	rows, err := ledger(pool, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].RC != 3 || rows[0].Owner != "t" || rows[0].Intent != "probe" {
		t.Fatalf("ledger rows = %+v", rows)
	}
	if filepath.Dir(res.Log) != filepath.Join(poolDir(pool), "logs") {
		t.Fatalf("log written outside the pool's logs dir: %s", res.Log)
	}
}
