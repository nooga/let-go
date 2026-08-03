/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

// This file exists so that runInstallers() fires after every other init()
// in the package has registered. Go runs init() functions in source-name
// order; the "zz_" prefix guarantees this file sorts last, so by the time
// its init() runs every other rt/*.go::init has finished appending to the
// installers slice (see installers.go).
//
// Why not call runInstallers from lang.go::init: lang.go sorts in the
// middle alphabetically, so files like math.go, os.go, pods.go,
// system.go, term.go, etc. wouldn't have registered yet when lang.go's
// init fired.

func init() {
	runInstallers()
	// Builtins register AFTER the installer drain so builtin versions (better
	// signatures) override the generated primitives for overlapping names
	// (vector, cons, contains?, array-map, nth). This preserves the EPIC-012
	// seam ordering now that the core primitive registrar self-registers via
	// the installer queue (drained just above) instead of an explicit call.
	registerBuiltinsModule()
}
