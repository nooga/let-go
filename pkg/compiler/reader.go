/*
 * Copyright (c) 2021-2026 Marcin Gasperowicz <xnooga@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/nooga/let-go/pkg/errors"
	"github.com/nooga/let-go/pkg/rt"

	"github.com/nooga/let-go/pkg/vm"
)

type TokenKind int

const (
	TokenString TokenKind = iota
	TokenNumber
	TokenKeyword
	TokenSymbol
	TokenChar
	TokenSpecial
	TokenComment
	TokenPunctuation
)

type Token struct {
	Start int
	End   int
	Kind  TokenKind
}

type LispReader struct {
	inputName  string
	pos        int
	line       int
	column     int
	lastCol    int
	lastRune   rune
	maxPercent int
	r          *bufio.Reader

	Tokens     []Token
	tokenizing bool
	splicing   bool
}

func NewLispReader(r io.Reader, inputName string) *LispReader {
	return &LispReader{
		inputName: inputName,
		r:         bufio.NewReader(r),
	}
}

func NewLispReaderTokenizing(r io.Reader, inputName string) *LispReader {
	return &LispReader{
		inputName:  inputName,
		r:          bufio.NewReader(r),
		Tokens:     []Token{},
		tokenizing: true,
	}
}

func (r *LispReader) openToken() {
	if !r.tokenizing {
		return
	}
	if len(r.Tokens) > 0 && r.Tokens[len(r.Tokens)-1].End == -1 {
		r.Tokens[len(r.Tokens)-1].Start = r.pos - 1
		return
	}
	r.Tokens = append(r.Tokens, Token{Start: r.pos - 1, End: -1})
}

func (r *LispReader) discardToken() {
	if !r.tokenizing {
		return
	}
	r.Tokens = r.Tokens[:len(r.Tokens)-1]
}

func (r *LispReader) closeToken(kind TokenKind) {
	if !r.tokenizing {
		return
	}
	if r.Tokens[len(r.Tokens)-1].End != -1 {
		return
	}
	r.Tokens[len(r.Tokens)-1].End = r.pos
	r.Tokens[len(r.Tokens)-1].Kind = kind
}

func (r *LispReader) addToken(kind TokenKind) {
	if !r.tokenizing {
		return
	}
	r.openToken()
	r.closeToken(kind)
}

func (r *LispReader) next() (rune, error) {
	c, _, err := r.r.ReadRune()
	if err == nil {
		if c == '\n' {
			r.line++
			r.lastCol = r.column
			r.column = 0
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
	if err == nil {
		r.pos--
		if r.lastRune == '\n' {
			r.line--
			r.column = r.lastCol
		} else {
			r.column--
		}
	}
	return err
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

func appendNonVoid(r *LispReader, vs []vm.Value, v vm.Value) []vm.Value {
	if v.Type() == vm.VoidType {
		r.splicing = false
		return vs
	}
	if r.splicing {
		r.splicing = false
		seq, ok := v.(vm.Seq)
		if !ok {
			return vs
		}
		for seq != nil && seq != vm.EmptyList {
			vs = append(vs, seq.First())
			seq = seq.Next()
		}
		return vs
	}
	return append(vs, v)
}

func (r *LispReader) Read() (vm.Value, error) {
	ch, err := r.eatWhitespace()
	if err != nil {
		return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
	}
	r.openToken()
	if isDigit(ch) {
		return readNumber(r, ch)
	}
	macro, ok := macros[ch]
	if ok {
		r.closeToken(TokenPunctuation)
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
	if err != nil {
		return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
	}
	return interpretToken(r, token)
}

func interpretToken(r *LispReader, t vm.Value) (vm.Value, error) {
	s, ok := t.(vm.Symbol)
	if !ok {
		return vm.NIL, NewReaderError(r, fmt.Sprintf("%v is not a symbol", t))
	}
	ss := string(s)
	if ss[0] == ':' {
		nom := ss[1:]
		if nom == "" {
			return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid token: %s", ss))
		}
		if nom[0] == ':' {
			// Auto-resolved keyword (::foo or ::alias/foo)
			onom := nom[1:]
			if strings.Contains(onom, "/") {
				// ::alias/name — resolve alias to full namespace name
				parts := strings.SplitN(onom, "/", 2)
				alias := parts[0]
				name := parts[1]
				if alias == "" || name == "" || strings.ContainsAny(name, ":/") {
					return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid token: %s", ss))
				}
				// Look up alias in current namespace
				cns := rt.CurrentNS.Deref().(*vm.Namespace)
				resolved := cns.ResolveAlias(vm.Symbol(alias))
				if resolved == nil {
					return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid token: %s, namespace alias %s not found", ss, alias))
				}
				nom = resolved.Name() + "/" + name
			} else {
				// ::foo — resolve to current namespace
				if strings.ContainsAny(onom, ":") {
					return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid token: %s", ss))
				}
				nom = rt.CurrentNS.Deref().(*vm.Namespace).Name() + "/" + onom
			}
		}
		if strings.ContainsAny(nom, ":") {
			return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid token: %s", ss))
		}
		r.closeToken(TokenKeyword)
		return vm.Keyword(nom), nil
	}
	if ss == "nil" {
		r.closeToken(TokenSpecial)
		return vm.NIL, nil
	}
	if ss == "true" {
		r.closeToken(TokenSpecial)
		return vm.TRUE, nil
	}
	if ss == "false" {
		r.closeToken(TokenSpecial)
		return vm.FALSE, nil
	}
	if ss[0] == '%' {
		var n int
		var err error
		if ss == "%" {
			n = 1
		} else {
			n, err = strconv.Atoi(ss[1:])
		}
		if err == nil && n >= 0 && n > r.maxPercent {
			r.maxPercent = n
		}
		if ss == "%" {
			return vm.Symbol("%1"), nil
		}
	}
	if _, ok := specialForms[t.(vm.Symbol)]; ok {
		r.closeToken(TokenSpecial)
	} else {
		r.closeToken(TokenSymbol)
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
	r.discardToken()
	r.openToken()
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
			case 'u':
				hex := ""
				for i := 0; i < 4; i++ {
					ch, err := r.next()
					if err != nil || !isHexDigit(ch) {
						break
					}
					hex += string(ch)
				}
				if len(hex) < 4 {
					return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid escape sequence \\u%s", hex)).Wrap(err)
				}
				var hexi int
				n, err := fmt.Sscanf(hex, "%x", &hexi)
				if n != 1 || (hexi >= 0xD800 && hexi <= 0xDFFF) {
					return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid escape sequence \\u%s", hex)).Wrap(err)
				}
				if err != nil {
					return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid escape sequence \\u%s", hex)).Wrap(err)
				}
				s.WriteRune(rune(hexi))
				continue
			default:
				return vm.NIL, NewReaderError(r, fmt.Sprintf("unknown escape sequence \\%c", ch)).Wrap(err)
			}
		}
		if ch == '"' {
			r.closeToken(TokenString)
			return vm.String(s.String()), nil
		}
		s.WriteRune(ch)
	}
}

func readRegex(r *LispReader, _ rune) (vm.Value, error) {
	// Read regex pattern literally — backslashes are regex escapes, not string escapes
	r.discardToken()
	r.openToken()
	s := strings.Builder{}
	for {
		ch, err := r.next()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unterminated regex").Wrap(err)
		}
		if ch == '\\' {
			// Pass through backslash + next char literally
			ch2, err := r.next()
			if err != nil {
				return vm.NIL, NewReaderError(r, "unterminated regex escape").Wrap(err)
			}
			s.WriteRune('\\')
			s.WriteRune(ch2)
			continue
		}
		if ch == '"' {
			r.closeToken(TokenString)
			return vm.ListType.Box([]vm.Value{vm.Symbol("re-pattern"), vm.String(s.String())})
		}
		s.WriteRune(ch)
	}
}

func isHexDigit(ch rune) bool {
	if unicode.IsDigit(ch) {
		return true
	}
	if ch >= 'a' && ch <= 'f' {
		return true
	}
	if ch >= 'A' && ch <= 'F' {
		return true
	}
	return false
}

func readChar(r *LispReader, _ rune) (vm.Value, error) {
	ch, err := r.next()
	if err != nil {
		return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
	}
	tok, err := readToken(r, ch)
	if err != nil {
		return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
	}
	toks := tok.Unbox().(string)
	ru, s := utf8.DecodeRuneInString(toks)
	switch {
	case len(toks) == s:
		return vm.Char(ru), nil
	case toks == "space":
		return vm.Char(' '), nil
	case toks == "tab":
		return vm.Char('\t'), nil
	case toks == "backspace":
		return vm.Char('\b'), nil
	case toks == "newline":
		return vm.Char('\n'), nil
	case toks == "formfeed":
		return vm.Char('\f'), nil
	case toks == "return":
		return vm.Char('\r'), nil
	case toks[0] == 'u':
		hex := toks[1:]
		if len(hex) < 4 {
			goto fail // LOL I'm using goto in 2021 because in Go it actually makes sense
		}
		var hexi int
		n, err := fmt.Sscanf(hex, "%x", &hexi)
		if n != 1 || (hexi >= 0xD800 && hexi <= 0xDFFF) {
			goto fail
		}
		if err != nil {
			return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid char constant \\%s", toks)).Wrap(err)
		}
		return vm.Char(rune(hexi)), nil
	case toks[0] == 'o':
		hex := toks[1:]
		if len(hex) > 3 {
			goto fail
		}
		var hexi int
		n, err := fmt.Sscanf(hex, "%o", &hexi)
		if n != 1 || hexi > 0377 {
			goto fail
		}
		if err != nil {
			return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid char constant \\%s", toks)).Wrap(err)
		}
		return vm.Char(rune(hexi)), nil
	}
fail:
	return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid char constant \\%s", toks))
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
	// Check for explicit BigInt suffix: 123N
	if len(sn) > 1 && (sn[len(sn)-1] == 'N' || sn[len(sn)-1] == 'n') {
		numStr := sn[:len(sn)-1]
		if bi, ok := vm.NewBigIntFromString(numStr); ok {
			r.closeToken(TokenNumber)
			return bi, nil
		}
	}
	// Check for BigDecimal suffix: 123.456M → Float approximation
	if len(sn) > 1 && (sn[len(sn)-1] == 'M' || sn[len(sn)-1] == 'm') {
		numStr := sn[:len(sn)-1]
		if f, err := strconv.ParseFloat(numStr, 64); err == nil {
			r.closeToken(TokenNumber)
			return vm.Float(f), nil
		}
		if i, err := strconv.Atoi(numStr); err == nil {
			r.closeToken(TokenNumber)
			return vm.Float(float64(i)), nil
		}
	}
	// Handle optional negative sign for special number formats
	numStr := sn
	neg := false
	if len(numStr) > 0 && numStr[0] == '-' {
		neg = true
		numStr = numStr[1:]
	}
	// Hex literal: 0xff or -0xff
	if len(numStr) > 2 && numStr[0] == '0' && (numStr[1] == 'x' || numStr[1] == 'X') {
		if i, err := strconv.ParseUint(numStr[2:], 16, 64); err == nil {
			n := int64(i)
			if neg {
				n = -n
			}
			r.closeToken(TokenNumber)
			return vm.MakeInt(int(n)), nil
		}
	}
	// Octal literal: 0377
	if len(numStr) > 1 && numStr[0] == '0' && numStr[1] >= '0' && numStr[1] <= '7' {
		if i, err := strconv.ParseInt(numStr[1:], 8, 64); err == nil {
			if neg {
				i = -i
			}
			r.closeToken(TokenNumber)
			return vm.MakeInt(int(i)), nil
		}
	}
	// Radix literal: 2r1010, 16rFF
	if idx := strings.IndexByte(numStr, 'r'); idx > 0 {
		if radix, err := strconv.Atoi(numStr[:idx]); err == nil && radix >= 2 && radix <= 36 {
			if i, err := strconv.ParseInt(numStr[idx+1:], radix, 64); err == nil {
				if neg {
					i = -i
				}
				r.closeToken(TokenNumber)
				return vm.MakeInt(int(i)), nil
			}
		}
	}
	// Ratio literal: 1/2, 22/7
	if idx := strings.IndexByte(sn, '/'); idx > 0 {
		num, nerr := strconv.ParseFloat(sn[:idx], 64)
		den, derr := strconv.ParseFloat(sn[idx+1:], 64)
		if nerr == nil && derr == nil && den != 0 {
			r.closeToken(TokenNumber)
			return vm.Float(num / den), nil
		}
	}
	// Try int first
	i, err := strconv.Atoi(sn)
	if err == nil {
		r.closeToken(TokenNumber)
		return vm.MakeInt(i), nil
	}
	// If Atoi failed due to range, try BigInt
	if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
		if bi, ok := vm.NewBigIntFromString(sn); ok {
			r.closeToken(TokenNumber)
			return bi, nil
		}
	}
	// Try float
	f, ferr := strconv.ParseFloat(sn, 64)
	if ferr == nil {
		r.closeToken(TokenNumber)
		return vm.Float(f), nil
	}
	return vm.NIL, NewReaderError(r, fmt.Sprintf("invalid number: %s", sn))
}

func readList(r *LispReader, _ rune) (vm.Value, error) {
	startLine := r.line
	startCol := r.column - 1 // -1 because '(' was already consumed
	if startCol < 0 {
		startCol = 0
	}
	var ret []vm.Value
	for {
		ch2, err := r.eatWhitespace()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		if ch2 == ')' {
			r.addToken(TokenPunctuation)
			break
		}
		if err = r.unread(); err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		form, err := r.Read()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		ret = appendNonVoid(r, ret, form)
	}
	result, err := vm.ListType.Box(ret)
	if err != nil {
		return vm.NIL, err
	}
	vm.FormSource.Set(result.(vm.Value), vm.SourceInfo{
		File: r.inputName, Line: startLine, Column: startCol,
	})
	return result, nil
}

func readVector(r *LispReader, _ rune) (vm.Value, error) {
	startLine := r.line
	startCol := r.column - 1
	if startCol < 0 {
		startCol = 0
	}
	ret := make([]vm.Value, 0)
	for {
		ch2, err := r.eatWhitespace()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		if ch2 == ']' {
			r.addToken(TokenPunctuation)
			break
		}
		if err = r.unread(); err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		form, err := r.Read()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		ret = appendNonVoid(r, ret, form)
	}
	result := vm.ArrayVector(ret)
	vm.FormSource.Set(result, vm.SourceInfo{
		File: r.inputName, Line: startLine, Column: startCol,
	})
	return result, nil
}

func readMap(r *LispReader, _ rune) (vm.Value, error) {
	startLine := r.line
	startCol := r.column - 1
	if startCol < 0 {
		startCol = 0
	}
	ret := make([]vm.Value, 0)
	for {
		ch2, err := r.eatWhitespace()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		if ch2 == '}' {
			r.addToken(TokenPunctuation)
			break
		}
		if err = r.unread(); err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		form, err := r.Read()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		ret = appendNonVoid(r, ret, form)
	}
	if len(ret)%2 != 0 {
		return vm.NIL, NewReaderError(r, "map literal must contain even number of forms")
	}
	result := vm.NewMap(ret)
	vm.FormSource.Set(result, vm.SourceInfo{
		File: r.inputName, Line: startLine, Column: startCol,
	})
	return result, nil
}

func readSet(r *LispReader, _ rune) (vm.Value, error) {
	startLine := r.line
	startCol := r.column - 2 // -2 because '#' and '{' were consumed
	if startCol < 0 {
		startCol = 0
	}
	ret := vm.EmptyList
	for {
		ch2, err := r.eatWhitespace()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		if ch2 == '}' {
			r.addToken(TokenPunctuation)
			break
		}
		if err = r.unread(); err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		form, err := r.Read()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
		}
		if form.Type() != vm.VoidType {
			ret = ret.Conj(form).(*vm.List)
		}
	}
	result := ret.Cons(vm.Symbol("hash-set"))
	vm.FormSource.Set(result, vm.SourceInfo{
		File: r.inputName, Line: startLine, Column: startCol,
	})
	return result, nil
}

func readQuote(r *LispReader, _ rune) (vm.Value, error) {
	form, err := r.Read()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading quoted form").Wrap(err)
	}
	quote := vm.Symbol("quote")
	ret, err := vm.ListType.Box([]vm.Value{quote, form})
	if err != nil {
		return vm.NIL, NewReaderError(r, "boxing quoted form").Wrap(err)
	}
	return ret, nil
}

func readDeref(r *LispReader, _ rune) (vm.Value, error) {
	form, err := r.Read()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading deref form").Wrap(err)
	}
	deref := vm.Symbol("deref")
	ret, err := vm.ListType.Box([]vm.Value{deref, form})
	if err != nil {
		return vm.NIL, NewReaderError(r, "boxing deref form").Wrap(err)
	}
	return ret, nil
}

type gensymEnv struct {
	syms map[vm.Symbol]vm.Symbol
	cnt  int
}

func (g *gensymEnv) Get(s vm.Symbol) vm.Value {
	y, ok := g.syms[s]
	if !ok {
		return vm.NIL
	}
	return y
}

func (g *gensymEnv) Set(s vm.Symbol) vm.Symbol {
	y := vm.Symbol(fmt.Sprintf("%s__G__%d", s[0:len(s)-2], g.cnt))
	g.syms[s] = y
	g.cnt += 1
	return y
}

func readSyntaxQuote(r *LispReader, _ rune) (vm.Value, error) {
	form, err := r.Read()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading quoted form").Wrap(err)
	}
	return syntaxQuote(r, form, &gensymEnv{
		syms: map[vm.Symbol]vm.Symbol{},
		cnt:  0,
	})
}

func syntaxQuote(r *LispReader, form vm.Value, env *gensymEnv) (vm.Value, error) {
	switch {
	case form.Type() == vm.SymbolType:
		sform := form.(vm.Symbol)
		if _, ok := specialForms[sform]; ok {
			ret, err := vm.ListType.Box([]vm.Value{vm.Symbol("quote"), form})
			if err != nil {
				return vm.NIL, NewReaderError(r, "boxing syntax-quoted special form").Wrap(err)
			}
			return ret, nil
		}
		ns, _ := sform.Namespaced()
		if ns == vm.NIL && sform[len(sform)-1] == '#' {
			gsym := env.Get(sform)
			if gsym == vm.NIL {
				gsym = env.Set(sform)
			}
			ret, err := vm.ListType.Box([]vm.Value{vm.Symbol("quote"), gsym})
			if err != nil {
				return vm.NIL, NewReaderError(r, "boxing syntax-quoted gensym form").Wrap(err)
			}
			return ret, nil
		}
		ret, err := vm.ListType.Box([]vm.Value{vm.Symbol("quote"), form})
		if err != nil {
			return vm.NIL, NewReaderError(r, "boxing syntax-quoted special form").Wrap(err)
		}
		return ret, nil
	case isUnquote(form):
		vl := form.(*vm.List)
		return vl.Next().First(), nil
	case isUnquoteSplicing(form):
		return vm.NIL, NewReaderError(r, "unquote-splicing used outside of a list")
	case form.Type() == vm.ArrayVectorType:
		uq, err := expandUnquotes(r, form, env)
		if err != nil {
			return vm.NIL, NewReaderError(r, "expanding unquotes for vector")
		}
		vv, err := vm.ListType.Box([]vm.Value{vm.Symbol("apply"), vm.Symbol("concat*"), uq})
		if err != nil {
			return vm.NIL, NewReaderError(r, "boxing unquoted vector form")
		}
		return vm.ListType.Box([]vm.Value{vm.Symbol("apply"), vm.Symbol("vector"), vv})
	case form.Type() == vm.MapType:
		lform := flattenMap(form)
		uq, err := expandUnquotes(r, lform, env)
		if err != nil {
			return vm.NIL, NewReaderError(r, "expanding unquotes for vector")
		}
		vv, err := vm.ListType.Box([]vm.Value{vm.Symbol("apply"), vm.Symbol("concat*"), uq})
		if err != nil {
			return vm.NIL, NewReaderError(r, "boxing unquoted vector form")
		}
		return vm.ListType.Box([]vm.Value{vm.Symbol("apply"), vm.Symbol("hash-map"), vv})
	case form.Type() == vm.ListType:
		uq, err := expandUnquotes(r, form, env)
		if err != nil {
			return vm.NIL, NewReaderError(r, "expanding unquotes for list")
		}
		vv, err := vm.ListType.Box([]vm.Value{vm.Symbol("apply"), vm.Symbol("concat*"), uq})
		if err != nil {
			return vm.NIL, NewReaderError(r, "boxing unquoted list form")
		}
		return vv, nil
	default:
		ret, err := vm.ListType.Box([]vm.Value{vm.Symbol("quote"), form})
		if err != nil {
			return vm.NIL, NewReaderError(r, "boxing syntax-quoted form").Wrap(err)
		}
		return ret, nil
	}
}

func flattenMap(form vm.Value) vm.Value {
	sq, ok := form.(vm.Sequable)
	if !ok {
		return vm.ArrayVector{}
	}
	ret := vm.ArrayVector{}
	for s := sq.Seq(); s != nil && s != vm.EmptyList; s = s.Next() {
		entry := s.First().(vm.ArrayVector)
		ret = append(ret, entry[0], entry[1])
	}
	return ret
}

func expandUnquotes(r *LispReader, form vm.Value, env *gensymEnv) (vm.Value, error) {
	seqable, ok := form.(vm.Sequable)
	if !ok {
		return vm.NIL, NewReaderError(r, "expandUnquotes expected Sequable")
	}
	fseq := seqable.Seq()
	if fseq == nil || fseq == vm.EmptyList {
		return vm.ArrayVector{}, nil
	}
	ret := vm.ArrayVector{}
	for fseq != nil && fseq != vm.EmptyList {
		v := fseq.First()
		switch {
		case isUnquote(v):
			ret = append(ret, vm.ArrayVector{v.(vm.Sequable).Seq().Next().First()})
		case isUnquoteSplicing(v):
			ret = append(ret, v.(vm.Sequable).Seq().Next().First())
		default:
			vq, err := syntaxQuote(r, v, env)
			if err != nil {
				return vm.NIL, err
			}
			ret = append(ret, vm.ArrayVector{vq})
		}
		fseq = fseq.Next()
	}
	return ret, nil
}

func isUnquote(v vm.Value) bool {
	if v.Type() != vm.ListType {
		return false
	}
	vl := v.(*vm.List)
	if vl.RawCount() != 2 {
		return false
	}
	return vl.First() == vm.Symbol("unquote")
}

func isUnquoteSplicing(v vm.Value) bool {
	if v.Type() != vm.ListType {
		return false
	}
	vl := v.(*vm.List)
	if vl.RawCount() != 2 {
		return false
	}
	return vl.First() == vm.Symbol("unquote-splicing")
}

func readUnquote(r *LispReader, _ rune) (vm.Value, error) {
	ch, err := r.next()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading unquote prefix").Wrap(err)
	}
	sym := vm.Symbol("unquote")
	if ch == '@' {
		sym = "unquote-splicing"
	} else {
		err = r.unread()
		if err != nil {
			return vm.NIL, NewReaderError(r, "unreading unquoted form").Wrap(err)
		}
	}
	form, err := r.Read()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading unquoted form").Wrap(err)
	}
	ret, err := vm.ListType.Box([]vm.Value{sym, form})
	if err != nil {
		return vm.NIL, NewReaderError(r, "boxing unquoted form").Wrap(err)
	}
	return ret, nil
}

func readVarQuote(r *LispReader, _ rune) (vm.Value, error) {
	form, err := r.Read()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading quoted var").Wrap(err)
	}
	if form.Type() != vm.SymbolType {
		return vm.NIL, NewReaderError(r, "invalid var quote")
	}
	quote := vm.Symbol("var")
	ret, err := vm.ListType.Box([]vm.Value{quote, form})
	if err != nil {
		return vm.NIL, NewReaderError(r, "boxing quoted var").Wrap(err)
	}
	return ret, nil
}

func readShortFn(r *LispReader, _ rune) (vm.Value, error) {
	startLine := r.line
	startCol := r.column - 2 // -2 because '#' and '(' were consumed
	if startCol < 0 {
		startCol = 0
	}
	var ret []vm.Value
	r.maxPercent = 0
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
		ret = appendNonVoid(r, ret, form)
	}
	var percents []vm.Value
	for i := 1; i <= r.maxPercent; i++ {
		percents = append(percents, vm.Symbol(fmt.Sprintf("%%%d", i)))
	}
	r.maxPercent = 0
	body, err := vm.ListType.Box(ret)
	if err != nil {
		return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
	}
	fn, err := vm.ListType.Box([]vm.Value{
		vm.Symbol("fn*"),
		vm.NewArrayVector(percents),
		body,
	})
	if err != nil {
		return vm.NIL, NewReaderError(r, "unexpected error").Wrap(err)
	}
	vm.FormSource.Set(fn.(vm.Value), vm.SourceInfo{
		File: r.inputName, Line: startLine, Column: startCol,
	})
	return fn, nil
}

const lgConditionalTag = vm.Keyword("lg")
const defaultConditionalTag = vm.Keyword("default")

func readConditional(r *LispReader, s rune) (vm.Value, error) {
	n, err := r.next()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading reader conditional")
	}
	splicing := false
	if n == '@' {
		splicing = true
		n, err = r.next()
		if err != nil {
			return vm.NIL, NewReaderError(r, "reading splicing reader conditional")
		}
	}
	if n != '(' {
		return vm.NIL, NewReaderError(r, "reading reader conditional, list expected")
	}

	// Parse key-value pairs manually so we can skip unreadable branches.
	var form vm.Value = vm.VOID
	found := false
	for {
		ch, err := r.eatWhitespace()
		if err != nil {
			return vm.NIL, NewReaderError(r, "reading reader conditional")
		}
		if ch == ')' {
			break
		}
		// Read the platform key (must be a keyword)
		if err = r.unread(); err != nil {
			return vm.NIL, NewReaderError(r, "reading reader conditional")
		}
		key, err := r.Read()
		if err != nil {
			return vm.NIL, NewReaderError(r, "reading reader conditional key")
		}
		isMatch := !found && (key == lgConditionalTag || key == defaultConditionalTag)
		if isMatch {
			// Read the value form normally
			val, err := r.Read()
			if err != nil {
				return vm.NIL, NewReaderError(r, "reading reader conditional value")
			}
			form = val
			found = true
		} else {
			// Skip the value form — it may contain dialect-specific syntax
			// we can't parse. Count balanced parens/brackets/braces.
			skipReaderForm(r)
		}
	}
	r.splicing = splicing
	return form, nil
}

// skipReaderForm consumes a single form from the reader, handling balanced
// delimiters. Used to skip unmatched reader conditional branches that may
// contain syntax our reader doesn't support.
func skipReaderForm(r *LispReader) {
	ch, err := r.eatWhitespace()
	if err != nil {
		return
	}
	switch ch {
	case '(', '[', '{':
		close := map[rune]rune{'(': ')', '[': ']', '{': '}'}[ch]
		depth := 1
		inString := false
		for depth > 0 {
			c, err := r.next()
			if err != nil {
				return
			}
			if inString {
				if c == '\\' {
					r.next() // skip escaped char
				} else if c == '"' {
					inString = false
				}
				continue
			}
			switch c {
			case '"':
				inString = true
			case '(', '[', '{':
				depth++
			case close:
				depth--
			case ')', ']', '}':
				// closing a different delimiter type — still decrement if matching any open
				depth--
			}
		}
	case '"':
		// Skip string
		for {
			c, err := r.next()
			if err != nil || c == '"' {
				return
			}
			if c == '\\' {
				r.next()
			}
		}
	case '#':
		// Hash dispatch — skip the next form too
		c, err := r.next()
		if err != nil {
			return
		}
		if c == '(' || c == '{' {
			r.unread()
			skipReaderForm(r)
		} else {
			// #something — skip to whitespace/delimiter
			for {
				c, err := r.next()
				if err != nil {
					return
				}
				if isWhitespace(c) || isTerminatingMacro(c) {
					r.unread()
					return
				}
			}
		}
	default:
		// Atom (symbol, number, keyword, etc.) — skip to whitespace/delimiter
		for {
			c, err := r.next()
			if err != nil {
				return
			}
			if isWhitespace(c) || isTerminatingMacro(c) {
				r.unread()
				return
			}
		}
	}
}

func readMeta(r *LispReader, _ rune) (vm.Value, error) {
	ch, err := r.next()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading meta")
	}
	var m vm.Value = vm.EmptyPersistentMap
	for {
		switch ch {
		case '{':
			m2, err := readMap(r, ch)
			if err != nil {
				return vm.NIL, NewReaderError(r, "reading meta map")
			}
			if sq, ok := m2.(vm.Sequable); ok {
				for s := sq.Seq(); s != nil && s != vm.EmptyList; s = s.Next() {
					entry := s.First().(vm.ArrayVector)
					m = m.(*vm.PersistentMap).Assoc(entry[0], entry[1]).(*vm.PersistentMap)
				}
			}
		case ':':
			k, err := readToken(r, ':')
			if err != nil {
				return vm.NIL, NewReaderError(r, "reading meta keyword")
			}
			k, err = interpretToken(r, k)
			if err != nil {
				return vm.NIL, NewReaderError(r, "reading meta keyword")
			}
			m = m.(*vm.PersistentMap).Assoc(k, vm.TRUE).(*vm.PersistentMap)
		default:
			return vm.NIL, NewReaderError(r, "unsupported meta form")
		}
		ch, err = r.eatWhitespace()
		if err != nil {
			return vm.NIL, NewReaderError(r, "reading meta")
		}
		if ch != '^' {
			break
		}
		ch, err = r.next()
		if err != nil {
			return vm.NIL, NewReaderError(r, "reading meta")
		}
	}
	err = r.unread()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading meta")
	}
	form, err := r.Read()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading meta")
	}
	return vm.NewList([]vm.Value{vm.Symbol("with-meta"), form, m}), nil
}

func readSymbolicValue(r *LispReader, _ rune) (vm.Value, error) {
	ch, err := r.next()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading symbolic value")
	}
	token, err := readToken(r, ch)
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading symbolic value")
	}
	s := string(token.(vm.Symbol))
	switch s {
	case "Inf":
		return vm.Float(math.Inf(1)), nil
	case "-Inf":
		return vm.Float(math.Inf(-1)), nil
	case "NaN":
		return vm.Float(math.NaN()), nil
	default:
		return vm.NIL, NewReaderError(r, fmt.Sprintf("unknown symbolic value: ##%s", s))
	}
}

func readHashMacro(r *LispReader, _ rune) (vm.Value, error) {
	ch, err := r.next()
	if err != nil {
		return vm.NIL, NewReaderError(r, "reading hash macro")
	}
	macro, ok := hashMacros[ch]
	if !ok {
		return vm.NIL, NewReaderError(r, "invalid hash macro")
	}
	return macro(r, ch)
}

func unmatchedDelimReader(ru rune) readerFunc {
	return func(r *LispReader, _ rune) (vm.Value, error) {
		return nil, NewReaderError(r, fmt.Sprintf("unmatched delimiter %c", ru))
	}
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r) || r == ','
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
var hashMacros map[rune]readerFunc

// readerInit must be called in compiler package init before everything else !
func readerInit() {
	macros = map[rune]readerFunc{
		'(':  readList,
		')':  unmatchedDelimReader(')'),
		'[':  readVector,
		']':  unmatchedDelimReader(']'),
		'{':  readMap,
		'}':  unmatchedDelimReader('}'),
		'"':  readString,
		'\\': readChar,
		'\'': readQuote,
		'@':  readDeref,
		'`':  readSyntaxQuote,
		'~':  readUnquote,
		';':  readLineComment,
		'#':  readHashMacro,
		'^':  readMeta,
	}

	hashMacros = map[rune]readerFunc{
		'\'': readVarQuote,
		'_':  readFormComment,
		'(':  readShortFn,
		'{':  readSet,
		'"':  readRegex,
		'?':  readConditional,
		'#':  readSymbolicValue,
	}
}

func readLineComment(r *LispReader, _ rune) (vm.Value, error) {
	for {
		ch, err := r.next()
		if err == io.EOF || ch == '\n' || ch == '\r' {
			return vm.VOID, nil
		}
		if err != nil {
			return vm.NIL, NewReaderError(r, "unexpected error while reading line comment").Wrap(err)
		}
	}
}
func readFormComment(r *LispReader, _ rune) (vm.Value, error) {
	_, err := r.Read()
	if errors.IsCausedBy(err, io.EOF) {
		return vm.NIL, err
	}
	return vm.VOID, nil
}
