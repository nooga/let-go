/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// alter-var-root writes back a wrapper the watch already made. That must not
// wrap it again, or ProgramEntryRan would restore the outer wrapper instead
// of the entry's own fn.
func TestProgramEntryWatchRestoresTheEntryAfterAWriteBack(t *testing.T) {
	ns := NS("entrywatchtest")
	entry := vm.NewCtxNativeFn("-main", func(*vm.ExecContext, []vm.Value) (vm.Value, error) { return vm.NIL, nil })
	WatchProgramEntry("entrywatchtest", "-main")
	v := ns.Def("-main", entry)
	if v.Root() == vm.Value(entry) {
		t.Fatal("the entry wasn't wrapped while watched")
	}
	identity := vm.NewCtxNativeFn("identity", func(_ *vm.ExecContext, vs []vm.Value) (vm.Value, error) { return vs[0], nil })
	if _, err := v.AlterRoot(identity); err != nil {
		t.Fatal(err)
	}
	if ProgramEntryRan("entrywatchtest", "-main") {
		t.Error("ProgramEntryRan is true, but nothing called the entry")
	}
	if v.Root() != vm.Value(entry) {
		t.Errorf("after the watch the var holds %v, want the entry's own fn", v.Root())
	}
}
