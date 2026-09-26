package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestJJWorkspaceSHAResolvesThisWorkspace pins the property the capture SHA
// exists for: it must name a commit that describes the tree being measured.
//
// The bug this guards is specific to jj and invisible to `git rev-parse`: a
// colocated jj repo keeps ONE git HEAD for every workspace, so each workspace
// reports whatever the default workspace last exported. This test asserts the
// resolved SHA is this workspace's own commit, which is what makes it differ
// from git HEAD in the failing case.
//
// It needs a jj repo to mean anything, so it skips elsewhere rather than
// asserting a fallback it cannot set up.
func TestJJWorkspaceSHAResolvesThisWorkspace(t *testing.T) {
	if _, err := exec.LookPath("jj"); err != nil {
		t.Skip("jj not installed")
	}
	if err := exec.Command("jj", "workspace", "root").Run(); err != nil {
		t.Skip("not inside a jj workspace")
	}

	got := jjWorkspaceSHA()
	if got == "" {
		t.Fatal("jjWorkspaceSHA returned empty inside a jj workspace")
	}
	if len(got) != 12 {
		t.Fatalf("jjWorkspaceSHA() = %q, want a 12-character id", got)
	}

	// The answer must be a commit jj can resolve. A stamp that resolves to
	// nothing is the whole defect, so assert resolvability rather than shape.
	if err := exec.Command("jj", "log", "--no-graph", "--ignore-working-copy",
		"-r", got, "-T", "commit_id").Run(); err != nil {
		t.Fatalf("jjWorkspaceSHA() = %q, which jj cannot resolve: %v", got, err)
	}

	// And it must be THIS workspace's commit: @ when the working copy differs
	// from its parent, the parent when it does not. Deriving the expectation
	// separately here, with --ignore-working-copy, keeps the assertion from
	// simply restating the implementation's own command.
	empty := jjOut(t, "--ignore-working-copy", "-r", "@", "-T", `if(empty,"y","n")`)
	rev := "@"
	if empty == "y" {
		rev = "@-"
	}
	want := jjOut(t, "--ignore-working-copy", "-r", rev, "-T", "commit_id.short(12)")
	if strings.Fields(want) == nil || len(strings.Fields(want)) != 1 {
		t.Skipf("%s is a merge or has no single parent; nothing to compare", rev)
	}
	if got != want {
		t.Errorf("jjWorkspaceSHA() = %q, want %q (workspace %s, @ empty = %q)", got, want, rev, empty)
	}
}

// TestCaptureSHAPrefersWorkspaceOverSharedGitHEAD documents the regression
// directly: where the shared git HEAD disagrees with this workspace's commit,
// gitShortSHA must follow the workspace. When they agree the test still runs
// but proves only consistency, so it says which case it observed.
func TestCaptureSHAPrefersWorkspaceOverSharedGitHEAD(t *testing.T) {
	if _, err := exec.LookPath("jj"); err != nil {
		t.Skip("jj not installed")
	}
	if err := exec.Command("jj", "workspace", "root").Run(); err != nil {
		t.Skip("not inside a jj workspace")
	}
	ws := jjWorkspaceSHA()
	if ws == "" {
		t.Fatal("jjWorkspaceSHA returned empty inside a jj workspace")
	}
	if got := gitShortSHA(); got != ws {
		t.Errorf("gitShortSHA() = %q, want the workspace commit %q", got, ws)
	}

	out, err := exec.Command("git", "rev-parse", "--short=12", "HEAD").Output()
	if err != nil {
		t.Skip("no git HEAD to compare against")
	}
	if head := strings.TrimSpace(string(out)); head == ws {
		t.Logf("git HEAD (%s) happens to match this workspace; "+
			"the divergent case is not exercised in this checkout", head)
	} else {
		t.Logf("git HEAD (%s) differs from this workspace (%s) — the regression case", head, ws)
	}
}

// TestCaptureSHAIsResolvableByPlainGit pins the property that makes a recorded
// stamp evidence rather than a label: someone with git and no jj must be able
// to resolve it. jj stores commits in the git object store, so this holds for
// a jj commit id too — but only once jj has exported, which is why
// jjWorkspaceSHA verifies it rather than assuming it.
func TestCaptureSHAIsResolvableByPlainGit(t *testing.T) {
	sha := gitShortSHA()
	if sha == "" {
		t.Skip("no SHA resolvable here")
	}
	if err := exec.Command("git", "cat-file", "-e", sha+"^{commit}").Run(); err != nil {
		t.Fatalf("recorded SHA %q is not a git commit object: %v", sha, err)
	}
	// Type, not just existence: a tree or blob id would satisfy -e on its own
	// spelling but is not something anyone can check out.
	out, err := exec.Command("git", "cat-file", "-t", sha).Output()
	if err != nil {
		t.Fatalf("git cat-file -t %q: %v", sha, err)
	}
	if got := strings.TrimSpace(string(out)); got != "commit" {
		t.Errorf("recorded SHA %q is a %s, want a commit", sha, got)
	}
}

// TestPublishedCheckAgreesWithGitReachability pins capturedSHAIsPublished
// against git's own answer, since the warning it drives is the only signal an
// operator gets that a baseline they are about to write cannot be reproduced
// elsewhere. A commit reachable from a remote-tracking ref must read as
// published; main@upstream's tip is the case that must never warn.
func TestPublishedCheckAgreesWithGitReachability(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	out, err := exec.Command("git", "for-each-ref", "--format=%(objectname:short=12)",
		"--count=1", "refs/remotes/").Output()
	remote := strings.TrimSpace(string(out))
	if err != nil || remote == "" {
		t.Skip("no remote-tracking refs to test against")
	}
	if !capturedSHAIsPublished(remote) {
		t.Errorf("capturedSHAIsPublished(%q) = false for a remote-tracking tip, want true", remote)
	}
}

// TestShellAndGoSHAResolversAgree pins scripts/current-commit.sh against
// jjWorkspaceSHA. The Makefile cannot call Go code to stamp -X main.commit, so
// the same rule necessarily exists twice; two copies of a rule drift, and a
// drift here is silent — each would keep returning a plausible sha. Since that
// stamp becomes let-go.semver/current-commit, which SHA-pinned require-letgo
// constraints prefix-match against, a divergence would make those constraints
// evaluate against a commit the binary was not built from.
func TestShellAndGoSHAResolversAgree(t *testing.T) {
	script := filepath.Join("..", "..", "scripts", "current-commit.sh")
	if _, err := os.Stat(script); err != nil {
		t.Skipf("script not present: %v", err)
	}
	out, err := exec.Command("sh", script).Output()
	if err != nil {
		t.Fatalf("current-commit.sh: %v", err)
	}
	shell := strings.TrimSpace(string(out))
	if shell == "" {
		t.Fatal("current-commit.sh printed nothing")
	}

	goSHA := gitShortSHA()
	if goSHA == "" {
		t.Skip("no SHA resolvable here")
	}
	if shell != goSHA {
		t.Errorf("current-commit.sh = %q, jjWorkspaceSHA/gitShortSHA = %q — the shell and Go resolvers have drifted", shell, goSHA)
	}
	// Whatever they agree on must be a commit plain git can resolve; that is
	// the property a non-jj user depends on.
	if shell != "none" {
		if err := exec.Command("git", "cat-file", "-e", shell+"^{commit}").Run(); err != nil {
			t.Errorf("resolved SHA %q is not a git commit object: %v", shell, err)
		}
	}
}

// TestDefaultRunPathUsesResolvedSHA pins that the capture filename carries the
// same id, since that filename is how a .jsonl is matched back to its commit
// when a baseline is rebuilt later.
func TestDefaultRunPathUsesResolvedSHA(t *testing.T) {
	sha := gitShortSHA()
	if sha == "" {
		t.Skip("no SHA resolvable here")
	}
	base := filepath.Base(defaultRunPath())
	if !strings.HasPrefix(base, sha+"-") {
		t.Errorf("defaultRunPath() = %q, want it to start with %q", base, sha+"-")
	}
	if filepath.Ext(base) != ".jsonl" {
		t.Errorf("defaultRunPath() = %q, want a .jsonl", base)
	}
}

// TestResolversWorkInASecondaryWorkspaceOutsideThePrimary pins the layout the
// review reproduced: `jj workspace add` to a directory outside the primary
// checkout gives a working tree with .jj but no .git and no Git-owning
// ancestor, so every bare `git` call there fails. The resolvers must reach the
// shared object store through `jj git root` instead.
//
// The fixture is a throwaway colocated repository, never the one under test.
// Its secondary workspace commits something the primary never checks out, so
// an answer that fell back to the shared git HEAD would name the wrong commit
// and fail here rather than pass by coincidence.
func TestResolversWorkInASecondaryWorkspaceOutsideThePrimary(t *testing.T) {
	if _, err := exec.LookPath("jj"); err != nil {
		t.Skip("jj not installed")
	}
	script, err := filepath.Abs(filepath.Join("..", "..", "scripts", "current-commit.sh"))
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	primary := filepath.Join(base, "primary")
	secondary := filepath.Join(base, "secondary")
	jj := func(dir string, args ...string) string {
		t.Helper()
		cmd := exec.Command("jj", append([]string{
			"--config", "user.name=ratchet-test",
			"--config", "user.email=ratchet-test@example.invalid",
		}, args...)...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("jj %s (in %s): %v\n%s", strings.Join(args, " "), dir, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	jj(base, "git", "init", "--colocate", primary)
	if err := os.WriteFile(filepath.Join(primary, "a.txt"), []byte("a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	jj(primary, "commit", "-m", "a")
	jj(primary, "workspace", "add", "--name", "secondary", secondary)
	if err := os.WriteFile(filepath.Join(secondary, "b.txt"), []byte("b\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	jj(secondary, "commit", "-m", "b")
	want := jj(secondary, "log", "--no-graph", "-r", "@-", "-T", "commit_id.short(12)")

	if _, err := os.Stat(filepath.Join(secondary, ".git")); err == nil {
		t.Fatal("fixture is wrong: the secondary workspace has a .git, so it does not reproduce the layout")
	}
	sharedHead, _ := exec.Command("git", "--git-dir="+filepath.Join(primary, ".git"),
		"rev-parse", "--short=12", "HEAD").Output()
	if strings.TrimSpace(string(sharedHead)) == want {
		t.Fatal("fixture is wrong: the shared git HEAD already names the secondary's commit, so a HEAD fallback would pass")
	}

	t.Chdir(secondary)
	if got := jjWorkspaceSHA(); got != want {
		t.Errorf("jjWorkspaceSHA() = %q, want the secondary workspace's commit %q", got, want)
	}
	if !gitHasCommit(want) {
		t.Errorf("gitHasCommit(%q) = false from the secondary workspace", want)
	}
	// The fixture has no remotes, so the commit is unpublished. Reporting it as
	// published would mean the reachability query never reached the object
	// store, since that is the answer given when git cannot be asked.
	if capturedSHAIsPublished(want) {
		t.Errorf("capturedSHAIsPublished(%q) = true, want false: the fixture has no remote-tracking refs", want)
	}

	cmd := exec.Command("sh", script)
	cmd.Dir = secondary
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("current-commit.sh: %v", err)
	}
	if got := strings.TrimSpace(string(out)); got != want {
		t.Errorf("current-commit.sh printed %q in the secondary workspace, want %q", got, want)
	}
}

// jjOut runs a jj log command and returns its trimmed output.
func jjOut(t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.Command("jj", append([]string{"log", "--no-graph"}, args...)...).Output()
	if err != nil {
		t.Fatalf("jj log %v: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}
