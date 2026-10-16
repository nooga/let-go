/*
 * Copyright (c) 2026 let-go contributors
 * SPDX-License-Identifier: MIT
 */

package rt

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
)

func unregisterEmbeddedSource(name string) {
	auxEmbeddedMu.Lock()
	delete(auxEmbeddedSources, name)
	auxEmbeddedMu.Unlock()
}

func TestRegisterEmbeddedSource(t *testing.T) {
	const name = "test.register-embedded-source"
	t.Cleanup(func() { unregisterEmbeddedSource(name) })
	RegisterEmbeddedSource(name, "(ns test.register-embedded-source)")
	if src, ok := EmbeddedSource(name); !ok || src != "(ns test.register-embedded-source)" {
		t.Fatalf("EmbeddedSource(%q) = %q, %v", name, src, ok)
	}
	// Registered sources stay out of the core lowering universe (#1063).
	if slices.Contains(EmbeddedNSNames(), name) {
		t.Fatalf("EmbeddedNSNames lists %q; registered sources must not be lowered with core", name)
	}
}

func TestRegisterEmbeddedSourceRefuses(t *testing.T) {
	const dup = "test.register-embedded-source-dup"
	t.Cleanup(func() { unregisterEmbeddedSource(dup) })
	RegisterEmbeddedSource(dup, "(ns test.register-embedded-source-dup)")

	for _, tc := range []struct{ desc, name, src, want string }{
		{"empty source", "test.register-embedded-source-empty", "", "is empty"},
		{"core namespace", "walk", "(ns walk)", "would shadow the core namespace"},
		{"native namespace", "os", "(ns os)", "native namespace of that name"},
		{"second registration", dup, "(ns test.register-embedded-source-dup)", "registered twice"},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil || !strings.Contains(fmt.Sprint(r), tc.want) || !strings.Contains(fmt.Sprint(r), tc.name) {
					t.Fatalf("panic = %v, want one naming %q and containing %q", r, tc.name, tc.want)
				}
			}()
			RegisterEmbeddedSource(tc.name, tc.src)
		})
	}
	for _, n := range []string{"walk", "os"} {
		if _, shadowed := auxEmbeddedSources[n]; shadowed {
			t.Fatalf("the refused registration was recorded over %s", n)
		}
	}
}

// Registration may run after the resolver has started reading (pkg/compiler's
// init reads embedded sources), so it must be safe alongside reads. Run with
// -race to make this meaningful.
func TestRegisterEmbeddedSourceConcurrentWithReads(t *testing.T) {
	var names []string
	for i := range 8 {
		names = append(names, fmt.Sprintf("test.register-embedded-source-race-%d", i))
	}
	t.Cleanup(func() {
		for _, n := range names {
			unregisterEmbeddedSource(n)
		}
	})
	var wg sync.WaitGroup
	for _, n := range names {
		wg.Go(func() { RegisterEmbeddedSource(n, "(ns "+n+")") })
		wg.Go(func() { EmbeddedSource(n); EmbeddedSource("walk") })
	}
	wg.Wait()
	for _, n := range names {
		if _, ok := EmbeddedSource(n); !ok {
			t.Errorf("%s not registered", n)
		}
	}
}
