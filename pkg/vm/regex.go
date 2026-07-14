/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type theRegexType struct {
}

func (t *theRegexType) String() string  { return t.Name() }
func (t *theRegexType) Type() ValueType { return TypeType }
func (t *theRegexType) Unbox() any      { return reflect.TypeFor[*theRegexType]() }

func (t *theRegexType) Name() string { return "let-go.lang.Regex" }

func (t *theRegexType) Box(bare any) (Value, error) {
	raw, ok := bare.(*regexp.Regexp)
	if !ok {
		return NIL, NewTypeError(bare, "can't be boxed as", t)
	}
	return &Regex{re: raw}, nil
}

// RegexType is the type of RegexValues
var RegexType *theRegexType = &theRegexType{}

// Regex is boxed int
//
// A Regex may be shared freely across goroutines/ExecContexts: regex literals
// compile once at read time into constants visible to all threads, which is
// safe because every method below is a read-only matching API (Go's regexp is
// goroutine-safe except mutating config like Longest). Never add a mutating
// configuration method here without copy-on-write (regexp.Copy).
type Regex struct {
	re *regexp.Regexp
}

// Type implements Value
func (l *Regex) Type() ValueType { return RegexType }

// Unbox implements Unbox
func (l *Regex) Unbox() any {
	return l
}

func (l *Regex) String() string {
	return fmt.Sprintf("#%q", l.re)
}

func (l *Regex) ReplaceAll(s string, replacement string) string {
	return l.re.ReplaceAllString(s, replacement)
}

// ReplaceFirst replaces only the first match, expanding $-group references in
// the replacement (as ReplaceAll does), matching clojure.string/replace-first.
func (l *Regex) ReplaceFirst(s string, replacement string) string {
	loc := l.re.FindStringSubmatchIndex(s)
	if loc == nil {
		return s
	}
	out := l.re.ExpandString([]byte(s[:loc[0]]), replacement, s, loc)
	return string(out) + s[loc[1]:]
}

// ReplaceAllFunc replaces every non-overlapping match using f, matching
// clojure.string/replace with a function replacement. For each match, f
// receives the submatch groups: groups[0] is the whole match, groups[i] the
// i-th capture group; present[i] is false for a group that did not
// participate (so the caller can pass nil, as Clojure's re-groups does).
func (l *Regex) ReplaceAllFunc(s string, f func(groups []string, present []bool) (string, error)) (string, error) {
	return l.replaceFunc(s, f, false)
}

// ReplaceFirstFunc is ReplaceAllFunc limited to the first match, matching
// clojure.string/replace-first with a function replacement.
func (l *Regex) ReplaceFirstFunc(s string, f func(groups []string, present []bool) (string, error)) (string, error) {
	return l.replaceFunc(s, f, true)
}

func (l *Regex) replaceFunc(s string, f func(groups []string, present []bool) (string, error), first bool) (string, error) {
	n := -1
	if first {
		n = 1
	}
	matches := l.re.FindAllStringSubmatchIndex(s, n)
	if matches == nil {
		return s, nil
	}
	var b strings.Builder
	last := 0
	for _, m := range matches {
		b.WriteString(s[last:m[0]])
		ng := len(m) / 2
		groups := make([]string, ng)
		present := make([]bool, ng)
		for i := 0; i < ng; i++ {
			if m[2*i] >= 0 {
				groups[i] = s[m[2*i]:m[2*i+1]]
				present[i] = true
			}
		}
		rep, err := f(groups, present)
		if err != nil {
			return "", err
		}
		b.WriteString(rep)
		last = m[1]
	}
	b.WriteString(s[last:])
	return b.String(), nil
}

func (l *Regex) FindStringSubmatch(s string) []string {
	return l.re.FindStringSubmatch(s)
}

func (l *Regex) FindStringSubmatchIndex(s string) []int {
	return l.re.FindStringSubmatchIndex(s)
}

func (l *Regex) FindAllString(s string, n int) []string {
	return l.re.FindAllString(s, n)
}

func (l *Regex) FindAllStringSubmatch(s string, n int) [][]string {
	return l.re.FindAllStringSubmatch(s, n)
}

func (l *Regex) FindAllStringSubmatchIndex(s string, n int) [][]int {
	return l.re.FindAllStringSubmatchIndex(s, n)
}

func (l *Regex) Split(s string, n int) []string {
	return l.re.Split(s, n)
}

// Pattern returns the regex pattern string.
func (l *Regex) Pattern() string { return l.re.String() }

func NewRegex(s string) (Value, error) {
	re, err := regexp.Compile(s)
	if err != nil {
		return NIL, err
	}
	return &Regex{
		re: re,
	}, nil
}

// MustRegexFromReadable reconstructs a Regex from its readable #"..." form
// (String()'s output: '#' followed by a Go-quoted pattern). Generated code
// uses it to rebuild read-time-compiled regex literal constants; the pattern
// was validated at read time, so failure here indicates a codegen bug.
func MustRegexFromReadable(s string) *Regex {
	pat, err := strconv.Unquote(s[1:])
	if err != nil {
		panic(fmt.Sprintf("MustRegexFromReadable: bad readable form %s: %v", s, err))
	}
	v, err := NewRegex(pat)
	if err != nil {
		panic(fmt.Sprintf("MustRegexFromReadable: %s: %v", s, err))
	}
	return v.(*Regex)
}
