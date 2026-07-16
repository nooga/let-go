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
//
// A Regex may be shared freely across goroutines/ExecContexts: regex literals
// compile once at read time into constants visible to all threads, which is
// safe because every method below is a read-only matching API (Go's regexp is
// goroutine-safe except mutating config like Longest). Never add a mutating
// configuration method here without copy-on-write (regexp.Copy).
type Regex struct {
	re             *regexp.Regexp
	pattern        string
	lookaheadGroup int
	startAnchor    anchorKind // leading-anchor class of the lookahead fallback pattern
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

func (l *Regex) ReplaceAll(s string, replacement string) (string, error) {
	if err := l.validateReplacement(replacement); err != nil {
		return "", err
	}
	if l.lookaheadGroup > 0 {
		return l.replaceLookahead(s, replacement, false), nil
	}
	return l.re.ReplaceAllString(s, replacement), nil
}

// ReplaceFirst replaces only the first match, expanding $-group references in
// the replacement (as ReplaceAll does), matching clojure.string/replace-first.
func (l *Regex) ReplaceFirst(s string, replacement string) (string, error) {
	if err := l.validateReplacement(replacement); err != nil {
		return "", err
	}
	if l.lookaheadGroup > 0 {
		return l.replaceLookahead(s, replacement, true), nil
	}
	loc := l.re.FindStringSubmatchIndex(s)
	if loc == nil {
		return s, nil
	}
	out := l.re.ExpandString([]byte(s[:loc[0]]), replacement, s, loc)
	return string(out) + s[loc[1]:], nil
}

// visibleGroupCount is the number of capture groups the USER's pattern has:
// the synthetic group a lookahead fallback appends is an implementation
// detail and not addressable from a replacement template.
func (l *Regex) visibleGroupCount() int {
	n := l.re.NumSubexp()
	if l.lookaheadGroup > 0 {
		n--
	}
	return n
}

func isExpandNameByte(c byte) bool {
	return c == '_' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9'
}

// validateReplacement checks every $-group reference in a replacement
// template against the pattern's capture groups, so a nonexistent reference
// fails loud — Clojure (java.util.regex) throws for a group that doesn't
// exist, while Go's Expand silently substitutes "", which turns a template
// typo into a plausible-looking wrong answer. References are parsed exactly
// as Go's Expand parses them ($$ literal; ${name} or $name with the longest
// run of letters, digits, and underscores; an all-digit name is a group
// index), so anything this rejects is precisely what expansion would have
// silently emptied.
func (l *Regex) validateReplacement(template string) error {
	names := l.re.SubexpNames()
	for i := 0; i < len(template); {
		if template[i] != '$' {
			i++
			continue
		}
		i++
		if i < len(template) && template[i] == '$' { // $$ → literal $
			i++
			continue
		}
		braced := i < len(template) && template[i] == '{'
		if braced {
			i++
		}
		start := i
		for i < len(template) && isExpandNameByte(template[i]) {
			i++
		}
		name := template[start:i]
		if braced {
			if i >= len(template) || template[i] != '}' {
				continue // malformed ${...}: Expand emits it literally
			}
			i++
		}
		if name == "" {
			continue // bare $ before a non-name byte: Expand emits it literally
		}
		if idx, err := strconv.Atoi(name); err == nil {
			if idx > l.visibleGroupCount() { // $0 is the whole match
				return fmt.Errorf("no group %d in regex #%q", idx, l.Pattern())
			}
			continue
		}
		found := false
		for gi, gn := range names {
			// gi != lookaheadGroup also holds trivially when no lookahead
			// group exists (lookaheadGroup 0 = the whole-match slot, whose
			// name is always "" and name here never is).
			if gn == name && gi != l.lookaheadGroup {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no group named %q in regex #%q", name, l.Pattern())
		}
	}
	return nil
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
	return &Regex{
		re:             re,
		pattern:        s,
		lookaheadGroup: re.SubexpIndex(groupName),
		startAnchor:    startAnchorKind(prefix),
	}, nil
}

// startAnchorKind classifies the fallback pattern's leading anchor so
// lookaheadMatches can preserve its semantics across resumed searches.
// Go's regexp has no "match at offset with left context" API — resumed
// searches slice the input, which turns the slice start into a fresh
// text start and lets `^` match at every resume point. Classifying the
// anchor up front lets the iterator constrain resume positions instead.
type anchorKind int

const (
	anchorNone anchorKind = iota // no leading ^/\A: every offset is a valid resume point
	anchorText                   // ^ or \A without (?m): only position 0 can match
	anchorLine                   // (?m)^: only line starts can match
)

func startAnchorKind(prefix string) anchorKind {
	multiline := strings.Contains(prefix, "(?m")
	rest := prefix
	for strings.HasPrefix(rest, "(?") { // skip leading flag groups like (?m) / (?is)
		end := strings.IndexByte(rest, ')')
		if end < 0 || strings.ContainsAny(rest[2:end], ":<P=!") {
			break // a real group, not a bare flag setter
		}
		rest = rest[end+1:]
	}
	switch {
	case strings.HasPrefix(rest, `\A`):
		return anchorText
	case strings.HasPrefix(rest, "^") && multiline:
		return anchorLine
	case strings.HasPrefix(rest, "^"):
		return anchorText
	default:
		return anchorNone
	}
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
		// Resumed searches slice the input, which would let a leading `^`
		// re-anchor at every resume point (Java: ^ matches only at the true
		// string start, or line starts under (?m)). Constrain resume
		// positions per the pattern's classified anchor instead; slicing at
		// a REAL text/line start keeps `^` semantics exact. (`\b` at a
		// resume boundary still sees a fresh text start — a known limitation
		// of Go's offset-less matching API.)
		switch l.startAnchor {
		case anchorText:
			if offset > 0 {
				return result
			}
		case anchorLine:
			for offset > 0 && offset <= len(s) && s[offset-1] != '\n' {
				_, width := utf8.DecodeRuneInString(s[offset:])
				if width == 0 {
					return result
				}
				offset += width
			}
			if offset > len(s) {
				return result
			}
		}
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
		// Expand against the VISIBLE indices (synthetic lookahead group
		// removed): template group numbers then mean the ORIGINAL pattern's
		// groups, and a reference to the synthetic group's slot is out of
		// range — rejected up front by validateReplacement, like any other
		// nonexistent group — instead of leaking the assertion's matched
		// text into the replacement.
		b.Write(l.re.ExpandString(nil, replacement, s, match.visible))
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
