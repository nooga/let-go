//go:build linux || darwin || freebsd || openbsd || netbsd || dragonfly || illumos

package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
	a, err := acquire("a", modeShared, pool, "t", 5*time.Second, 4, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer a.release()
	b, err := acquire("b", modeShared, pool, "t", 500*time.Millisecond, 4, false, nil)
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
	a, err := acquire("a", modeShared, pool, "t", time.Second, 1, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer a.release()
	if _, err := acquire("b", modeShared, pool, "t", 300*time.Millisecond, 1, false, nil); err == nil {
		t.Fatal("second shared lease was granted with maxShared=1")
	}
}

// The shared-lane capacity belongs to the pool: a lease that omits
// --max-shared uses the recorded value, a different explicit value is refused
// while the pool is in use, and an idle pool accepts a new value.
func TestPoolCapacityIsAnInvariant(t *testing.T) {
	pool := testPool(t)
	a, err := acquire("a", modeShared, pool, "t", time.Second, 1, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := status(pool).MaxShared; got != 1 {
		t.Fatalf("status max_shared = %d, want the recorded 1", got)
	}
	if _, err := acquire("b", modeShared, pool, "t", 300*time.Millisecond, 0, false, nil); err == nil {
		t.Fatal("a lease without --max-shared widened the lane past the recorded capacity of 1")
	}
	_, err = acquire("c", modeShared, pool, "t", 300*time.Millisecond, 4, false, nil)
	if err == nil || !strings.Contains(err.Error(), "max_shared is 1") {
		t.Fatalf("mismatched --max-shared while in use: error = %v, want a refusal naming max_shared 1", err)
	}
	a.release()
	d, err := acquire("d", modeShared, pool, "t", time.Second, 4, false, nil)
	if err != nil {
		t.Fatalf("idle pool refused a new capacity: %v", err)
	}
	defer d.release()
	if got := status(pool).MaxShared; got != 4 {
		t.Fatalf("status max_shared = %d, want the newly recorded 4", got)
	}
}

// An exclusive lease waits for every shared holder to finish.
func TestExclusiveWaitsForShared(t *testing.T) {
	pool := testPool(t)
	sh, err := acquire("build", modeShared, pool, "t", time.Second, 4, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := acquire("bench", modeExclusive, pool, "t", 300*time.Millisecond, 4, false, nil); err == nil {
		t.Fatal("exclusive lease was granted while a shared lease was held")
	}
	sh.release()
	ex, err := acquire("bench", modeExclusive, pool, "t", time.Second, 4, false, nil)
	if err != nil {
		t.Fatalf("exclusive lease not granted after shared release: %v", err)
	}
	ex.release()
}

// A waiting exclusive closes the gate: shared work arriving behind it is
// deferred until the exclusive has run, so builds cannot starve a benchmark.
func TestWaitingExclusiveClosesGate(t *testing.T) {
	pool := testPool(t)
	first, err := acquire("build-1", modeShared, pool, "t", time.Second, 4, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	granted := make(chan time.Time, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ex, err := acquire("bench", modeExclusive, pool, "t", 5*time.Second, 4, false, nil)
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
	if _, err := acquire("build-2", modeShared, pool, "t", 200*time.Millisecond, 4, false, nil); err == nil {
		t.Fatal("shared lease was granted while an exclusive was waiting at the gate")
	}
	first.release()
	wg.Wait()
	select {
	case <-granted:
	default:
		t.Fatal("exclusive never ran after the shared holder released")
	}
	late, err := acquire("build-3", modeShared, pool, "t", time.Second, 4, false, nil)
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

// backgroundChild is a shell command that starts `sleep 30` in the background,
// records its PID, and then either waits for it or exits at once.
func backgroundChild(pidFile string, wait bool) []string {
	script := "sleep 30 & echo $! > " + pidFile
	if wait {
		script += "; wait"
	}
	return []string{"sh", "-c", script}
}

func readPID(t *testing.T, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("child PID file: %v", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		t.Fatalf("child PID %q: %v", data, err)
	}
	return pid
}

// requireStopped gives an orphaned, killed child a moment to be reaped, then
// fails if it is still running. A survivor is killed so the test leaks nothing.
func requireStopped(t *testing.T, pid int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for alive(pid) {
		if time.Now().After(deadline) {
			_ = syscall.Kill(pid, syscall.SIGKILL)
			t.Fatalf("process %d outlived its lease", pid)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// A timeout stops the whole process group, not only the direct child, before
// the lease is released.
func TestTimeoutStopsTheProcessGroup(t *testing.T) {
	pool := testPool(t)
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	res, err := run(runOptions{
		argv: backgroundChild(pidFile, true), cwd: t.TempDir(), mode: modeShared, pool: pool, owner: "t",
		timeout: 500 * time.Millisecond, leaseTimeout: time.Second, maxShared: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.RC == 0 {
		t.Fatal("a timed-out command reported rc 0")
	}
	requireStopped(t, readPID(t, pidFile))
	if st := status(pool); len(st.Holders) != 0 {
		t.Fatalf("lease still recorded after the run: %+v", st.Holders)
	}
}

// A command that exits while a process it started keeps running leaves
// nothing behind once the lease is released.
func TestExitedCommandLeavesNoProcesses(t *testing.T) {
	pool := testPool(t)
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	res, err := run(runOptions{
		argv: backgroundChild(pidFile, false), cwd: t.TempDir(), mode: modeShared, pool: pool, owner: "t",
		leaseTimeout: time.Second, maxShared: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.RC != 0 {
		t.Fatalf("rc = %d, want 0", res.RC)
	}
	requireStopped(t, readPID(t, pidFile))
}

// waitForWaiter blocks until the pool records a waiter with this owner, so a
// cancellation test signals a worker that has really reached the queue.
func waitForWaiter(t *testing.T, pool, owner string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for {
		for _, w := range status(pool).Waiters {
			if w.Owner == owner {
				return
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("worker %q never appeared as a waiter in pool %s", owner, pool)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// A worker cancelled while it is queued behind an exclusive holder returns at
// once, never runs its command, and leaves neither a waiter record nor a held
// lock behind.
func TestCancelWhileQueuedForExclusive(t *testing.T) {
	pool := testPool(t)
	holder, err := acquire("holder", modeExclusive, pool, "holder", time.Second, 4, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "ran")
	stop := make(chan os.Signal, 1)
	done := make(chan error, 1)
	go func() {
		_, err := run(runOptions{
			argv: []string{"sh", "-c", "touch " + marker}, cwd: t.TempDir(),
			mode: modeExclusive, pool: pool, owner: "queued", stop: stop,
		})
		done <- err
	}()
	waitForWaiter(t, pool, "queued")
	stop <- syscall.SIGTERM

	select {
	case err := <-done:
		var c *canceled
		if !errors.As(err, &c) {
			t.Fatalf("run returned %v, want a cancellation error", err)
		}
		if want := 128 + int(syscall.SIGTERM); c.exitCode() != want {
			t.Fatalf("exit code = %d, want %d", c.exitCode(), want)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("a queued worker sent SIGTERM did not return")
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("the cancelled worker ran its command")
	}
	if st := status(pool); len(st.Waiters) != 0 {
		t.Fatalf("cancelled worker left a waiter record: %+v", st.Waiters)
	}
	holder.release()
	// The cancelled worker must hold no gate or lease lock: if it did, this
	// shared lease could not be granted.
	after, err := acquire("after", modeShared, pool, "t", 2*time.Second, 4, false, nil)
	if err != nil {
		t.Fatalf("pool still locked by the cancelled worker: %v", err)
	}
	after.release()
}

// A signal that arrives before the worker even joins the queue cancels it
// too: an idle pool is no reason to start the command anyway.
func TestCancelBeforeQueueing(t *testing.T) {
	pool := testPool(t)
	marker := filepath.Join(t.TempDir(), "ran")
	stop := make(chan os.Signal, 1)
	stop <- syscall.SIGHUP
	_, err := run(runOptions{
		argv: []string{"sh", "-c", "touch " + marker}, cwd: t.TempDir(),
		mode: modeExclusive, pool: pool, owner: "early", stop: stop,
	})
	var c *canceled
	if !errors.As(err, &c) {
		t.Fatalf("run returned %v, want a cancellation error", err)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("a worker signalled before queueing still ran its command")
	}
	if st := status(pool); len(st.Waiters) != 0 || len(st.Holders) != 0 {
		t.Fatalf("pool not clean after an early cancellation: %+v", st)
	}
}

// A shared worker cancelled while it waits for a slot gives up the shared
// lease it already holds, so cancelling it really does stop it delaying an
// exclusive benchmark.
func TestCancelWhileWaitingForSharedSlot(t *testing.T) {
	pool := testPool(t)
	occupant, err := acquire("occupant", modeShared, pool, "occupant", time.Second, 1, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "ran")
	stop := make(chan os.Signal, 1)
	done := make(chan error, 1)
	go func() {
		_, err := run(runOptions{
			argv: []string{"sh", "-c", "touch " + marker}, cwd: t.TempDir(),
			mode: modeShared, pool: pool, owner: "queued", maxShared: 1, stop: stop,
		})
		done <- err
	}()
	waitForWaiter(t, pool, "queued")
	stop <- syscall.SIGINT

	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "canceled") {
			t.Fatalf("run returned %v, want a cancellation error", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("a shared worker waiting for a slot did not return after SIGINT")
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("the cancelled worker ran its command")
	}
	occupant.release()
	// The cancelled waiter must no longer hold the shared lease, or this
	// exclusive lease would keep waiting for it.
	ex, err := acquire("bench", modeExclusive, pool, "t", 2*time.Second, 0, false, nil)
	if err != nil {
		t.Fatalf("cancelled shared waiter still holds the lease: %v", err)
	}
	ex.release()
	if st := status(pool); len(st.Waiters) != 0 {
		t.Fatalf("cancelled worker left a waiter record: %+v", st.Waiters)
	}
}

// A signal to the baton stops the wrapped process group and reports
// 128+signal, the shell convention for a command ended by a signal.
func TestStopSignalStopsTheProcessGroup(t *testing.T) {
	pool := testPool(t)
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	stop := make(chan os.Signal, 1)
	go func() {
		time.Sleep(500 * time.Millisecond)
		stop <- syscall.SIGTERM
	}()
	res, err := run(runOptions{
		argv: backgroundChild(pidFile, true), cwd: t.TempDir(), mode: modeShared, pool: pool, owner: "t",
		leaseTimeout: time.Second, maxShared: 4, stop: stop,
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := 128 + int(syscall.SIGTERM); res.RC != want {
		t.Fatalf("rc = %d, want %d", res.RC, want)
	}
	requireStopped(t, readPID(t, pidFile))
}
