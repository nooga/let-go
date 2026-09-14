/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * SPDX-License-Identifier: MIT
 */

package compiler

import (
	"errors"
	"fmt"
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
