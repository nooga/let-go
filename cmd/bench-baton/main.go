//go:build linux || darwin || freebsd || openbsd || netbsd || dragonfly || illumos

// Command bench-baton is a machine-quiescence lease for benchmarks, with a
// concurrent lane for builds and tests.
//
// The scarce resource on a developer machine is a quiet CPU, not the build
// tree. A benchmark or a perf gate (make bench-ratchet) that runs while a
// test suite or a second benchmark is executing produces numbers that are
// wrong in ways a ratchet cannot distinguish from a regression. The baton
// arbitrates one pool per machine-level resource with two lanes:
//
//	bench  (exclusive)  benchmarks, perf gates, timing-sensitive runs.
//	                    Holds the pool alone: no other exclusive run and no
//	                    shared run overlaps it.
//	build  (shared)     builds, test suites, regeneration, lint. Runs
//	                    concurrently with other shared work up to the pool's
//	                    max_shared slots, but never while an exclusive run
//	                    holds or is waiting for the pool.
//
// A waiting exclusive closes the gate, so a stream of builds cannot starve
// a benchmark: shared work arriving behind it is deferred until it has run.
//
// Every lease is an OS file lock (flock), so a holder that dies releases
// the pool without cleanup; the JSON state file is bookkeeping only, and
// `reap` tidies records whose process is gone.
//
// The wrapped command runs in its own process group. A timeout, a signal to
// the baton, or the command exiting while processes it started live on all
// stop the whole group before the lease is released, so no workload outlives
// its lease.
//
// Usage:
//
//	go run ./cmd/bench-baton bench --owner pr641 --cwd WT -- make bench-ratchet
//	go run ./cmd/bench-baton build --owner pr641 --cwd WT -- go test ./...
//	go run ./cmd/bench-baton run --mode shared -- make generate
//	go run ./cmd/bench-baton status
//	go run ./cmd/bench-baton ledger -n 20
//	go run ./cmd/bench-baton reap
//
// build and bench exit with the wrapped command's exit code and print its
// tail, so they drop into existing scripts unchanged; the full output is in
// the log file the ledger names.
//
// State lives under $BENCH_BATON_HOME, default $XDG_CACHE_HOME/let-go/batons
// (os.UserCacheDir), one directory per --pool (default "letgo"): lease.lock,
// gate.lock, slotN.lock, state.json, ledger.jsonl, logs/. Point
// BENCH_BATON_HOME at an existing pool directory to keep its ledger.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	modeExclusive = "exclusive"
	modeShared    = "shared"

	defaultPool      = "letgo"
	defaultMaxShared = 4
	pollInterval     = 250 * time.Millisecond
	stopGrace        = 2 * time.Second // SIGTERM to SIGKILL for a stopped process group
	logHeaderLines   = 5               // the "# mode/intent/cwd/cmd" preamble plus its blank line
)

// token is one holder or waiter. Field names match the JSON the Python
// predecessor wrote, so an existing pool's state.json and ledger.jsonl
// stay readable.
type token struct {
	ID      string  `json:"id"`
	Owner   string  `json:"owner"`
	Intent  string  `json:"intent"`
	Mode    string  `json:"mode"`
	PID     int     `json:"pid"`
	Since   float64 `json:"since"`
	Granted float64 `json:"granted,omitempty"`
	Slot    *int    `json:"slot,omitempty"`
}

type state struct {
	Holders []token `json:"holders"`
	Waiters []token `json:"waiters"`
	// MaxShared is the pool's shared-lane capacity, recorded by the first
	// lease of an idle pool. Zero means not yet recorded.
	MaxShared int `json:"max_shared,omitempty"`
}

// statusView is token plus liveness, as `status` prints it.
type statusView struct {
	token
	Alive    bool    `json:"alive"`
	HeldS    float64 `json:"held_s,omitempty"`
	WaitingS float64 `json:"waiting_s,omitempty"`
}

type statusReport struct {
	Pool             string       `json:"pool"`
	MaxShared        int          `json:"max_shared"`
	ExclusiveHeld    bool         `json:"exclusive_held"`
	ExclusiveWaiting bool         `json:"exclusive_waiting"`
	Holders          []statusView `json:"holders"`
	Waiters          []statusView `json:"waiters"`
	Dir              string       `json:"dir"`
}

func rootDir() string {
	if h := os.Getenv("BENCH_BATON_HOME"); h != "" {
		return h
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		cache = filepath.Join(os.Getenv("HOME"), ".cache")
	}
	return filepath.Join(cache, "let-go", "batons")
}

func poolDir(pool string) string {
	d := filepath.Join(rootDir(), pool)
	_ = os.MkdirAll(filepath.Join(d, "logs"), 0o755)
	return d
}

func poolFile(pool, name string) string { return filepath.Join(poolDir(pool), name) }

func alive(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

func now() float64 { return float64(time.Now().UnixNano()) / 1e9 }

// flockUntil takes `how` (LOCK_EX or LOCK_SH) on f, polling until deadline.
// A zero deadline blocks indefinitely.
func flockUntil(f *os.File, how int, deadline time.Time, what string) error {
	if deadline.IsZero() {
		return syscall.Flock(int(f.Fd()), how)
	}
	for {
		err := syscall.Flock(int(f.Fd()), how|syscall.LOCK_NB)
		if err == nil {
			return nil
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) && !errors.Is(err, syscall.EAGAIN) {
			return err
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out acquiring %s", what)
		}
		time.Sleep(pollInterval)
	}
}

// mutateState applies fn to the pool's bookkeeping under the state lock and
// writes it back atomically.
func mutateState(pool string, fn func(*state)) error {
	lock, err := os.OpenFile(poolFile(pool, "state.lock"), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN) //nolint:errcheck
	st := readState(pool)
	fn(&st)
	data, err := json.MarshalIndent(st, "", " ")
	if err != nil {
		return err
	}
	tmp := poolFile(pool, "state.json.tmp")
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, poolFile(pool, "state.json"))
}

func readState(pool string) state {
	var st state
	data, err := os.ReadFile(poolFile(pool, "state.json"))
	if err == nil {
		_ = json.Unmarshal(data, &st)
	}
	if st.Holders == nil {
		st.Holders = []token{}
	}
	if st.Waiters == nil {
		st.Waiters = []token{}
	}
	return st
}

// capacity is the pool's shared-lane capacity: the recorded value, or the
// default for a pool that has not recorded one.
func capacity(st state) int {
	if st.MaxShared > 0 {
		return st.MaxShared
	}
	return defaultMaxShared
}

// settleCapacity fixes the shared-lane capacity a new lease uses. The slot
// count is a property of the pool, not of one invocation: a worker asking for
// a different --max-shared must not silently widen or narrow the lane under
// another worker's lease. An idle pool records the requested value (or the
// default); a pool with holders or waiters keeps its recorded value, and an
// explicit different request is refused. requested <= 0 means "use the pool's".
func settleCapacity(st *state, requested int) (int, error) {
	busy := len(st.Holders) > 0 || len(st.Waiters) > 0
	switch {
	case requested <= 0:
		if st.MaxShared == 0 {
			st.MaxShared = defaultMaxShared
		}
		return st.MaxShared, nil
	case st.MaxShared == requested:
		return requested, nil
	case st.MaxShared == 0 || !busy:
		st.MaxShared = requested
		return requested, nil
	default:
		return 0, fmt.Errorf("pool max_shared is %d while the pool is in use; requested %d (omit --max-shared, or run `reap` if its holders are gone)",
			st.MaxShared, requested)
	}
}

func status(pool string) statusReport {
	st := readState(pool)
	t := now()
	rep := statusReport{Pool: pool, MaxShared: capacity(st), Dir: poolDir(pool),
		Holders: []statusView{}, Waiters: []statusView{}}
	for _, h := range st.Holders {
		rep.Holders = append(rep.Holders, statusView{token: h, Alive: alive(h.PID), HeldS: round1(t - h.Granted)})
		if h.Mode == modeExclusive {
			rep.ExclusiveHeld = true
		}
	}
	for _, w := range st.Waiters {
		rep.Waiters = append(rep.Waiters, statusView{token: w, Alive: alive(w.PID), WaitingS: round1(t - w.Since)})
		if w.Mode == modeExclusive {
			rep.ExclusiveWaiting = true
		}
	}
	return rep
}

func round1(x float64) float64 { return float64(int(x*10+0.5)) / 10 }

// reap drops holders and waiters whose process is gone. Their locks were
// released by the kernel when the process exited; only the record lingers.
func reap(pool string) ([]token, error) {
	var dropped []token
	err := mutateState(pool, func(st *state) {
		keep := func(ts []token) []token {
			out := ts[:0:0]
			for _, t := range ts {
				if alive(t.PID) {
					out = append(out, t)
				} else {
					dropped = append(dropped, t)
				}
			}
			return out
		}
		st.Holders = keep(st.Holders)
		st.Waiters = keep(st.Waiters)
	})
	return dropped, err
}

// lease is a granted pool lease; release returns it.
type lease struct {
	tok   token
	pool  string
	files []*os.File // gate, lease, slot (shared only); unlocked and closed on release
	t0    time.Time
	quiet bool
}

// acquire queues for the pool in `mode` and returns once granted.
//
// Exclusive: take the gate exclusively (so new shared work queues behind
// us), then the lease exclusively (so current shared holders drain).
// Shared: take the gate shared (blocked only while an exclusive holds or
// waits for it), drop it, take the lease shared, then one of the pool's
// max_shared slots. maxShared <= 0 uses the pool's recorded capacity.
func acquire(intent, mode, pool, owner string, timeout time.Duration, maxShared int, verbose bool) (*lease, error) {
	if mode != modeExclusive && mode != modeShared {
		return nil, fmt.Errorf("mode must be %q or %q", modeExclusive, modeShared)
	}
	if owner == "" {
		owner = os.Getenv("BENCH_BATON_OWNER")
	}
	if owner == "" {
		owner = fmt.Sprintf("pid%d", os.Getpid())
	}
	tok := token{ID: newID(), Owner: owner, Intent: intent, Mode: mode, PID: os.Getpid(), Since: now()}
	var deadline time.Time
	if timeout > 0 {
		deadline = time.Now().Add(timeout)
	}
	// Settle the capacity and join the queue under one state lock, so two
	// idle-pool leases cannot each record a different capacity.
	var capErr error
	if err := mutateState(pool, func(st *state) {
		c, err := settleCapacity(st, maxShared)
		if err != nil {
			capErr = err
			return
		}
		maxShared = c
		st.Waiters = append(st.Waiters, tok)
	}); err != nil {
		return nil, err
	}
	if capErr != nil {
		return nil, capErr
	}
	if verbose {
		s := status(pool)
		fmt.Printf("[baton:%s] queued %s %s: %s (holders=%d, waiters=%d)\n",
			pool, mode, owner, intent, len(s.Holders), len(s.Waiters))
	}
	l := &lease{tok: tok, pool: pool, t0: time.Now(), quiet: !verbose}
	open := func(name string) (*os.File, error) {
		f, err := os.OpenFile(poolFile(pool, name), os.O_CREATE|os.O_RDWR, 0o644)
		if err == nil {
			l.files = append(l.files, f)
		}
		return f, err
	}
	fail := func(err error) (*lease, error) {
		l.unlockAll()
		_ = mutateState(pool, func(st *state) { st.Waiters = without(st.Waiters, tok.ID) })
		return nil, err
	}
	gate, err := open("gate.lock")
	if err != nil {
		return fail(err)
	}
	lock, err := open("lease.lock")
	if err != nil {
		return fail(err)
	}
	if mode == modeExclusive {
		if err := flockUntil(gate, syscall.LOCK_EX, deadline, "gate"); err != nil {
			return fail(err)
		}
		if err := flockUntil(lock, syscall.LOCK_EX, deadline, "exclusive lease"); err != nil {
			return fail(err)
		}
	} else {
		if err := flockUntil(gate, syscall.LOCK_SH, deadline, "gate"); err != nil {
			return fail(err)
		}
		_ = syscall.Flock(int(gate.Fd()), syscall.LOCK_UN)
		if err := flockUntil(lock, syscall.LOCK_SH, deadline, "shared lease"); err != nil {
			return fail(err)
		}
		slot, idx, err := acquireSlot(pool, maxShared, deadline)
		if err != nil {
			return fail(err)
		}
		l.files = append(l.files, slot)
		l.tok.Slot = &idx
	}
	l.tok.Granted = now()
	granted := l.tok
	if err := mutateState(pool, func(st *state) {
		st.Waiters = without(st.Waiters, tok.ID)
		st.Holders = append(st.Holders, granted)
	}); err != nil {
		return fail(err)
	}
	if verbose {
		fmt.Printf("[baton:%s] HELD %s by %s: %s (waited %.1fs)\n",
			pool, mode, owner, intent, time.Since(l.t0).Seconds())
	}
	return l, nil
}

func acquireSlot(pool string, maxShared int, deadline time.Time) (*os.File, int, error) {
	for {
		for i := 0; i < maxShared; i++ {
			f, err := os.OpenFile(poolFile(pool, fmt.Sprintf("slot%d.lock", i)), os.O_CREATE|os.O_RDWR, 0o644)
			if err != nil {
				return nil, 0, err
			}
			if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err == nil {
				return f, i, nil
			}
			f.Close()
		}
		if !deadline.IsZero() && time.Now().After(deadline) {
			return nil, 0, fmt.Errorf("no shared slot in pool %s (max_shared=%d)", pool, maxShared)
		}
		time.Sleep(pollInterval)
	}
}

func (l *lease) unlockAll() {
	for _, f := range l.files {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
	}
	l.files = nil
}

func (l *lease) release() {
	_ = mutateState(l.pool, func(st *state) {
		st.Holders = without(st.Holders, l.tok.ID)
		st.Waiters = without(st.Waiters, l.tok.ID)
	})
	l.unlockAll()
	if !l.quiet && l.tok.Granted > 0 {
		fmt.Printf("[baton:%s] RELEASED %s %s after %.1fs\n",
			l.pool, l.tok.Mode, l.tok.Owner, now()-l.tok.Granted)
	}
}

func without(ts []token, id string) []token {
	out := ts[:0:0]
	for _, t := range ts {
		if t.ID != id {
			out = append(out, t)
		}
	}
	return out
}

func newID() string {
	var b [4]byte
	f, err := os.Open("/dev/urandom")
	if err == nil {
		_, _ = f.Read(b[:])
		f.Close()
	}
	return fmt.Sprintf("%x", b)
}

// groupAlive reports whether any process in process group pgid still exists.
func groupAlive(pgid int) bool {
	err := syscall.Kill(-pgid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

// waitGroupGone polls until process group pgid is empty or d elapses, and
// reports whether it emptied.
func waitGroupGone(pgid int, d time.Duration) bool {
	deadline := time.Now().Add(d)
	for groupAlive(pgid) {
		if time.Now().After(deadline) {
			return false
		}
		time.Sleep(20 * time.Millisecond)
	}
	return true
}

// stopGroup ends every process in process group pgid and waits until the
// group is empty: SIGTERM, a grace period, then SIGKILL. Killing only the
// direct child would leave anything it started running after the lease is
// released.
func stopGroup(pgid int) {
	if !groupAlive(pgid) {
		return
	}
	_ = syscall.Kill(-pgid, syscall.SIGTERM)
	if waitGroupGone(pgid, stopGrace) {
		return
	}
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
	waitGroupGone(pgid, stopGrace)
}

// runResult is one ledger row; JSON names match the Python ledger.
type runResult struct {
	TS       string   `json:"ts"`
	Owner    string   `json:"owner"`
	RC       int      `json:"rc"`
	Mode     string   `json:"mode"`
	Intent   string   `json:"intent"`
	Cmd      []string `json:"cmd"`
	Cwd      string   `json:"cwd"`
	Duration float64  `json:"duration_s"`
	Waited   float64  `json:"waited_s"`
	Log      string   `json:"log"`
	Tail     string   `json:"-"`
}

type runOptions struct {
	argv         []string
	cwd          string
	intent       string
	mode         string
	pool         string
	owner        string
	timeout      time.Duration // wrapped command; zero = none
	leaseTimeout time.Duration // waiting for the pool; zero = block
	maxShared    int           // zero = the pool's recorded capacity
	tailLines    int
	verbose      bool
	// stop delivers a signal sent to the baton itself. Receiving one stops the
	// wrapped command's process group before the lease is released. Nil means
	// the run is never interrupted.
	stop <-chan os.Signal
}

// run executes argv in the requested lane. Full output goes to a log under
// the pool; the tail comes back for the caller to print.
func run(o runOptions) (runResult, error) {
	if o.intent == "" {
		o.intent = strings.Join(o.argv, " ")
		if len(o.intent) > 120 {
			o.intent = o.intent[:120]
		}
	}
	if o.tailLines <= 0 {
		o.tailLines = 40
	}
	if o.cwd == "" {
		o.cwd = "."
	}
	logPath := filepath.Join(poolDir(o.pool), "logs",
		fmt.Sprintf("%s-%s.log", time.Now().UTC().Format("20060102T150405"), newID()[:6]))
	tReq := time.Now()
	l, err := acquire(o.intent, o.mode, o.pool, o.owner, o.leaseTimeout, o.maxShared, o.verbose)
	if err != nil {
		return runResult{}, err
	}
	waited := time.Since(tReq)
	logf, err := os.Create(logPath)
	if err != nil {
		l.release()
		return runResult{}, err
	}
	fmt.Fprintf(logf, "# mode: %s\n# intent: %s\n# cwd: %s\n# cmd: %q\n\n", o.mode, o.intent, o.cwd, o.argv) // logHeaderLines lines
	cmd := exec.Command(o.argv[0], o.argv[1:]...)
	cmd.Dir = o.cwd
	cmd.Stdout = logf
	cmd.Stderr = logf
	cmd.Env = os.Environ()
	// Its own process group, so stopGroup reaches everything the command
	// starts, not only the direct child.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	t0 := time.Now()
	runErr := cmd.Start()
	rc := 0
	var interrupted syscall.Signal
	if runErr == nil {
		pgid := cmd.Process.Pid
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		var timedOut <-chan time.Time
		if o.timeout > 0 {
			timer := time.NewTimer(o.timeout)
			defer timer.Stop()
			timedOut = timer.C
		}
		select {
		case runErr = <-done:
		case <-timedOut:
			runErr = fmt.Errorf("timed out after %s", o.timeout)
			stopGroup(pgid)
			<-done
		case sig := <-o.stop:
			if s, ok := sig.(syscall.Signal); ok {
				interrupted = s
			}
			runErr = fmt.Errorf("interrupted by %v", sig)
			stopGroup(pgid)
			<-done
		}
		// The direct child can exit while processes it started keep running;
		// none of them may outlive the lease.
		stopGroup(pgid)
	}
	var exitErr *exec.ExitError
	switch {
	case interrupted != 0:
		rc = 128 + int(interrupted)
		fmt.Fprintf(logf, "\n# bench-baton: %v\n", runErr)
	case runErr == nil:
	case errors.As(runErr, &exitErr):
		rc = exitErr.ExitCode()
	default:
		rc = 127
		fmt.Fprintf(logf, "\n# bench-baton: %v\n", runErr)
	}
	dur := time.Since(t0)
	logf.Close()
	l.release()

	data, _ := os.ReadFile(logPath)
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	// The tail is the command's output, not the log's own header.
	if len(lines) > logHeaderLines {
		lines = lines[logHeaderLines:]
	} else {
		lines = nil
	}
	if len(lines) > o.tailLines {
		lines = lines[len(lines)-o.tailLines:]
	}
	res := runResult{
		TS: time.Now().Format("2006-01-02T15:04:05"), Owner: l.tok.Owner, RC: rc, Mode: o.mode,
		Intent: o.intent, Cmd: o.argv, Cwd: o.cwd, Duration: round2(dur.Seconds()),
		Waited: round2(waited.Seconds()), Log: logPath, Tail: strings.Join(lines, "\n"),
	}
	if row, err := json.Marshal(res); err == nil {
		if lf, err := os.OpenFile(poolFile(o.pool, "ledger.jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644); err == nil {
			lf.Write(append(row, '\n')) //nolint:errcheck
			lf.Close()
		}
	}
	if o.verbose {
		verdict := "PASS"
		if rc != 0 {
			verdict = "FAIL"
		}
		fmt.Printf("[baton:%s] %s rc=%d %.1fs (%s) :: %s\n", o.pool, verdict, rc, dur.Seconds(), o.mode, o.intent)
	}
	return res, nil
}

func round2(x float64) float64 { return float64(int(x*100+0.5)) / 100 }

func ledger(pool string, n int) ([]runResult, error) {
	data, err := os.ReadFile(poolFile(pool, "ledger.jsonl"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var rows []runResult
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		var r runResult
		if json.Unmarshal([]byte(line), &r) == nil {
			rows = append(rows, r)
		}
	}
	if n > 0 && len(rows) > n {
		rows = rows[len(rows)-n:]
	}
	return rows, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: bench-baton [--pool P] (build|bench|run [--mode M]) [--owner O] [--cwd D] [--intent S] [--timeout SEC] [--lease-timeout SEC] [--max-shared N] -- CMD ...")
	fmt.Fprintln(os.Stderr, "       bench-baton [--pool P] status | reap | ledger [-n N]")
	os.Exit(2)
}

func main() {
	os.Exit(mainWithArgs(os.Args[1:]))
}

func mainWithArgs(args []string) int {
	pool := defaultPool
	if len(args) >= 2 && args[0] == "--pool" {
		pool, args = args[1], args[2:]
	}
	if len(args) == 0 {
		usage()
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "build", "bench", "run":
		fs := flag.NewFlagSet(sub, flag.ExitOnError)
		mode := modeShared
		if sub == "bench" {
			mode = modeExclusive
		}
		if sub == "run" {
			fs.StringVar(&mode, "mode", modeShared, "shared or exclusive")
		}
		owner := fs.String("owner", "", "stable worker name for status/ledger")
		cwd := fs.String("cwd", ".", "working directory for CMD")
		intent := fs.String("intent", "", "what this run is for (default: the command)")
		timeout := fs.Float64("timeout", 0, "stop CMD and its process group after this many seconds")
		leaseTimeout := fs.Float64("lease-timeout", 0, "give up waiting for the pool after this many seconds")
		maxShared := fs.Int("max-shared", 0, fmt.Sprintf("shared-lane capacity for an idle pool (default: the pool's recorded value, or %d)", defaultMaxShared))
		// Split at "--": flags before it, CMD after.
		var cmdArgs []string
		for i, a := range rest {
			if a == "--" {
				cmdArgs = rest[i+1:]
				rest = rest[:i]
				break
			}
		}
		_ = fs.Parse(rest)
		if len(cmdArgs) == 0 {
			cmdArgs = fs.Args()
		}
		if len(cmdArgs) == 0 {
			fmt.Fprintln(os.Stderr, "bench-baton: no command given after --")
			return 2
		}
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer signal.Stop(stop)
		res, err := run(runOptions{
			argv: cmdArgs, cwd: *cwd, intent: *intent, mode: mode, pool: pool, owner: *owner,
			timeout: secs(*timeout), leaseTimeout: secs(*leaseTimeout), maxShared: *maxShared,
			verbose: true, stop: stop,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "bench-baton: %v\n", err)
			return 1
		}
		fmt.Println(res.Tail)
		return res.RC
	case "status":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", " ")
		_ = enc.Encode(status(pool))
		return 0
	case "reap":
		dropped, err := reap(pool)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bench-baton: %v\n", err)
			return 1
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", " ")
		_ = enc.Encode(map[string]any{"dropped": dropped})
		return 0
	case "ledger":
		fs := flag.NewFlagSet("ledger", flag.ExitOnError)
		n := fs.Int("n", 20, "rows to show")
		_ = fs.Parse(rest)
		rows, err := ledger(pool, *n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bench-baton: %v\n", err)
			return 1
		}
		for _, r := range rows {
			fmt.Printf("%s %-4s rc=%d %.1fs w=%.1f %s :: %s\n", r.TS, r.Mode[:4], r.RC, r.Duration, r.Waited, r.Owner, r.Intent)
		}
		return 0
	}
	usage()
	return 2
}

func secs(f float64) time.Duration { return time.Duration(f * float64(time.Second)) }
