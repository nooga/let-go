/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

// installers is the queue of namespace-installation functions to run from
// lang.go's init(). Files that contribute a Go-side namespace add themselves
// via their own init():
//
//	func init() { RegisterInstaller(installFooNS) }
//	func installFooNS() {
//	    ns := vm.NewNamespace("foo")
//	    ns.Def(...)
//	    RegisterNS(ns)
//	}
//
// This removes the need to maintain an explicit list of install calls in
// lang.go::init — adding a binding is one new file with one init() line.
//
// Ordering: Go evaluates init() functions in source-filename order within
// a package, so the slice is appended in deterministic order per build. No
// installer in pkg/rt today depends on another running first; each
// constructs its own Namespace and registers it independently. If that ever
// changes, treat the cross-installer dependency as a `require` at the
// resolver layer, not as init-order coupling.
var installers []func()

// RegisterInstaller queues fn to be called from lang.go's init(). Call this
// from your file's init() — it is not safe to call after init has run.
func RegisterInstaller(fn func()) {
	installers = append(installers, fn)
}

// runInstallers invokes every queued installer in registration order.
// Called once from lang.go::init.
func runInstallers() {
	for _, fn := range installers {
		fn()
	}
}
