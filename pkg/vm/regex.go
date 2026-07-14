/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package vm

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
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
	re             *regexp.Regexp
	pattern        string
	lookaheadGroup int
}

// Type implements Value
func (l *Regex) Type() ValueType { return RegexType }

// Unbox implements Unbox
func (l *Regex) Unbox() any {
	return l
}

func (l *Regex) String() string {
	return fmt.Sprintf("#%q", l.Pattern())
}

func (l *Regex) ReplaceAll(s string, replacement string) string {
	if l.lookaheadGroup > 0 {
		return l.replaceLookahead(s, replacement, false)
	}
	return l.re.ReplaceAllString(s, replacement)
}

// ReplaceFirst replaces only the first match, expanding $-group references in
// the replacement (as ReplaceAll does), matching clojure.string/replace-first.
func (l *Regex) ReplaceFirst(s string, replacement string) string {
	if l.lookaheadGroup > 0 {
		return l.replaceLookahead(s, replacement, true)
	}
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
	if l.lookaheadGroup > 0 {
		return l.replaceLookaheadFunc(s, f, first)
	}
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
	if l.lookaheadGroup > 0 {
		matches := l.lookaheadMatches(s, 1)
		if len(matches) == 0 {
			return nil
		}
		return stringsFromMatch(s, matches[0].visible)
	}
	return l.re.FindStringSubmatch(s)
}

func (l *Regex) FindStringSubmatchIndex(s string) []int {
	if l.lookaheadGroup > 0 {
		matches := l.lookaheadMatches(s, 1)
		if len(matches) == 0 {
			return nil
		}
		return matches[0].visible
	}
	return l.re.FindStringSubmatchIndex(s)
}

func (l *Regex) FindAllString(s string, n int) []string {
	if l.lookaheadGroup > 0 {
		matches := l.lookaheadMatches(s, n)
		if len(matches) == 0 {
			return nil
		}
		result := make([]string, len(matches))
		for i, match := range matches {
			result[i] = s[match.visible[0]:match.visible[1]]
		}
		return result
	}
	return l.re.FindAllString(s, n)
}

func (l *Regex) FindAllStringSubmatch(s string, n int) [][]string {
	if l.lookaheadGroup > 0 {
		matches := l.lookaheadMatches(s, n)
		if len(matches) == 0 {
			return nil
		}
		result := make([][]string, len(matches))
		for i, match := range matches {
			result[i] = stringsFromMatch(s, match.visible)
		}
		return result
	}
	return l.re.FindAllStringSubmatch(s, n)
}

func (l *Regex) FindAllStringSubmatchIndex(s string, n int) [][]int {
	if l.lookaheadGroup > 0 {
		matches := l.lookaheadMatches(s, n)
		if len(matches) == 0 {
			return nil
		}
		result := make([][]int, len(matches))
		for i, match := range matches {
			result[i] = match.visible
		}
		return result
	}
	return l.re.FindAllStringSubmatchIndex(s, n)
}

func (l *Regex) Split(s string, n int) []string {
	if l.lookaheadGroup > 0 {
		if n == 0 {
			return nil
		}
		matchLimit := -1
		if n > 0 {
			matchLimit = n - 1
		}
		matches := l.lookaheadMatches(s, matchLimit)
		result := make([]string, 0, len(matches)+1)
		last := 0
		for _, match := range matches {
			result = append(result, s[last:match.visible[0]])
			last = match.visible[1]
		}
		return append(result, s[last:])
	}
	return l.re.Split(s, n)
}

// Pattern returns the regex pattern string.
func (l *Regex) Pattern() string {
	if l.pattern != "" {
		return l.pattern
	}
	return l.re.String()
}

func NewRegex(s string) (Value, error) {
	re, err := regexp.Compile(s)
	if err == nil {
		return &Regex{re: re, pattern: s}, nil
	}

	prefix, assertion, ok := terminalPositiveLookahead(s)
	if !ok {
		return NIL, err
	}
	assertionRE, assertionErr := regexp.Compile(assertion)
	if assertionErr != nil || assertionRE.NumSubexp() != 0 {
		return NIL, err
	}
	const groupName = "__let_go_terminal_lookahead"
	re, fallbackErr := regexp.Compile(prefix + "(?P<" + groupName + ">" + assertion + ")")
	if fallbackErr != nil {
		return NIL, err
	}
	return &Regex{re: re, pattern: s, lookaheadGroup: re.SubexpIndex(groupName)}, nil
}

type lookaheadMatch struct {
	raw     []int
	visible []int
}

func terminalPositiveLookahead(pattern string) (string, string, bool) {
	if !strings.HasSuffix(pattern, ")") {
		return "", "", false
	}
	start := strings.LastIndex(pattern, "(?=")
	if start < 0 || start+3 == len(pattern)-1 {
		return "", "", false
	}
	return pattern[:start], pattern[start+3 : len(pattern)-1], true
}

func (l *Regex) lookaheadMatches(s string, n int) []lookaheadMatch {
	if n == 0 {
		return nil
	}
	result := make([]lookaheadMatch, 0)
	for offset := 0; offset <= len(s) && (n < 0 || len(result) < n); {
		raw := l.re.FindStringSubmatchIndex(s[offset:])
		if raw == nil {
			break
		}
		for i, index := range raw {
			if index >= 0 {
				raw[i] = index + offset
			}
		}
		assertionStart := raw[2*l.lookaheadGroup]
		visible := make([]int, 0, len(raw)-2)
		visible = append(visible, raw[:2*l.lookaheadGroup]...)
		visible = append(visible, raw[2*l.lookaheadGroup+2:]...)
		visible[1] = assertionStart
		result = append(result, lookaheadMatch{raw: raw, visible: visible})

		next := assertionStart
		if next <= offset {
			if offset == len(s) {
				next = offset + 1
			} else {
				_, width := utf8.DecodeRuneInString(s[offset:])
				next = offset + width
			}
		}
		offset = next
	}
	return result
}

func stringsFromMatch(s string, indices []int) []string {
	groups := make([]string, len(indices)/2)
	for i := range groups {
		if indices[2*i] >= 0 {
			groups[i] = s[indices[2*i]:indices[2*i+1]]
		}
	}
	return groups
}

func (l *Regex) replaceLookahead(s, replacement string, first bool) string {
	limit := -1
	if first {
		limit = 1
	}
	matches := l.lookaheadMatches(s, limit)
	if len(matches) == 0 {
		return s
	}
	var b strings.Builder
	last := 0
	for _, match := range matches {
		b.WriteString(s[last:match.visible[0]])
		expandIndices := append([]int(nil), match.raw...)
		expandIndices[1] = match.visible[1]
		b.Write(l.re.ExpandString(nil, replacement, s, expandIndices))
		last = match.visible[1]
	}
	b.WriteString(s[last:])
	return b.String()
}

func (l *Regex) replaceLookaheadFunc(s string, f func(groups []string, present []bool) (string, error), first bool) (string, error) {
	limit := -1
	if first {
		limit = 1
	}
	matches := l.lookaheadMatches(s, limit)
	if len(matches) == 0 {
		return s, nil
	}
	var b strings.Builder
	last := 0
	for _, match := range matches {
		b.WriteString(s[last:match.visible[0]])
		groups := stringsFromMatch(s, match.visible)
		present := make([]bool, len(groups))
		for i := range present {
			present[i] = match.visible[2*i] >= 0
		}
		replacement, err := f(groups, present)
		if err != nil {
			return "", err
		}
		b.WriteString(replacement)
		last = match.visible[1]
	}
	b.WriteString(s[last:])
	return b.String(), nil
}
