// Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLintGoSpans(t *testing.T) {
	const source = `package fixture

func First(
    x string,
) string {
    _ = "{"
    return x
}

func Second() {
    // }
}

func Asm()
`
	path := filepath.Join(t.TempDir(), "fixture.go")
	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := spans([]string{path})
	if err != nil {
		t.Fatal(err)
	}
	want := []span{{path, 3, 8}, {path, 10, 12}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("spans = %#v, want %#v", got, want)
	}

	cmd := exec.Command("go", "run", "lint_go_spans.go", "--", path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("helper CLI: %v\n%s", err, out)
	}
	wantOut := strings.Join([]string{
		path + "\t3\t8",
		path + "\t10\t12",
		"",
	}, "\n")
	if string(out) != wantOut {
		t.Errorf("CLI output = %q, want %q", out, wantOut)
	}
}
