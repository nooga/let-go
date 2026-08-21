/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	"github.com/nooga/let-go/pkg/vm"
)

func TestTaggedReaderRegistryDispatchesDataReader(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	var called atomic.Bool
	if err := registry.RegisterData("probe", func(value vm.Value) (vm.Value, error) {
		called.Store(true)
		if value.String() != "[1 2]" {
			t.Fatalf("data reader received %v, want [1 2]", value)
		}
		return vm.String("data-reader-result"), nil
	}); err != nil {
		t.Fatal(err)
	}

	reader := NewLispReaderWithTaggedReaders(strings.NewReader("#probe ; ignored\n [1 2]"), "probe.lg", registry)
	got, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if !called.Load() || got != vm.String("data-reader-result") {
		t.Fatalf("data reader dispatch: called=%v got=%v", called.Load(), got)
	}
}

func TestTaggedReaderRegistryRejectsDuplicateAndUnknownTags(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterData("same", func(value vm.Value) (vm.Value, error) { return value, nil }); err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterData("same", func(value vm.Value) (vm.Value, error) { return value, nil }); err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Fatalf("duplicate registration error = %v", err)
	}
	reader := NewLispReaderWithTaggedReaders(strings.NewReader("#missing 1"), "probe.lg", registry)
	if _, err := reader.Read(); err == nil || !strings.Contains(err.Error(), "unknown tagged literal #missing") {
		t.Fatalf("unknown tag error = %v", err)
	}
}

func TestTaggedReaderRegistryConcurrentRegistrationHasOneWinner(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	var successes atomic.Int32
	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if registry.RegisterData("shared", func(value vm.Value) (vm.Value, error) { return value, nil }) == nil {
				successes.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := successes.Load(); got != 1 {
		t.Fatalf("successful registrations = %d, want 1", got)
	}
}

func TestTaggedReaderRegistryValidatesDispatchableTags(t *testing.T) {
	var registry TaggedReaderRegistry
	identity := func(value vm.Value) (vm.Value, error) { return value, nil }
	for _, tag := range []string{"tools/go", "uuid", "inst"} {
		if err := registry.RegisterData(tag, identity); err != nil {
			t.Fatalf("register %q: %v", tag, err)
		}
	}
	for _, tag := range []string{"", "1x", "a b", "go{", "é"} {
		if err := registry.RegisterData(tag, identity); err == nil || !strings.Contains(err.Error(), "cannot be read") {
			t.Fatalf("unreachable tag %q error = %v", tag, err)
		}
	}
}

func TestTaggedReaderRegistryOverridesBuiltinTags(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	for _, tag := range []string{"uuid", "inst"} {
		tag := tag
		if err := registry.RegisterData(tag, func(vm.Value) (vm.Value, error) {
			return vm.String("custom-" + tag), nil
		}); err != nil {
			t.Fatal(err)
		}
	}
	for _, tag := range []string{"uuid", "inst"} {
		reader := NewLispReaderWithTaggedReaders(strings.NewReader("#"+tag+" \"invalid\""), "probe.lg", registry)
		got, err := reader.Read()
		if err != nil {
			t.Fatalf("#%s override: %v", tag, err)
		}
		if got != vm.String("custom-"+tag) {
			t.Fatalf("#%s override = %v", tag, got)
		}
	}
}

func TestTaggedReaderRegistryPropagatesHandlerAndPayloadErrors(t *testing.T) {
	if err := errors.Join(errors.New("other"), io.EOF); !IsErrorEOF(err) {
		t.Fatalf("IsErrorEOF(errors.Join(other, EOF)) = false: %v", err)
	}

	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterData("fail", func(vm.Value) (vm.Value, error) {
		return vm.NIL, errors.New("handler failed")
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := NewLispReaderWithTaggedReaders(strings.NewReader("#fail 1"), "probe.lg", registry).Read(); err == nil || !strings.Contains(err.Error(), "handler failed") {
		t.Fatalf("handler error = %v", err)
	}

	if err := registry.RegisterData("need-value", func(value vm.Value) (vm.Value, error) { return value, nil }); err != nil {
		t.Fatal(err)
	}
	_, err := NewLispReaderWithTaggedReaders(strings.NewReader("#need-value"), "probe.lg", registry).Read()
	if err == nil || !IsErrorEOF(err) || !strings.Contains(err.Error(), "unexpected EOF") {
		t.Fatalf("truncated tagged literal error = %v, want EOF syntax error", err)
	}

	if err := registry.RegisterData("wrapped-eof", func(vm.Value) (vm.Value, error) {
		return vm.NIL, fmt.Errorf("handler needs more input: %w", io.EOF)
	}); err != nil {
		t.Fatal(err)
	}
	_, err = NewLispReaderWithTaggedReaders(strings.NewReader("#wrapped-eof 1"), "probe.lg", registry).Read()
	if err == nil || !IsErrorEOF(err) {
		t.Fatalf("wrapped handler EOF error = %v, want EOF causality", err)
	}
	if !strings.Contains(err.Error(), "handler needs more input") {
		t.Fatalf("wrapped handler EOF lost handler message: %v", err)
	}

	if err := registry.RegisterData("joined-eof", func(vm.Value) (vm.Value, error) {
		return vm.NIL, errors.Join(errors.New("joined handler failed"), io.EOF)
	}); err != nil {
		t.Fatal(err)
	}
	_, err = NewLispReaderWithTaggedReaders(strings.NewReader("#joined-eof 1"), "probe.lg", registry).Read()
	if err == nil || !IsErrorEOF(err) {
		t.Fatalf("joined handler EOF error = %v, want EOF causality", err)
	}
	if !strings.Contains(err.Error(), "joined handler failed") {
		t.Fatalf("joined handler EOF lost handler message: %v", err)
	}
}

func compileWithTaggedReaders(t *testing.T, registry *TaggedReaderRegistry, src string) (vm.Value, error) {
	t.Helper()
	ctx := NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS)).SetTaggedReaders(registry)
	_, got, err := ctx.CompileMultiple(strings.NewReader(src))
	return got, err
}

func TestCompilerTaggedReadersReachRuntimeStringReaders(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterData("scoped", func(vm.Value) (vm.Value, error) {
		return vm.Keyword("explicit"), nil
	}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		src  string
		want string
	}{
		{name: "read-string", src: `(read-string "#scoped 1")`, want: ":explicit"},
		{name: "read-all-string", src: `(read-all-string "#scoped 1 #scoped 2")`, want: "[:explicit :explicit]"},
		{name: "load-string", src: `(load-string "#scoped 3")`, want: ":explicit"},
		{name: "eval", src: `(eval '(read-string "#scoped 4"))`, want: ":explicit"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compileWithTaggedReaders(t, registry, tt.src)
			if err != nil {
				t.Fatal(err)
			}
			if got.String() != tt.want {
				t.Fatalf("result = %s, want %s", got.String(), tt.want)
			}
		})
	}
}

func TestCompilerTaggedReadersTakePrecedenceDuringRuntimeReads(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	for _, tag := range []string{"scoped", "uuid"} {
		if err := registry.RegisterData(tag, func(vm.Value) (vm.Value, error) {
			return vm.Keyword("explicit"), nil
		}); err != nil {
			t.Fatal(err)
		}
	}

	for _, tt := range []struct {
		src  string
		want string
	}{
		{src: `(binding [*data-readers* {'scoped (fn [_] :dynamic)}] (read-string "#scoped 1"))`, want: ":explicit"},
		{src: `(binding [*data-readers* {'scoped (fn [_] :dynamic)}] (read-all-string "#scoped 1"))`, want: "[:explicit]"},
		{src: `(binding [*data-readers* {'scoped (fn [_] :dynamic)}] (load-string "#scoped 1"))`, want: ":explicit"},
		{src: `(binding [*data-readers* {'scoped (fn [_] :dynamic)}] (eval '(read-string "#scoped 1")))`, want: ":explicit"},
		{src: `(binding [*data-readers* {'uuid (fn [_] :dynamic)}] (read-string "#uuid \"not-a-uuid\""))`, want: ":explicit"},
		{src: `(binding [*data-readers* {'uuid (fn [_] :dynamic)}] (eval '(read-string "#uuid \"not-a-uuid\"")))`, want: ":explicit"},
	} {
		got, err := compileWithTaggedReaders(t, registry, tt.src)
		if err != nil {
			t.Fatalf("%s: %v", tt.src, err)
		}
		if got.String() != tt.want {
			t.Fatalf("%s result = %v, want %s", tt.src, got, tt.want)
		}
	}
}

func TestDataReadersDynamicBindingMustRemainAMap(t *testing.T) {
	_, err := compileWithTaggedReaders(t, nil, `(binding [*data-readers* []] (read-string "#scoped 1"))`)
	if err == nil || !strings.Contains(err.Error(), "*data-readers* must be a map") {
		t.Fatalf("non-map *data-readers* error = %v", err)
	}
}

func TestCompilerTaggedReadersDoNotLeakBetweenCompilers(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterData("scoped", func(vm.Value) (vm.Value, error) {
		return vm.Keyword("explicit"), nil
	}); err != nil {
		t.Fatal(err)
	}

	got, err := compileWithTaggedReaders(t, registry, `(read-string "#scoped 1")`)
	if err != nil || got != vm.Keyword("explicit") {
		t.Fatalf("explicit compiler result = %v, err = %v", got, err)
	}
	got, err = compileWithTaggedReaders(t, nil, `(read-string "#scoped 1")`)
	if err != nil || got != vm.Int(1) {
		t.Fatalf("default compiler result = %v, err = %v; explicit registry leaked", got, err)
	}
	got, err = compileWithTaggedReaders(t, registry, `(eval '(read-string "#scoped 2"))`)
	if err != nil || got != vm.Keyword("explicit") {
		t.Fatalf("explicit compiler eval result = %v, err = %v", got, err)
	}
	got, err = compileWithTaggedReaders(t, nil, `(eval '(read-string "#scoped 2"))`)
	if err != nil || got != vm.Int(2) {
		t.Fatalf("default compiler eval result = %v, err = %v; explicit registry leaked", got, err)
	}
}

func TestCompilerTaggedReadersAreConcurrentEvaluationLocal(t *testing.T) {
	const workers = 16
	type result struct {
		want string
		got  vm.Value
		err  error
	}
	contexts := make([]*Context, workers)
	sources := make([]string, workers)
	wants := make([]string, workers)
	for i := range workers {
		want := vm.Keyword(fmt.Sprintf("registry-%d", i))
		registry := NewTaggedReaderRegistry()
		if err := registry.RegisterData("scoped", func(vm.Value) (vm.Value, error) {
			return want, nil
		}); err != nil {
			t.Fatal(err)
		}
		contexts[i] = NewCompiler(vm.NewConsts(), rt.NS(rt.NameCoreNS)).SetTaggedReaders(registry)
		switch i % 4 {
		case 0:
			sources[i] = `(read-string "#scoped 1")`
			wants[i] = want.String()
		case 1:
			sources[i] = `(read-all-string "#scoped 1")`
			wants[i] = "[" + want.String() + "]"
		case 2:
			sources[i] = `(load-string "#scoped 1")`
			wants[i] = want.String()
		case 3:
			sources[i] = `(eval '(read-string "#scoped 1"))`
			wants[i] = want.String()
		}
	}

	results := make(chan result, workers)
	var wg sync.WaitGroup
	for i := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, got, err := contexts[i].CompileMultiple(strings.NewReader(sources[i]))
			results <- result{want: wants[i], got: got, err: err}
		}()
	}
	wg.Wait()
	close(results)
	for result := range results {
		if result.err != nil || result.got.String() != result.want {
			t.Errorf("concurrent result = %v, want %v, err = %v", result.got, result.want, result.err)
		}
	}
}

func TestLegacyReaderKeepsUnknownTagBestEffortBehavior(t *testing.T) {
	reader := NewLispReader(strings.NewReader("#probe 42"), "probe.lg")
	got, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if got != vm.Int(42) {
		t.Fatalf("legacy unknown tag result = %v, want 42", got)
	}
}

func TestRawGoReaderCapturesLexicallyBalancedFragment(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterRaw("go", ReadRawGoFragment); err != nil {
		t.Fatal(err)
	}
	input := "#go{" + `
if ready {
    println("brace } and escaped quote \"")
    println(` + "`raw { brace }`" + `)
    r := '}'
    // comment closes } and opens {
    /* block { comment } */
    _ = r
}
` + "} 42"
	reader := NewLispReaderWithTaggedReaders(strings.NewReader(input), "probe.lg", registry)
	got, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	want := input[len("#go{") : len(input)-len("} 42")]
	if got != vm.String(want) {
		t.Fatalf("raw fragment:\n got %q\nwant %q", got, want)
	}
	wrapped := "package probe\nfunc _() {" + string(got.(vm.String)) + "\n}"
	if _, err := parser.ParseFile(token.NewFileSet(), "expected.go", wrapped, parser.AllErrors); err != nil {
		t.Fatalf("captured #go fragment did not parse as Go statements: %v", err)
	}
	next, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if next != vm.Int(42) {
		t.Fatalf("form after raw fragment = %v, want 42", next)
	}
}

func TestRawGoReaderIsAvailableByDefault(t *testing.T) {
	reader := NewLispReader(strings.NewReader(`#go{if ready { return "}" }}`), "probe.lg")
	got, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if got != vm.String(`if ready { return "}" }`) {
		t.Fatalf("default #go result = %q", got)
	}
}

func TestRawGoReaderAllowsWhitespaceBeforeOpeningBrace(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterRaw("go", ReadRawGoFragment); err != nil {
		t.Fatal(err)
	}
	reader := NewLispReaderWithTaggedReaders(strings.NewReader("#go  { return 1 }"), "probe.lg", registry)
	got, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if got != vm.String(" return 1 ") {
		t.Fatalf("raw fragment = %q, want %q", got, " return 1 ")
	}
}

func TestRawGoReaderRejectsMissingOrUnterminatedFragment(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterRaw("go", ReadRawGoFragment); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, input, want string
	}{
		{"missing opener", "#go [1 2]", "requires an opening {"},
		{"unterminated", "#go{if x { println(x) }", "unterminated #go fragment"},
		{"unterminated string", "#go{println(\"x)}", "unterminated #go fragment"},
		{"unterminated block comment", "#go{/* comment }", "unterminated #go fragment"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := NewLispReaderWithTaggedReaders(strings.NewReader(tc.input), "probe.lg", registry)
			if _, err := reader.Read(); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want substring %q", err, tc.want)
			} else if readerErr, ok := err.(*ReaderError); ok && readerErr.IsEOF() {
				t.Fatalf("malformed #go was reported as clean EOF: %v", err)
			}
		})
	}
}

func TestRawGoReaderSkipsUnselectedReaderConditionalBranch(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterRaw("go", ReadRawGoFragment); err != nil {
		t.Fatal(err)
	}
	input := "#?(:clj #go{if ready { println(`brace }`) }} :default 7) 8"
	reader := NewLispReaderWithTaggedReaders(strings.NewReader(input), "probe.lg", registry)
	first, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	second, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if first != vm.Int(7) || second != vm.Int(8) {
		t.Fatalf("forms after skipped #go branch = %v, %v; want 7, 8", first, second)
	}
}

func TestExplicitDataGoReaderOverridesDefaultWhileSkippingConditional(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterData("go", func(value vm.Value) (vm.Value, error) { return value, nil }); err != nil {
		t.Fatal(err)
	}
	reader := NewLispReaderWithTaggedReaders(
		strings.NewReader("#?(:clj #go [1 2] :default 7) #go [3 4]"),
		"probe.lg",
		registry,
	)
	first, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	second, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}
	if first != vm.Int(7) || second.String() != "[3 4]" {
		t.Fatalf("explicit data #go results = %v, %v; want 7, [3 4]", first, second)
	}
}

// A tagged reader's own failure is a syntax error, never end-of-input. The
// EOF side of this contract is the #770 review's: a truncated tagged literal
// (`#need-value`) or a handler that reports EOF IS incomplete input, exactly
// like an unterminated `(defn`, so a REPL keeps prompting — see
// TestTaggedReaderRegistryPropagatesHandlerAndPayloadErrors. What must not
// happen is the converse: a non-EOF handler error, or a malformed payload on
// complete input, reading as EOF and silently ending a read.
func TestTaggedReaderErrorsCannotMasqueradeAsCleanEOF(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterData("fail", func(vm.Value) (vm.Value, error) {
		return vm.NIL, errors.New("handler failed")
	}); err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterRaw("raw", func(input TaggedRawInput) (vm.Value, error) {
		return vm.NIL, errors.New("raw payload rejected")
	}); err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{"#fail 1", "#raw{x}", "#go{"} {
		reader := NewLispReaderWithTaggedReaders(strings.NewReader(input), "probe.lg", registry)
		_, err := reader.Read()
		if err == nil {
			t.Fatalf("%q returned nil error", input)
		}
		if input != "#go{" && IsErrorEOF(err) {
			t.Fatalf("%q was reported as EOF: %v", input, err)
		}
	}
}

func TestRawGoReaderReportsMalformedSkippedConditionalBranch(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterRaw("go", ReadRawGoFragment); err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{
		"#?(:clj #go [1 2] :default 7)",
		"#?(:clj #go{if ready { println(1) } :default 7)",
	} {
		reader := NewLispReaderWithTaggedReaders(strings.NewReader(input), "probe.lg", registry)
		_, err := reader.Read()
		if err == nil {
			t.Fatalf("malformed skipped branch %q returned nil error", input)
		}
		if readerErr, ok := err.(*ReaderError); ok && readerErr.IsEOF() {
			t.Fatalf("malformed skipped branch %q was reported as clean EOF: %v", input, err)
		}
	}
}

func TestRawGoReaderRejectsEveryTruncatedPrefix(t *testing.T) {
	registry := NewTaggedReaderRegistry()
	if err := registry.RegisterRaw("go", ReadRawGoFragment); err != nil {
		t.Fatal(err)
	}
	full := `#go{if ready { println("}", ` + "`{`" + `, '\'') /* } */ }}`
	for cut := len("#go{"); cut < len(full); cut++ {
		reader := NewLispReaderWithTaggedReaders(strings.NewReader(full[:cut]), "probe.lg", registry)
		_, err := reader.Read()
		if err == nil {
			t.Fatalf("prefix %d/%d unexpectedly parsed: %q", cut, len(full), full[:cut])
		}
		if readerErr, ok := err.(*ReaderError); ok && readerErr.IsEOF() {
			t.Fatalf("prefix %d/%d was reported as clean EOF: %v", cut, len(full), err)
		}
	}
}
