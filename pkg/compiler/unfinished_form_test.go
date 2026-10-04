/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// unfinishedSources each end inside an open form after zero or more
// complete forms.
var unfinishedSources = []string{
	"(def incomplete",
	"(+ 1 2) (def x",
	"(+ 1 2)\n(defn f [x]\n  (inc x)",
	"[1 2",
	"{:a 1",
	"#{1 2",
	"\"abc",
	"(+ 1 2) '(",
}

// EOF inside an unfinished form is a syntax error, not the end of input
// (#807). The error still counts as EOF so frontends can keep reading lines.
func TestCompileMultipleRejectsUnfinishedForm(t *testing.T) {
	for _, src := range unfinishedSources {
		_, err := Eval(src)
		if err == nil {
			t.Errorf("Eval(%q) succeeded, want a syntax error", src)
			continue
		}
		if !IsErrorEOF(err) {
			t.Errorf("Eval(%q) error is not EOF-caused: %v", src, err)
		}
	}
}

func TestLoadStringRejectsUnfinishedForm(t *testing.T) {
	if _, err := Eval(`(load-string "(+ 1 2) (def x")`); err == nil {
		t.Fatal("load-string of an unfinished form succeeded, want a syntax error")
	}
}

// EOF at a form boundary, after whitespace or no-value forms, still ends the
// input cleanly.
func TestCompileMultipleAcceptsCleanEOF(t *testing.T) {
	for _, src := range []string{
		"",
		"   \n\t",
		"(+ 1 2)",
		"(+ 1 2)  \n",
		"(+ 1 2) ;; done",
		"#!/usr/bin/env lg\n(+ 1 2)",
		"#!/usr/bin/env lg",
		"+",
		"-",
	} {
		if _, err := Eval(src); err != nil {
			t.Errorf("Eval(%q): %v", src, err)
		}
	}
}

// The reader looks past a leading + or - for a digit. At end of input the
// sign is a complete symbol, not an unfinished form.
func TestSignSymbolAtEOF(t *testing.T) {
	for _, src := range []string{"42 #_ +", "42 #_ -"} {
		v, err := Eval(src)
		if err != nil {
			t.Fatalf("Eval(%q): %v", src, err)
		}
		if n, ok := v.(vm.Int); !ok || int(n) != 42 {
			t.Fatalf("Eval(%q) = %v, want 42", src, v)
		}
	}
	for _, src := range []string{"+", "-"} {
		v, err := ReadDataString(src)
		if err != nil {
			t.Fatalf("ReadDataString(%q): %v", src, err)
		}
		if v != vm.Symbol(src) {
			t.Fatalf("ReadDataString(%q) = %v, want the symbol %s", src, v, src)
		}
	}
}

func TestSplitTopLevelFormsRejectsUnfinishedForm(t *testing.T) {
	for _, src := range unfinishedSources {
		forms, err := SplitTopLevelForms(src, "<test>")
		if err == nil {
			t.Errorf("SplitTopLevelForms(%q) = %d forms, want a syntax error", src, len(forms))
		}
	}
	for src, want := range map[string]int{
		"(+ 1 2) ;; done\n[3]  ": 2,
		"42 #_ +":                1,
		"42 #_ -":                1,
		"42 +":                   2,
	} {
		forms, err := SplitTopLevelForms(src, "<test>")
		if err != nil {
			t.Errorf("SplitTopLevelForms(%q): %v", src, err)
			continue
		}
		if len(forms) != want {
			t.Errorf("SplitTopLevelForms(%q) = %d forms, want %d", src, len(forms), want)
		}
	}
}
