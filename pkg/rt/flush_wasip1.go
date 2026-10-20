//go:build wasip1

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import "github.com/nooga/let-go/pkg/vm"

// flushOut syncs the active *out* binding for clojure.core/flush. wasip1 has
// no term namespace; a rebound buffered *out* still needs its Flush.
func flushOut(ec *vm.ExecContext) error {
	if h := resolveIOHandleVar(ec, "*out*"); h != nil {
		return h.Sync()
	}
	return nil
}
