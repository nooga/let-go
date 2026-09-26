package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runHook runs scripts/bench-ratchet-if-affected.sh from the repository root
// with a fake `make` first on PATH, and reports whether the hook ran the
// ratchet together with its stderr. The fake stands in for `make
// bench-ratchet`, so these tests exercise the hook's decision without running
// a benchmark.
func runHook(t *testing.T, env ...string) (ran bool, stderr string) {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	marker := filepath.Join(bin, "ran")
	fake := "#!/bin/sh\necho \"$@\" > '" + marker + "'\n"
	if err := os.WriteFile(filepath.Join(bin, "make"), []byte(fake), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("bash", filepath.Join(root, "scripts", "bench-ratchet-if-affected.sh"))
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	cmd.Env = append(cmd.Env, env...)
	var errb strings.Builder
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		t.Fatalf("hook failed: %v\n%s", err, errb.String())
	}
	got, err := os.ReadFile(marker)
	if err != nil {
		return false, errb.String()
	}
	if strings.TrimSpace(string(got)) != "bench-ratchet" {
		t.Fatalf("hook invoked make %q, want bench-ratchet", got)
	}
	return true, errb.String()
}

// TestHookFailsClosedOnUnresolvableTip pins the fail-closed direction for the
// pushed ref. An all-zeros ref is what a pre-push hook receives for a ref that
// does not exist locally; nothing about the push can be decided from it.
func TestHookFailsClosedOnUnresolvableTip(t *testing.T) {
	ran, stderr := runHook(t, "PRE_COMMIT_TO_REF="+strings.Repeat("0", 40))
	if !ran {
		t.Fatalf("an unresolvable pushed tip must run the ratchet; stderr:\n%s", stderr)
	}
	if !strings.Contains(stderr, "could not determine the pushed range") {
		t.Errorf("stderr should say the range was undeterminable:\n%s", stderr)
	}
}

// TestHookRunsWhenPushedTipIsNotTheCheckout pins the review finding that the
// hook classified the checkout instead of the pushed ref. The closure and the
// deletion check both read the checked-out tree, so a push of any other tree
// cannot be scoped from here and must run.
func TestHookRunsWhenPushedTipIsNotTheCheckout(t *testing.T) {
	other := otherCommit(t)
	ran, stderr := runHook(t, "PRE_COMMIT_TO_REF="+other)
	if !ran {
		t.Fatalf("pushing %s, which is not the checked-out tree, must run the ratchet; stderr:\n%s", other, stderr)
	}
	if !strings.Contains(stderr, "differs from the checked-out tree") {
		t.Errorf("stderr should name the tip/checkout mismatch:\n%s", stderr)
	}
}

// otherCommit returns a commit whose tree differs from the checkout: the
// parent of the checkout's commit in a git checkout, and the grandparent of @
// under jj, whose @ is usually an empty child of the checked-out commit. It
// skips when history is too shallow to have one (a depth-1 CI clone).
func otherCommit(t *testing.T) string {
	t.Helper()
	root := filepath.Join("..", "..")
	if _, err := exec.LookPath("jj"); err == nil {
		c := exec.Command("jj", "workspace", "root")
		c.Dir = root
		if c.Run() == nil {
			return commandOutput(t, root, "jj", "log", "--no-graph", "--ignore-working-copy", "-r", "@--", "-T", "commit_id")
		}
	}
	return commandOutput(t, root, "git", "rev-parse", "--verify", "--quiet", "HEAD~1^{commit}")
}

func commandOutput(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	c := exec.Command(name, args...)
	c.Dir = dir
	out, err := c.Output()
	sha := strings.TrimSpace(string(out))
	if err != nil || sha == "" {
		t.Skipf("no commit other than the checkout is available (%s %s): %v", name, strings.Join(args, " "), err)
	}
	return sha
}
