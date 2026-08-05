/*
 * Copyright (c) 2026 Norman Nunley, Jr <nnunley@gmail.com>
 * Part of the let-go project; see CONTRIBUTORS for full list of authors.
 * SPDX-License-Identifier: MIT
 */

package rt

import "testing"

// knownShadowedHandRegistrations is the accepted set of names registered both
// by hand in installLangNS and by the generated //lg:native registrar. Only
// the generated body runs; the closure is dead code that still reads as live.
//
// Each entry is a duplicate implementation awaiting a decision: port whatever
// the closure has that the //lg:native decl lacks, then delete the closure.
// The two bodies agree for the ones spot-checked so far (get, some), but that
// is a property nothing enforces — see TestReduceRangeFastPath's history.
var knownShadowedHandRegistrations = []string{
	"clojure.core/conj",
	"clojure.core/deref",
	"clojure.core/get",
	"clojure.core/int",
	"clojure.core/name",
	"clojure.core/namespace",
	"clojure.core/nth",
	"clojure.core/pop-binding!",
	"clojure.core/push-binding!",
	"clojure.core/some",
	"clojure.core/str",
	"clojure.core/subs",
}

// TestNoNewShadowedHandRegistrations ratchets the duplicate-registration set.
//
// A name in both places is a silent hazard, not a style problem: the generated
// registrar drains last, so the hand-written closure never runs, and the two
// bodies are free to drift. `reduce` drifted — the closure grew ArrayVector and
// Range fast paths the //lg:native decl never got — and when #639's ns-alias
// fix made the generated registration resolve for the first time, reduce lost
// those fast paths and got 1.75x slower with the whole suite still green,
// including a test named TestReduceRangeFastPath (it asserts results, not the
// path). This test is the guard that was missing: adding a //lg:native decl for
// an already-hand-registered name fails here instead of years later in a
// benchmark.
func TestNoNewShadowedHandRegistrations(t *testing.T) {
	known := map[string]bool{}
	for _, k := range knownShadowedHandRegistrations {
		known[k] = true
	}
	got := ShadowedHandRegistrations()
	seen := map[string]bool{}
	for _, g := range got {
		seen[g] = true
		if !known[g] {
			t.Errorf("new shadowed hand registration %s: the generated //lg:native body wins and the installLangNS closure is dead. Port anything the closure has that the decl lacks, delete the closure, or add the name to knownShadowedHandRegistrations with a reason.", g)
		}
	}
	for _, k := range knownShadowedHandRegistrations {
		if !seen[k] {
			t.Errorf("%s is no longer double-registered — drop it from knownShadowedHandRegistrations", k)
		}
	}
}
