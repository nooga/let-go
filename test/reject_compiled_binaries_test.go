package test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func binaryGuardGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func binaryGuardCommit(t *testing.T, dir, message string) string {
	t.Helper()
	binaryGuardGit(t, dir, "add", "-A")
	binaryGuardGit(t, dir, "commit", "-q", "-m", message)
	return binaryGuardGit(t, dir, "rev-parse", "HEAD")
}

func binaryGuardNoRefEnv() []string {
	env := make([]string, 0, len(os.Environ()))
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "PRE_COMMIT_FROM_REF=") &&
			!strings.HasPrefix(entry, "PRE_COMMIT_TO_REF=") {
			env = append(env, entry)
		}
	}
	return env
}

func binaryGuardEnv(fromRef, toRef string) []string {
	return append(binaryGuardNoRefEnv(),
		"PRE_COMMIT_FROM_REF="+fromRef,
		"PRE_COMMIT_TO_REF="+toRef,
	)
}

func binaryGuardToOnlyEnv(toRef string) []string {
	return append(binaryGuardNoRefEnv(), "PRE_COMMIT_TO_REF="+toRef)
}

func TestRejectCompiledBinariesChecksIntermediatePushBlobs(t *testing.T) {
	script, err := filepath.Abs("../scripts/reject_compiled_binaries.py")
	if err != nil {
		t.Fatal(err)
	}

	repo := t.TempDir()
	binaryGuardGit(t, repo, "init", "-q")
	binaryGuardGit(t, repo, "config", "user.name", "Binary Guard Test")
	binaryGuardGit(t, repo, "config", "user.email", "binary-guard@example.invalid")

	if err := os.WriteFile(filepath.Join(repo, "tracked.txt"), []byte("safe\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fromRef := binaryGuardCommit(t, repo, "base")

	artifact := filepath.Join(repo, "accidental-build")
	if err := os.WriteFile(artifact, []byte("\x7fELFpayload"), 0o755); err != nil {
		t.Fatal(err)
	}
	binaryGuardCommit(t, repo, "add accidental executable")
	if err := os.Remove(artifact); err != nil {
		t.Fatal(err)
	}
	toRef := binaryGuardCommit(t, repo, "delete accidental executable")

	if diff := binaryGuardGit(t, repo, "diff", "--name-only", fromRef+"..."+toRef); diff != "" {
		t.Fatalf("test precondition: expected empty net diff, got %q", diff)
	}

	for _, args := range [][]string{nil, {"tracked.txt"}} {
		cmd := exec.Command("python3", append([]string{script}, args...)...)
		cmd.Dir = repo
		cmd.Env = binaryGuardEnv(fromRef, toRef)
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("pre-push guard accepted an ELF blob from an intermediate commit (args=%v)\n%s", args, out)
		}
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("run guard: %v", err)
		}
		if !strings.Contains(string(out), "ELF executable") ||
			!strings.Contains(string(out), "accidental-build") {
			t.Fatalf("guard failed without identifying the ELF blob and path:\n%s", out)
		}
	}
}

func TestRejectCompiledBinariesChecksOrphanHistoryWithoutFromRef(t *testing.T) {
	script, err := filepath.Abs("../scripts/reject_compiled_binaries.py")
	if err != nil {
		t.Fatal(err)
	}

	repo := t.TempDir()
	binaryGuardGit(t, repo, "init", "-q")
	binaryGuardGit(t, repo, "config", "user.name", "Binary Guard Test")
	binaryGuardGit(t, repo, "config", "user.email", "binary-guard@example.invalid")
	if err := os.WriteFile(filepath.Join(repo, "tracked.txt"), []byte("safe\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	binaryGuardCommit(t, repo, "root")

	artifact := filepath.Join(repo, "orphan-build")
	if err := os.WriteFile(artifact, []byte("\x7fELFpayload"), 0o755); err != nil {
		t.Fatal(err)
	}
	binaryGuardCommit(t, repo, "add accidental executable")
	if err := os.Remove(artifact); err != nil {
		t.Fatal(err)
	}
	toRef := binaryGuardCommit(t, repo, "delete accidental executable")

	cmd := exec.Command("python3", script)
	cmd.Dir = repo
	cmd.Env = binaryGuardToOnlyEnv(toRef)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("pre-push guard accepted orphan history without PRE_COMMIT_FROM_REF\n%s", out)
	}
	if !strings.Contains(string(out), "ELF executable") ||
		!strings.Contains(string(out), "orphan-build") {
		t.Fatalf("guard failed without identifying the orphan-history blob and path:\n%s", out)
	}
}

func TestRejectCompiledBinariesPreservesExecutableMagicChecks(t *testing.T) {
	script, err := filepath.Abs("../scripts/reject_compiled_binaries.py")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name  string
		magic []byte
		label string
	}{
		{"elf", []byte("\x7fELF"), "ELF executable"},
		{"macho-32-be", []byte("\xfe\xed\xfa\xce"), "Mach-O executable (32-bit)"},
		{"macho-64-be", []byte("\xfe\xed\xfa\xcf"), "Mach-O executable (64-bit)"},
		{"macho-32-le", []byte("\xce\xfa\xed\xfe"), "Mach-O executable (32-bit, LE)"},
		{"macho-64-le", []byte("\xcf\xfa\xed\xfe"), "Mach-O executable (64-bit, LE)"},
		{"macho-universal-be", []byte("\xca\xfe\xba\xbe"), "Mach-O universal binary"},
		{"macho-universal-le", []byte("\xbe\xba\xfe\xca"), "Mach-O universal binary (LE)"},
		{"pe", []byte("MZxx"), "PE/DOS executable"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			artifact := filepath.Join(dir, "artifact")
			if err := os.WriteFile(artifact, append(tc.magic, []byte("payload")...), 0o644); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command("python3", script, "artifact")
			cmd.Dir = dir
			cmd.Env = binaryGuardNoRefEnv()
			out, err := cmd.CombinedOutput()
			if err == nil || !strings.Contains(string(out), tc.label) {
				t.Fatalf("guard did not reject %s: err=%v\n%s", tc.label, err, out)
			}
		})
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "core.lgb"), []byte("LGB\x01payload"), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("python3", script, "core.lgb")
	cmd.Dir = dir
	cmd.Env = binaryGuardNoRefEnv()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("guard rejected legitimate LGB data: %v\n%s", err, out)
	}
}

func TestRejectCompiledBinariesAllowsSafePushedBlobs(t *testing.T) {
	script, err := filepath.Abs("../scripts/reject_compiled_binaries.py")
	if err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	binaryGuardGit(t, repo, "init", "-q")
	binaryGuardGit(t, repo, "config", "user.name", "Binary Guard Test")
	binaryGuardGit(t, repo, "config", "user.email", "binary-guard@example.invalid")
	if err := os.WriteFile(filepath.Join(repo, "source.go"), []byte("package safe\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fromRef := binaryGuardCommit(t, repo, "base")
	if err := os.WriteFile(filepath.Join(repo, "source.go"), []byte("package safer\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	toRef := binaryGuardCommit(t, repo, "safe change")

	cmd := exec.Command("python3", script, "source.go")
	cmd.Dir = repo
	cmd.Env = binaryGuardEnv(fromRef, toRef)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("guard rejected safe pushed blobs: %v\n%s", err, out)
	}
}

func TestRejectCompiledBinariesFailsClosedOnBadRefs(t *testing.T) {
	script, err := filepath.Abs("../scripts/reject_compiled_binaries.py")
	if err != nil {
		t.Fatal(err)
	}
	repo := t.TempDir()
	binaryGuardGit(t, repo, "init", "-q")
	binaryGuardGit(t, repo, "config", "user.name", "Binary Guard Test")
	binaryGuardGit(t, repo, "config", "user.email", "binary-guard@example.invalid")
	if err := os.WriteFile(filepath.Join(repo, "safe"), []byte("safe\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	toRef := binaryGuardCommit(t, repo, "base")

	cmd := exec.Command("python3", script)
	cmd.Dir = repo
	cmd.Env = binaryGuardEnv("not-a-ref", toRef)
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "reject_compiled_binaries: git rev-list") {
		t.Fatalf("guard did not fail closed for an invalid from_ref: err=%v\n%s", err, out)
	}
}

func TestRejectCompiledBinariesHookAlwaysRuns(t *testing.T) {
	config, err := os.ReadFile("../.pre-commit-config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	text := string(config)
	start := strings.Index(text, "- id: forbid-compiled-binaries")
	if start < 0 {
		t.Fatal("forbid-compiled-binaries hook is missing")
	}
	section := text[start:]
	if next := strings.Index(section[len("- id: forbid-compiled-binaries"):], "\n      - id:"); next >= 0 {
		section = section[:len("- id: forbid-compiled-binaries")+next]
	}
	if !strings.Contains(section, "always_run: true") ||
		!strings.Contains(section, "stages: [pre-commit, pre-push]") {
		t.Fatalf("binary guard must always run at pre-push even when the net diff is empty:\n%s", section)
	}
}

func TestPrePushGoTestTimeoutMatchesMakefile(t *testing.T) {
	makefile, err := os.ReadFile("../Makefile")
	if err != nil {
		t.Fatal(err)
	}
	config, err := os.ReadFile("../.pre-commit-config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	makeMatch := regexp.MustCompile(`(?m)^GO-TEST-TIMEOUT \?= (\S+)$`).FindSubmatch(makefile)
	hookMatch := regexp.MustCompile(`entry: go test -short -timeout (\S+) ./\.\.\.`).FindSubmatch(config)
	if len(makeMatch) != 2 || len(hookMatch) != 2 {
		t.Fatalf("could not find Go test timeout in Makefile or pre-push hook")
	}
	if string(makeMatch[1]) != string(hookMatch[1]) {
		t.Fatalf("pre-push Go timeout %s differs from Makefile timeout %s",
			hookMatch[1], makeMatch[1])
	}
}
