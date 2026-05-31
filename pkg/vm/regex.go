/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"regexp"
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

func (l *Regex) FindStringSubmatch(s string) []string {
	return l.re.FindStringSubmatch(s)
}

func (l *Regex) FindAllString(s string, n int) []string {
	return l.re.FindAllString(s, n)
}

func (l *Regex) FindAllStringSubmatch(s string, n int) [][]string {
	return l.re.FindAllStringSubmatch(s, n)
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
