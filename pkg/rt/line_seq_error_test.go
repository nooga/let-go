package rt

import (
	"bufio"
	"errors"
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// failAfter returns data on the first Read and err on every later Read, like
// a connection that delivers one chunk and then resets.
type failAfter struct {
	data string
	err  error
	sent bool
}

func (r *failAfter) Read(p []byte) (int, error) {
	if !r.sent {
		r.sent = true
		return copy(p, r.data), nil
	}
	return 0, r.err
}

// TestLineSeqPropagatesReadErrors: only io.EOF ends a line-seq. Any other read
// error is thrown when the sequence is realized, so a consumer can tell a
// finished stream from a failed one.
func TestLineSeqPropagatesReadErrors(t *testing.T) {
	ec := vm.NewExecContext()
	first := LookupVar("core", "first").Deref().(vm.Fn)
	rest := LookupVar("core", "rest").Deref().(vm.Fn)
	nth := func(t *testing.T, s vm.Value, n int) (vm.Value, error) {
		t.Helper()
		for i := 0; i < n; i++ {
			next, err := ec.Invoke(rest, []vm.Value{s})
			if err != nil {
				return nil, err
			}
			s = next
		}
		return ec.Invoke(first, []vm.Value{s})
	}

	t.Run("a read error is thrown, not treated as end of file", func(t *testing.T) {
		s := makeLineSeq(bufio.NewReader(&failAfter{data: "first\n", err: errors.New("connection reset by peer")}))
		head, err := nth(t, s, 0)
		if err != nil || head.String() != `"first"` {
			t.Fatalf("first line: %v %v", head, err)
		}
		_, err = nth(t, s, 1)
		if err == nil || !strings.Contains(err.Error(), "connection reset by peer") {
			t.Fatalf("line-seq must surface the read error, got %v", err)
		}
	})

	t.Run("end of file still ends the sequence and keeps an unterminated last line", func(t *testing.T) {
		s := makeLineSeq(bufio.NewReader(strings.NewReader("a\nb")))
		for i, want := range []string{`"a"`, `"b"`} {
			v, err := nth(t, s, i)
			if err != nil || v.String() != want {
				t.Fatalf("line %d: %v %v", i, v, err)
			}
		}
		v, err := nth(t, s, 2)
		if err != nil || v != vm.NIL {
			t.Fatalf("after the last line the sequence must end cleanly, got %v %v", v, err)
		}
	})
}
