/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package gofragment

import (
	"strconv"
	"strings"
	"testing"
)

const generatedGoProbe = `package probe

var Result = x + y*2

func Steps() {
	n := 1
	if n > 0 {
		n++
	}
}

func Compute(x int) int {
	return x * 2
}

func Nested(x int) int {
	trace(x)
	if x > 0 {
		doubled := x * 2
		log(doubled)
		return doubled
	}
	return 0
}

func Duplicate(x int, left, right bool) {
	a := x * 2
	b := x * 2
	if left {
		n := x + 1
		log(n)
	}
	if right {
		n := x + 1
		log(n)
	}
	_, _ = a, b
}

type Item struct {
	Name string
}
`

func TestGoFragmentMatcherFindsSubtreesInsideStableFunctionScope(t *testing.T) {
	tests := []struct {
		name string
		req  GoMatchRequest
	}{
		{
			"expression in nested block",
			GoMatchRequest{Kind: GoExpression, Expected: "x * 2", Target: "Nested"},
		},
		{
			"consecutive statements in nested larger block",
			GoMatchRequest{
				Kind:     GoStatements,
				Expected: "doubled := x * 2; log(doubled)",
				Target:   "Nested",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MatchGeneratedGoFragment(tt.req, generatedGoProbe); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGoFragmentMatcherSubtreeCardinalityRejectsZeroAndDuplicates(t *testing.T) {
	tests := []struct {
		name      string
		req       GoMatchRequest
		wantCount int
	}{
		{
			"missing expression",
			GoMatchRequest{Kind: GoExpression, Expected: "x * 3", Target: "Nested"},
			0,
		},
		{
			"duplicate expression",
			GoMatchRequest{Kind: GoExpression, Expected: "x * 2", Target: "Duplicate"},
			2,
		},
		{
			"missing statement sequence",
			GoMatchRequest{Kind: GoStatements, Expected: "n := x + 2; log(n)", Target: "Duplicate"},
			0,
		},
		{
			"duplicate statement sequence",
			GoMatchRequest{Kind: GoStatements, Expected: "n := x + 1; log(n)", Target: "Duplicate"},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MatchGeneratedGoFragment(tt.req, generatedGoProbe)
			want := string(tt.req.Kind) + ` target "` + tt.req.Target + `": found ` +
				strconv.Itoa(tt.wantCount) + ", want exactly 1"
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("cardinality error = %v, want substring %q", err, want)
			}
		})
	}
}

func TestGoFragmentMatcherMatchesParserWrappersPointwise(t *testing.T) {
	tests := []struct {
		name string
		req  GoMatchRequest
	}{
		{"expression", GoMatchRequest{Kind: GoExpression, Expected: "x + y*2", Target: "Result"}},
		{"statements", GoMatchRequest{Kind: GoStatements, Expected: "n := 1; if n > 0 { n++ }", Target: "Steps"}},
		{"function", GoMatchRequest{Kind: GoFunction, Expected: "func Compute(x int) int { return x * 2 }", Target: "Compute"}},
		{"declaration", GoMatchRequest{Kind: GoDeclaration, Expected: "type Item struct { Name string }", Target: "Item"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MatchGeneratedGoFragment(tt.req, generatedGoProbe); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGoFragmentMatcherDefaultsToExactlyOneTarget(t *testing.T) {
	missing := GoMatchRequest{Kind: GoFunction, Expected: "func Missing() {}", Target: "Missing"}
	if err := MatchGeneratedGoFragment(missing, generatedGoProbe); err == nil ||
		!strings.Contains(err.Error(), `function target "Missing": found 0, want exactly 1`) {
		t.Fatalf("missing target error = %v", err)
	}

	duplicateSource := generatedGoProbe + "\nfunc Compute() {}\n"
	duplicate := GoMatchRequest{Kind: GoFunction, Expected: "func Compute(x int) int { return x * 2 }", Target: "Compute"}
	if err := MatchGeneratedGoFragment(duplicate, duplicateSource); err == nil ||
		!strings.Contains(err.Error(), `function target "Compute": found 2, want exactly 1`) {
		t.Fatalf("duplicate target error = %v", err)
	}
}

func TestGoFragmentMatcherReportsActionableMismatch(t *testing.T) {
	req := GoMatchRequest{Kind: GoExpression, Expected: "x - y*2", Target: "Result"}
	err := MatchGeneratedGoFragment(req, generatedGoProbe)
	if err == nil {
		t.Fatal("expected mismatch")
	}
	for _, want := range []string{
		`expression target "Result" AST mismatch`,
		"--- expected",
		"+++ generated",
		"-x - y*2",
		"+x + y*2",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("mismatch %q does not contain %q", err, want)
		}
	}
}

func TestGoFragmentMatcherRejectsInvalidExpectedFragment(t *testing.T) {
	req := GoMatchRequest{Kind: GoStatements, Expected: "if {", Target: "Steps"}
	err := MatchGeneratedGoFragment(req, generatedGoProbe)
	if err == nil || !strings.Contains(err.Error(), "parse expected statements") {
		t.Fatalf("invalid fragment error = %v", err)
	}
}
