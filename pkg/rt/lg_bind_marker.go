/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

// pkg/rt OWNS clojure.core: its generated registrar must BIND vars, not just
// register direct-call metadata. This marker selects own mode in lgprimgen.
//lg:bind

package rt
