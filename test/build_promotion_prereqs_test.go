/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBuildPromotionTracksGateInputs(t *testing.T) {
	t.Setenv("MAKEFLAGS", "")
	t.Setenv("MFLAGS", "")
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	buildDir := filepath.Join(tmp, "build")
	binDir := filepath.Join(tmp, "bin")
	candidate := filepath.Join(buildDir, "lg")
	bootprobe := filepath.Join(buildDir, "bootprobe")
	promoted := filepath.Join(binDir, "lg")
	for _, path := range []string{candidate, bootprobe, promoted} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("fixture"), 0o755); err != nil {
			t.Fatal(err)
		}
		future := time.Now().Add(time.Hour)
		if err := os.Chtimes(path, future, future); err != nil {
			t.Fatal(err)
		}
	}

	cases := []struct {
		changed string
		wants   []string
	}{
		{
			changed: "cmd/bootprobe/main.go",
			wants: []string{
				"go build -o " + bootprobe + " ./cmd/bootprobe",
				"scripts/smoke.lg",
				"scripts/smoke-boot.sh",
				"install -m 0755 " + candidate + " " + promoted,
			},
		},
		{
			changed: "scripts/smoke.lg",
			wants: []string{
				"scripts/smoke.lg",
				"scripts/smoke-boot.sh",
				"install -m 0755 " + candidate + " " + promoted,
			},
		},
		{
			changed: "scripts/smoke-boot.sh",
			wants: []string{
				"scripts/smoke.lg",
				"scripts/smoke-boot.sh",
				"install -m 0755 " + candidate + " " + promoted,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.changed, func(t *testing.T) {
			args := []string{
				"-n", "--no-print-directory",
				"BUILD-DIR=" + buildDir,
				"BIN-DIR=" + binDir,
				"-W", tc.changed,
				"build",
			}
			cmd := exec.Command("make", args...)
			cmd.Dir = root
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("make %v: %v\n%s", args, err, out)
			}
			output := string(out)
			for _, want := range tc.wants {
				if !strings.Contains(output, want) {
					t.Errorf("make -W %s did not rerun %q; output:\n%s", tc.changed, want, output)
				}
			}
		})
	}
}
