/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"strings"
	"testing"

	"github.com/nooga/let-go/pkg/vm"
)

// Lowered code calls through CachedVarFn; a var that holds no function must
// surface as an error from the call, as on the VM, not as a Go panic.
func TestCachedVarFnReportsANonFunctionAsAnError(t *testing.T) {
	ns := NS("test-cached-var-fn")
	ns.Def("number", vm.Int(1))
	ns.LookupOrAdd(vm.Symbol("unbound")) // interned with no root, as (declare unbound) leaves it

	for _, tc := range []struct {
		sym, want string
	}{
		{"unbound", "Attempting to call unbound fn: #'test-cached-var-fn/unbound"},
		{"number", "is not a function"},
		{"missing", "Unable to resolve var: test-cached-var-fn/missing"},
	} {
		t.Run(tc.sym, func(t *testing.T) {
			var ptr *vm.Var
			var err error
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("CachedVarFn panicked: %v", r)
					}
				}()
				_, err = CachedVarFn(&ptr, "test-cached-var-fn", tc.sym).Invoke(nil)
			}()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("calling #'test-cached-var-fn/%s: got error %v, want one containing %q", tc.sym, err, tc.want)
			}
		})
	}
}
