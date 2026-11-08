/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package wasm

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// storageShimDriver runs the worker code from lg-host-core.js under node,
// sends it an init carrying a storage snapshot, and stops it at the
// wasm_exec.js eval, after _lgStorage is built and before any wasm loads.
const storageShimDriver = `
globalThis.posted = [];
globalThis.postMessage = (m) => posted.push(m);
%WORKER%
(async () => {
  try {
    await onmessage({data: {
      t: 'init',
      sab: new SharedArrayBuffer(128),
      storage: {'let-go:s:a': '1', 'let-go:s:b': '2'},
      wasmExecJS: 'throw new Error("stop before wasm")',
    }});
  } catch (e) {
    if (e.message !== 'stop before wasm') throw e;
  }
  const s = globalThis._lgStorage;
  const keys = () => { const ks = []; for (let i = 0; i <= s.length; i++) ks.push(s.key(i)); return ks; };
  const out = {};
  out.initial = keys();
  out.missing = s.getItem('let-go:s:none');
  s.setItem('let-go:s:c', 3);
  out.afterAdd = keys();
  out.added = s.getItem('let-go:s:c');
  s.setItem('let-go:s:a', 'x');
  out.afterOverwrite = keys();
  s.removeItem('let-go:s:b');
  out.afterRemove = keys();
  out.removed = s.getItem('let-go:s:b');
  s.removeItem('let-go:s:none');
  out.posted = posted;
  console.log(JSON.stringify(out));
})().catch((e) => { console.error(e); process.exit(1); });
`

func TestWorkerStorageShim(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node not on PATH; cannot run the worker storage shim")
	}

	const startMark, endMark = "const workerCode = `", "\n  `;"
	start := strings.Index(lgHostCoreJS, startMark)
	if start < 0 {
		t.Fatalf("lg-host-core.js: %q not found", startMark)
	}
	start += len(startMark)
	end := strings.Index(lgHostCoreJS[start:], endMark)
	if end < 0 {
		t.Fatalf("lg-host-core.js: end of workerCode not found")
	}
	worker := lgHostCoreJS[start : start+end]
	if !strings.Contains(worker, "globalThis._lgStorage") {
		t.Fatal("workerCode no longer defines globalThis._lgStorage")
	}

	driver := filepath.Join(t.TempDir(), "shim.js")
	if err := os.WriteFile(driver, []byte(strings.Replace(storageShimDriver, "%WORKER%", worker, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, node, driver).CombinedOutput()
	if err != nil {
		t.Fatalf("node: %v\n%s", err, out)
	}

	var got struct {
		Initial, AfterAdd, AfterOverwrite, AfterRemove []*string
		Missing, Added, Removed                        *string
		Posted                                         []map[string]any
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &got); err != nil {
		t.Fatalf("parse node output: %v\n%s", err, out)
	}

	// key(length) is included in each list, so a trailing nil checks that key
	// past the end reads null.
	keys := func(ks ...string) []*string {
		ps := make([]*string, 0, len(ks)+1)
		for _, k := range ks {
			ps = append(ps, &k)
		}
		return append(ps, nil)
	}
	str := func(p *string) any {
		if p == nil {
			return nil
		}
		return *p
	}
	strs := func(ps []*string) []any {
		vs := make([]any, len(ps))
		for i, p := range ps {
			vs[i] = str(p)
		}
		return vs
	}
	for _, c := range []struct {
		name      string
		got, want []*string
	}{
		{"initial", got.Initial, keys("let-go:s:a", "let-go:s:b")},
		{"after add", got.AfterAdd, keys("let-go:s:a", "let-go:s:b", "let-go:s:c")},
		{"after overwrite", got.AfterOverwrite, keys("let-go:s:a", "let-go:s:b", "let-go:s:c")},
		{"after remove", got.AfterRemove, keys("let-go:s:a", "let-go:s:c")},
	} {
		if !reflect.DeepEqual(strs(c.got), strs(c.want)) {
			t.Errorf("%s: keys %v, want %v", c.name, strs(c.got), strs(c.want))
		}
	}
	if got.Missing != nil {
		t.Errorf("getItem of a missing key = %q, want null", *got.Missing)
	}
	if str(got.Added) != "3" {
		t.Errorf("getItem after setItem(k, 3) = %v, want \"3\"", str(got.Added))
	}
	if got.Removed != nil {
		t.Errorf("getItem after removeItem = %q, want null", *got.Removed)
	}

	wantPosted := []map[string]any{
		{"t": "storage-set", "k": "let-go:s:c", "v": "3"},
		{"t": "storage-set", "k": "let-go:s:a", "v": "x"},
		{"t": "storage-remove", "k": "let-go:s:b"},
		{"t": "storage-remove", "k": "let-go:s:none"},
	}
	if !reflect.DeepEqual(got.Posted, wantPosted) {
		t.Errorf("posted %v, want %v", got.Posted, wantPosted)
	}
}
