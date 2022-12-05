/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
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

func (t *theRegexType) String() string     { return t.Name() }
func (t *theRegexType) Type() ValueType    { return TypeType }
func (t *theRegexType) Unbox() interface{} { return reflect.TypeOf(t) }

func (t *theRegexType) Name() string { return "let-go.lang.Regex" }

func (t *theRegexType) Box(bare interface{}) (Value, error) {
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
func (l *Regex) Unbox() interface{} {
	return l
}

func (l *Regex) String() string {
	return fmt.Sprintf("#%q", l.re)
}

func (l *Regex) ReplaceAll(s string, replacement string) string {
	return l.re.ReplaceAllString(s, replacement)
}

func NewRegex(s string) (Value, error) {
	re, err := regexp.Compile(s)
	if err != nil {
		return NIL, err
	}
	return &Regex{
		re: re,
	}, nil
}
