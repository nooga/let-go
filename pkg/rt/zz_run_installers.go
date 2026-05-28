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
}
