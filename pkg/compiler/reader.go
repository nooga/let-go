/*
 * Copyright (c) 2021 Marcin Gasperowicz <xnooga@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
 * persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package compiler

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/vm"
)

type LispReader struct {
	inputName string
	pos       int
	line      int
	column    int
	lastRune  rune
	r         *bufio.Reader
}

func NewLispReader(r io.Reader, inputName string) *LispReader {
	return &LispReader{
		inputName: inputName,
		r:         bufio.NewReader(r),
	}
}

func (r *LispReader) next() (rune, error) {
	c, _, err := r.r.ReadRune()
	if err != nil {
		if c == '\n' {
			r.line++
		} else {
			r.column++
		}
		r.pos++
		r.lastRune = c
	}
	return c, err
}

func (r *LispReader) unread() error {
	err := r.r.UnreadRune()
	if err != nil {
		r.pos--
		if r.lastRune == '\n' {
			r.line--
		} else {
			r.column--
		}
	}
	return err
}

func (r *LispReader) peek() (rune, error) {
	for peekBytes := 4; peekBytes > 0; peekBytes-- {
		b, err := r.r.Peek(peekBytes)
		if err == nil {
			ru, _ := utf8.DecodeRune(b)
			if ru == utf8.RuneError {
				return ru, NewReaderError(r, "peek failed - rune error")
			}
			return ru, nil
		}
	}
	return -1, io.EOF
}

func (r *LispReader) eatWhitespace() (rune, error) {
	ch, err := r.next()
	if err != nil {
		return -1, NewReaderError(r, "unexpected error").Wrap(err)
	}
	for isWhitespace(ch) {
		ch, err = r.next()
		if err != nil {
			return -1, NewReaderError(r, "unexpected error").Wrap(err)
		}
	}
	return ch, err
}

func (r *LispReader) Read() (vm.Value, error) {
	for {
		ch, err := r.eatWhitespace()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		if isDigit(ch) {
			return readNumber(r, ch)
		}
		macro, ok := macros[ch]
		if ok {
			return macro(r, ch)
		}
		if ch == '+' || ch == '-' {
			ch2, err := r.next()
			if err != nil {
				return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
			}
			if isDigit(ch2) {
				if err = r.unread(); err != nil {
					return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
				}
				return readNumber(r, ch)
			}
			if err = r.unread(); err != nil {
				return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
			}
		}
		token, err := readToken(r, ch)
		return interpretToken(r, token)
	}
}

func interpretToken(r *LispReader, t vm.Value) (vm.Value, error) {
	s, ok := t.(vm.Symbol)
	if !ok {
		return vm.NIL, NewReaderError(r, fmt.Sprintf("%v is not a symbol", t))
	}
	ss := string(s)
	if ss[0] == ':' {
		return vm.Keyword(ss[1:]), nil
	}
	if ss == "nil" {
		return vm.NIL, nil
	}
	if ss == "true" {
		return vm.TRUE, nil
	}
	if ss == "false" {
		return vm.FALSE, nil
	}
	return t, nil
}

func readToken(r *LispReader, ru rune) (vm.Value, error) {
	s := strings.Builder{}
	s.WriteRune(ru)
	for {
		ch, err := r.next()
		if err != nil {
			if err == io.EOF {
				return vm.Symbol(s.String()), nil
			}
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		if isWhitespace(ch) || isTerminatingMacro(ch) {
			if err = r.unread(); err != nil {
				return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
			}
			return vm.Symbol(s.String()), nil
		}
		s.WriteRune(ch)
	}
}

func readString(r *LispReader, _ rune) (vm.Value, error) {
	s := strings.Builder{}
	for {
		ch, err := r.next()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		if ch == '\\' {
			ch, err := r.next()
			if err != nil {
				return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
			}
			switch ch {
			case 't':
				s.WriteRune('\t')
				continue
			case 'r':
				s.WriteRune('\r')
				continue
			case 'n':
				s.WriteRune('\n')
				continue
			case 'b':
				s.WriteRune('\b')
				continue
			case 'f':
				s.WriteRune('\f')
				continue
			case '\\', '"':
				s.WriteRune(ch)
				continue
			default:
				return vm.NIL, NewReaderError(r, fmt.Sprintf("unknown escape sequence \\%c", ch)).Wrap(err)
			}
		}
		if ch == '"' {
			return vm.String(s.String()), nil
		}
		s.WriteRune(ch)
	}
}

func readNumber(r *LispReader, ru rune) (vm.Value, error) {
	s := strings.Builder{}
	s.WriteRune(ru)
	for {
		ch, err := r.next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return vm.NIL, err
		}
		if isWhitespace(ch) || isTerminatingMacro(ch) {
			if err = r.unread(); err != nil {
				return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
			}
			break
		}
		s.WriteRune(ch)
	}
	sn := s.String()
	i, err := strconv.Atoi(sn)
	if err != nil {
		return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
	}
	return vm.Int(i), nil
}

func readList(r *LispReader, _ rune) (vm.Value, error) {
	var ret []vm.Value
	for {
		ch2, err := r.eatWhitespace()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		if ch2 == ')' {
			break
		}
		if err = r.unread(); err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		form, err := r.Read()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		ret = append(ret, form)
	}
	return vm.ListType.Box(ret)
}

func unmatchedDelimReader(ru rune) readerFunc {
	return func(r *LispReader, _ rune) (vm.Value, error) {
		return nil, NewReaderError(r, fmt.Sprintf("unmatched delimiter %c", ru))
	}
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isTerminatingMacro(r rune) bool {
	return r != '#' && r != '\'' && r != '%' && isMacro(r)
}

func isMacro(r rune) bool {
	_, ok := macros[r]
	return ok
}

type readerFunc func(*LispReader, rune) (vm.Value, error)

var macros map[rune]readerFunc

func init() {
	macros = map[rune]readerFunc{
		'(': readList,
		')': unmatchedDelimReader(')'),
		'"': readString,
	}
}
