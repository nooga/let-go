package vm

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStderr redirects os.Stderr for the duration of fn and returns whatever
// was written. Namespace.Alias emits its collision warning via
// fmt.Fprintf(os.Stderr, ...), so this is the only seam to observe it.
func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stderr = w
	defer func() { os.Stderr = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("closing pipe: %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("reading captured stderr: %v", err)
	}
	return buf.String()
}

// TestAlias_SameNameReloadIsSilent is the regression guard for the collision
// guard comparing namespace NAMES rather than pointers.
//
// The guard used to compare the previous and new *Namespace by pointer. But two
// loads of the same namespace present distinct *Namespace objects with the same
// name (e.g. a benchmark re-running the whole corpus per iteration re-creates
// the ns), so pointer comparison spuriously warned
//
//	"alias X already points to lib, being repointed to lib"
//
// — note the SAME name on both sides — for what is really the idempotent-reload
// case. Because `go test` merges the test binary's stderr into stdout, that spam
// landed between a benchmark's name and its metrics line and made benchfmt drop
// the [bytecode] suite results entirely. Sameness must be judged by name.
func TestAlias_SameNameReloadIsSilent(t *testing.T) {
	consumer := NewNamespace("consumer")

	// Two distinct namespace objects for the SAME logical namespace — exactly
	// what a reload produces.
	lib1 := NewNamespace("mylib")
	lib2 := NewNamespace("mylib")
	if lib1 == lib2 {
		t.Fatal("test precondition: the two loads must be distinct pointers")
	}

	out := captureStderr(t, func() {
		consumer.Alias(Symbol("ml"), lib1) // first alias: always silent
		consumer.Alias(Symbol("ml"), lib2) // same-name reload: MUST be silent
	})
	if out != "" {
		t.Fatalf("re-aliasing to a same-named namespace (reload) must be silent, "+
			"but Alias warned: %q", out)
	}
}

// TestAlias_DifferentNameCollisionWarns is the positive control: the guard must
// still fire for a genuine collision (same short name repointed at a
// differently-named namespace), so the name-based comparison didn't just
// disable the warning.
func TestAlias_DifferentNameCollisionWarns(t *testing.T) {
	consumer := NewNamespace("consumer")
	first := NewNamespace("lib.one")
	second := NewNamespace("lib.two")

	out := captureStderr(t, func() {
		consumer.Alias(Symbol("p"), first)  // first alias: silent
		consumer.Alias(Symbol("p"), second) // real collision: MUST warn
	})
	if !strings.Contains(out, "already points to") {
		t.Fatalf("re-aliasing to a different-named namespace is a real collision "+
			"and must warn, but got: %q", out)
	}
	// The warning should name both targets so the conflict is diagnosable.
	if !strings.Contains(out, "lib.one") || !strings.Contains(out, "lib.two") {
		t.Fatalf("collision warning should name both namespaces, got: %q", out)
	}
}
