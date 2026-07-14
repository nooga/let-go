/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/nooga/let-go/pkg/vm"
)

var reflectionWarnings = struct {
	sync.Mutex
	writer io.Writer
	seen   map[string]struct{}
}{writer: os.Stderr, seen: make(map[string]struct{})}

// SetReflectionWarningWriter replaces the warning sink and returns a restore
// function. It is primarily useful to embedders and tests.
func SetReflectionWarningWriter(writer io.Writer) func() {
	reflectionWarnings.Lock()
	previous := reflectionWarnings.writer
	reflectionWarnings.writer = writer
	reflectionWarnings.Unlock()
	return func() {
		reflectionWarnings.Lock()
		reflectionWarnings.writer = previous
		reflectionWarnings.Unlock()
	}
}

// ResetReflectionWarnings clears the once-per-source-site deduplication set.
func ResetReflectionWarnings() {
	reflectionWarnings.Lock()
	clear(reflectionWarnings.seen)
	reflectionWarnings.Unlock()
}

func reflectionWarningsEnabled(ec *vm.ExecContext) bool {
	if CoreNS == nil {
		return false
	}
	warnVar := CoreNS.LookupLocal(vm.Symbol("*warn-on-reflection*"))
	if warnVar == nil {
		return false
	}
	return vm.IsTruthy(ec.Deref(warnVar))
}

func emitReflectionWarning(ec *vm.ExecContext, info *vm.SourceInfo, kind, reason string) {
	if info == nil || !reflectionWarningsEnabled(ec) {
		return
	}
	key := kind + "\x00" + info.String()
	reflectionWarnings.Lock()
	defer reflectionWarnings.Unlock()
	if _, exists := reflectionWarnings.seen[key]; exists {
		return
	}
	reflectionWarnings.seen[key] = struct{}{}
	fmt.Fprintf(reflectionWarnings.writer, "reflection warning, %s: %s (%s)\n", info.String(), reason, kind)
}

// EmitReflectionWarningForForm emits a warning using a compiler form's source.
func EmitReflectionWarningForForm(form vm.Value, kind, reason string) {
	emitReflectionWarning(vm.RootExecContext, vm.FormSource.Get(form), kind, reason)
}

func sourceInfoFromWarningValue(value vm.Value) *vm.SourceInfo {
	if boxed, ok := value.(*vm.Boxed); ok {
		switch info := boxed.Unbox().(type) {
		case vm.SourceInfo:
			return &info
		case *vm.SourceInfo:
			return info
		}
	}
	seq, err := seqOf(value)
	if err != nil {
		return nil
	}
	for current := seq; current != nil; current = current.Next() {
		if info := sourceInfoFromWarningValue(current.First()); info != nil && info.File != "" {
			return info
		}
	}
	return nil
}

func emitReflectionWarningValue(ec *vm.ExecContext, source vm.Value, kind, reason string) {
	emitReflectionWarning(ec, sourceInfoFromWarningValue(source), kind, reason)
}
