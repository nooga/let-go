package vm

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

// oldStringMethod is .charAt/.indexOf/.substring as they were before the
// in-place walk: every call round-tripped through []rune, which also turned
// invalid UTF-8 bytes into U+FFFD. The current methods must agree with it on
// every input, malformed ones included.
func oldStringMethod(l String, name string, args []Value) (Value, error) {
	switch name {
	case "charAt":
		i := args[0].(Int)
		rs := []rune(string(l))
		if i < 0 || int(i) >= len(rs) {
			return NIL, fmt.Errorf("string.charAt: index %d out of bounds for length %d", int(i), len(rs))
		}
		return Char(rs[i]), nil
	case "indexOf":
		rs := []rune(string(l))
		from := 0
		if len(args) == 2 {
			from = min(max(int(args[1].(Int)), 0), len(rs))
		}
		needle := string(args[0].(String))
		tail := string(rs[from:])
		idx := strings.Index(tail, needle)
		if idx < 0 {
			return Int(-1), nil
		}
		return Int(from + utf8.RuneCountInString(tail[:idx])), nil
	case "substring":
		rs := []rune(string(l))
		begin := args[0].(Int)
		end := Int(len(rs))
		if len(args) == 2 {
			end = args[1].(Int)
		}
		if begin < 0 || begin > end || int(end) > len(rs) {
			return NIL, fmt.Errorf("string.substring: range [%d, %d) out of bounds for length %d", int(begin), int(end), len(rs))
		}
		return String(rs[begin:end]), nil
	}
	panic(name)
}

func TestStringMethodsMatchRuneSliceSemantics(t *testing.T) {
	corpus := []String{
		"", "a", "abc", "héllo", "日本語", "a\xffb", "\xff", "ab\xe2\x82cd", "€\x82€",
		"a�b", "\xc3", "x\xe2\x82\xacy\xf0\x9f", String(strings.Repeat("é—", 5) + "\xfe"),
	}
	needles := []String{"", "a", "b", "é", "€", "�", "\xff", "\x82", "\xe2\x82", "ab"}
	same := func(t *testing.T, call string, got, want Value, gotErr, wantErr error) {
		t.Helper()
		if fmt.Sprint(gotErr) != fmt.Sprint(wantErr) || (wantErr == nil && !(got.Type() == want.Type() && fmt.Sprintf("%q", got) == fmt.Sprintf("%q", want))) {
			t.Errorf("%s = %#v, %v; want %#v, %v", call, got, gotErr, want, wantErr)
		}
	}
	for _, s := range corpus {
		n := utf8.RuneCountInString(string(s))
		for i := -2; i <= n+2; i++ {
			args := []Value{Int(i)}
			got, gotErr := s.InvokeMethod("charAt", args)
			want, wantErr := oldStringMethod(s, "charAt", args)
			same(t, fmt.Sprintf("%q.charAt(%d)", string(s), i), got, want, gotErr, wantErr)

			got, gotErr = s.InvokeMethod("substring", args)
			want, wantErr = oldStringMethod(s, "substring", args)
			same(t, fmt.Sprintf("%q.substring(%d)", string(s), i), got, want, gotErr, wantErr)
			for j := -2; j <= n+2; j++ {
				args := []Value{Int(i), Int(j)}
				got, gotErr := s.InvokeMethod("substring", args)
				want, wantErr := oldStringMethod(s, "substring", args)
				same(t, fmt.Sprintf("%q.substring(%d, %d)", string(s), i, j), got, want, gotErr, wantErr)
			}
			for _, nd := range needles {
				for _, args := range [][]Value{{nd}, {nd, Int(i)}} {
					got, gotErr := s.InvokeMethod("indexOf", args)
					want, wantErr := oldStringMethod(s, "indexOf", args)
					same(t, fmt.Sprintf("%q.indexOf%q", string(s), fmt.Sprint(args)), got, want, gotErr, wantErr)
				}
			}
		}
	}
}
