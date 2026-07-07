/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package wasm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveShell(t *testing.T) {
	if ct, x, err := ResolveShell("xterm"); err != nil || ct != "" || !x {
		t.Fatalf("xterm: got (%q,%v,%v)", ct, x, err)
	}
	if ct, x, err := ResolveShell("none"); err != nil || ct != "" || x {
		t.Fatalf("none: got (%q,%v,%v)", ct, x, err)
	}
	if _, _, err := ResolveShell("/no/such/template.html"); err == nil {
		t.Fatal("missing template should error")
	}

	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.html")
	if err := os.WriteFile(bad, []byte("<html>no marker</html>"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ResolveShell(bad); err == nil {
		t.Fatal("template without marker should error")
	}

	good := filepath.Join(dir, "good.html")
	if err := os.WriteFile(good, []byte("<script>"+HostBodyMarker+"</script>"), 0644); err != nil {
		t.Fatal(err)
	}
	if ct, x, err := ResolveShell(good); err != nil || ct != good || x {
		t.Fatalf("good: got (%q,%v,%v)", ct, x, err)
	}
}
