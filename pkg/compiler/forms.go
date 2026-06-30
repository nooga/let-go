package compiler

import (
	"strings"

	"github.com/nooga/let-go/pkg/vm"
)

// TopLevelForm holds one top-level form plus the exact source slice that
// produced it.
type TopLevelForm struct {
	Form      vm.Value
	Source    string
	Start     int
	End       int
	Line      int
	Column    int
	EndLine   int
	EndColumn int
}

// SplitTopLevelForms reads src as a sequence of top-level forms and returns
// their exact source slices in order. Leading whitespace and no-value reader
// forms are skipped; each returned Source starts at the first byte of the
// corresponding real form.
func SplitTopLevelForms(src, inputName string) ([]TopLevelForm, error) {
	r := NewLispReader(strings.NewReader(src), inputName)
	var out []TopLevelForm
	for {
		ch, err := r.eatWhitespace()
		if err != nil {
			if isErrorEOF(err) {
				return out, nil
			}
			return nil, err
		}
		if err := r.unread(); err != nil {
			return nil, err
		}
		_ = ch
		startPos := r.pos
		startLine := r.line
		startCol := r.column

		form, err := r.Read()
		if err != nil {
			if isErrorEOF(err) {
				return out, nil
			}
			return nil, err
		}
		if form.Type() == vm.VoidType {
			continue
		}
		endPos := r.pos
		if startPos < 0 {
			startPos = 0
		}
		if endPos > len(src) {
			endPos = len(src)
		}
		out = append(out, TopLevelForm{
			Form:      form,
			Source:    src[startPos:endPos],
			Start:     startPos,
			End:       endPos,
			Line:      startLine,
			Column:    startCol,
			EndLine:   r.line,
			EndColumn: r.column,
		})
	}
}
