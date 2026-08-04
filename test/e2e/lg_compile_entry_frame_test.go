/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLgCompileEntryFrameOptIn pins the #628 ownership contract: frame
// emission is opt-in via --entry-frame. Without the flag, lg-compile writes
// only lowered packages (historical behavior for Gloat and other
// orchestrators). With the flag, exactly one package-main frame is emitted,
// plus a private-entry bridge when applicable. Package lowering is identical
// in both modes (excluding the intentionally added bridge file).
func TestLgCompileEntryFrameOptIn(t *testing.T) {
	bin := buildLG(t)
	root := repoRoot(t)

	writeSrc := func(t *testing.T, body string) string {
		t.Helper()
		dir := t.TempDir()
		p := filepath.Join(dir, "app.lg")
		if err := os.WriteFile(p, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
		return p
	}

	runCompile := func(t *testing.T, outDir, src string, entryFrame bool) string {
		t.Helper()
		args := []string{"scripts/lg-compile"}
		if entryFrame {
			args = append(args, "--entry-frame")
		}
		args = append(args, outDir, "tmpmod", src)
		cmd := exec.Command(bin, args...)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "LG_SOURCE_PATHS=pkg/rt/core")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("lg-compile (entry-frame=%v): %v\n%s", entryFrame, err, out)
		}
		return string(out)
	}

	pkgRel := filepath.Join("app", "app.go")
	bridgeRel := filepath.Join("app", "native_entry.go")

	mustAbsent := func(t *testing.T, path string) {
		t.Helper()
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("%s should be absent; err=%v", path, err)
		}
	}
	mustPresent := func(t *testing.T, path string) []byte {
		t.Helper()
		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s missing: %v", path, err)
		}
		return b
	}

	t.Run("no flag + private entry → no main.go, no bridge", func(t *testing.T) {
		src := writeSrc(t, "(ns app)\n(defn- -main [] 1)\n")
		outDir := t.TempDir()
		log := runCompile(t, outDir, src, false)
		mustAbsent(t, filepath.Join(outDir, "main.go"))
		mustAbsent(t, filepath.Join(outDir, bridgeRel))
		mustPresent(t, filepath.Join(outDir, pkgRel))
		if strings.Contains(log, "native-entry frame ->") || strings.Contains(log, "native-entry bridge ->") {
			t.Fatalf("must not report frame/bridge emission without --entry-frame; log:\n%s", log)
		}
	})

	t.Run("flag + public entry → one main.go, no bridge", func(t *testing.T) {
		src := writeSrc(t, "(ns app)\n(defn -main [] 1)\n")
		outDir := t.TempDir()
		log := runCompile(t, outDir, src, true)
		mainBody := string(mustPresent(t, filepath.Join(outDir, "main.go")))
		mustAbsent(t, filepath.Join(outDir, bridgeRel))
		if c := strings.Count(mainBody, "func main()"); c != 1 {
			t.Fatalf("want exactly one func main(), got %d", c)
		}
		if !strings.Contains(log, "native-entry frame ->") {
			t.Fatalf("expected frame emission log; got:\n%s", log)
		}
		if strings.Contains(log, "native-entry bridge ->") {
			t.Fatalf("public entry must not emit a bridge; log:\n%s", log)
		}
	})

	t.Run("flag + private entry → one main.go, one bridge", func(t *testing.T) {
		src := writeSrc(t, "(ns app)\n(defn- -main [] 1)\n")
		outDir := t.TempDir()
		log := runCompile(t, outDir, src, true)
		mainBody := string(mustPresent(t, filepath.Join(outDir, "main.go")))
		bridgeBody := string(mustPresent(t, filepath.Join(outDir, bridgeRel)))
		if c := strings.Count(mainBody, "func main()"); c != 1 {
			t.Fatalf("want exactly one func main(), got %d", c)
		}
		if c := strings.Count(bridgeBody, "func NativeMain("); c != 1 {
			t.Fatalf("want exactly one NativeMain, got %d\n%s", c, bridgeBody)
		}
		if !strings.Contains(log, "native-entry frame ->") || !strings.Contains(log, "native-entry bridge ->") {
			t.Fatalf("expected frame + bridge emission logs; got:\n%s", log)
		}
	})

	t.Run("package lowering identical excluding bridge", func(t *testing.T) {
		src := writeSrc(t, "(ns app)\n(defn- -main [] 1)\n")
		outA := t.TempDir()
		outB := t.TempDir()
		runCompile(t, outA, src, false)
		runCompile(t, outB, src, true)
		a := mustPresent(t, filepath.Join(outA, pkgRel))
		b := mustPresent(t, filepath.Join(outB, pkgRel))
		if string(a) != string(b) {
			t.Fatalf("lowered package differed with/without --entry-frame\n--- no flag ---\n%s\n--- flag ---\n%s", a, b)
		}
		mustAbsent(t, filepath.Join(outA, bridgeRel))
		mustPresent(t, filepath.Join(outB, bridgeRel))
	})
}
