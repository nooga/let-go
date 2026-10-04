/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"testing"
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
	} {
		if _, err := Eval(src); err != nil {
			t.Errorf("Eval(%q): %v", src, err)
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
	forms, err := SplitTopLevelForms("(+ 1 2) ;; done\n[3]  ", "<test>")
	if err != nil {
		t.Fatalf("SplitTopLevelForms: %v", err)
	}
	if len(forms) != 2 {
		t.Fatalf("SplitTopLevelForms = %d forms, want 2", len(forms))
	}
}
