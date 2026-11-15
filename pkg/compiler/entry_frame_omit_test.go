/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

// compileEntryProgram compiles a program whose -main counts its runs, ending
// in the given top-level forms and then the count. Compiling evaluates each
// top-level form, so the result is how often -main ran at build time, which
// is how often the bundle would run it before the native frame calls it.
func compileEntryProgram(t *testing.T, i int, tail string, entry bool) (vm.Value, error) {
	t.Helper()
	ns := fmt.Sprintf("entryomit%d", i)
	src := fmt.Sprintf("(ns %s)\n(def runs (atom 0))\n(defn -main [& args] (swap! runs inc))\n%s\n@runs\n", ns, tail)
	ctx := NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS))
	var (
		res vm.Value
		err error
	)
	if entry {
		_, res, err = ctx.CompileMultipleForEntryFrame(strings.NewReader(src), ns, "-main")
	} else {
		_, res, err = ctx.CompileMultiple(strings.NewReader(src))
	}
	return res, err
}

func TestEntryFrameOmitsEntryCalls(t *testing.T) {
	for i, tail := range []string{
		`(-main)`,
		`(apply -main *command-line-args*)`,
		`(apply -main ["a" "b"])`,
		`(clojure.core/apply -main [])`,
		`(do (def before 1) (apply -main []))`,
		`(when-not *compiling-aot* (apply -main *command-line-args*))`,
		// A call in a body that runs only when called isn't a top-level
		// call, so it stays as written.
		`(defn again [] (-main))`,
		`(def later (fn [] (apply -main [])))`,
		`(def boot (delay (apply -main *command-line-args*)))`,
		"(defmulti cmd identity)\n(defmethod cmd :restart [_] (-main))",
		"(defmacro defcommand [n & body] `(def ~n (fn [] ~@body)))\n(defcommand restart (-main))",
		`(letfn [(again [] (-main))] nil)`,
		`(def raw (fn* [] (-main)))`,
		// case test constants are data, not calls.
		`(case (first *command-line-args*) (-main main) 1 2)`,
		`(def entry -main)`,
		`'(-main)`,
	} {
		res, err := compileEntryProgram(t, 100+i, tail, true)
		if err != nil {
			t.Errorf("%s: %v", tail, err)
			continue
		}
		if res != vm.MakeInt(0) {
			t.Errorf("%s: -main ran %v times at build time, want 0", tail, res)
		}
	}
}

func TestEntryFrameOmitsQualifiedApply(t *testing.T) {
	res, err := compileEntryProgram(t, 200, `(apply entryomit200/-main [])`, true)
	if err != nil || res != vm.MakeInt(0) {
		t.Fatalf("got %v, %v; want 0 runs", res, err)
	}
}

// Without an entry frame, every spelling still runs -main as written.
func TestOrdinaryCompileKeepsEntryCalls(t *testing.T) {
	for i, tail := range []string{`(apply -main [])`, `(when true (-main))`} {
		res, err := compileEntryProgram(t, 400+i, tail, false)
		if err != nil || res != vm.MakeInt(1) {
			t.Errorf("%s: got %v, %v; want 1 run", tail, res, err)
		}
	}
}
