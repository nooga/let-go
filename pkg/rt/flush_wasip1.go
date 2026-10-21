//go:build wasip1

/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"errors"
	"syscall"

	"github.com/nooga/let-go/pkg/vm"
)

// flushOut syncs the active *out* binding for clojure.core/flush. wasip1 has
// no term namespace; a rebound buffered *out* still needs its Flush.
//
// A WASI host may refuse fd_sync on stdio: Node's WASI returns ENOTCAPABLE when
// stdout is a pipe. Flushing such an fd is a no-op, so for a file-backed *out*
// those errnos are swallowed, as the native flushOut does for ENOTTY and
// friends. An embedder writer's Sync goes through Flush() and keeps its errors.
func flushOut(ec *vm.ExecContext) error {
	h := resolveIOHandleVar(ec, "*out*")
	if h == nil {
		return nil
	}
	err := h.Sync()
	if h.File() != nil && (errors.Is(err, syscall.ENOTCAPABLE) ||
		errors.Is(err, syscall.ENOTTY) ||
		errors.Is(err, syscall.EBADF) ||
		errors.Is(err, syscall.EINVAL) ||
		errors.Is(err, syscall.ENOTSUP)) {
		err = nil
	}
	return err
}
