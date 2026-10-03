/*
 * Copyright (c) 2026 let-go contributors; see CONTRIBUTORS.
 * SPDX-License-Identifier: MIT
 */

package e2e

import (
	"os/exec"
	"strings"
	"testing"
)

// A core namespace the bundle loads on demand runs its lowered Go under
// -tags gogen_ir, like clojure.core. The bundle decodes the namespace's vars
// before the lowered package registers its overrides, and the namespace's own
// chunk replays on the first require, so the overrides must outlast that
// replay.
func TestGogenLazyCoreNamespaceStaysLowered(t *testing.T) {
	if testing.Short() {
		t.Skip("builds lg with -tags gogen_ir; run without -short")
	}
	bin := buildLGTags(t, "gogen_ir")
	out, err := exec.Command(bin, "-e",
		`(do (require 'clojure.walk)
		     (println (str frequencies))
		     (println (str clojure.walk/postwalk)))`).CombinedOutput()
	if err != nil {
		t.Fatalf("lg -e: %v\n%s", err, out)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var core, walk string
	for _, l := range lines {
		switch {
		case strings.Contains(l, " frequencies "):
			core = l
		case strings.Contains(l, " postwalk "):
			walk = l
		}
	}
	if !strings.HasPrefix(core, "<native-fn") {
		t.Fatalf("clojure.core/frequencies is %q under gogen_ir; is pkg/rt/core_go_lowered generated (make lowered)?\n%s", core, out)
	}
	if !strings.HasPrefix(walk, "<native-fn") {
		t.Fatalf("clojure.walk/postwalk is %q after require, want its lowered Go\n%s", walk, out)
	}
}
