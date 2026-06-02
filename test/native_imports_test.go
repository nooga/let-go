package test

// Blank-imports so the TestRunner binary links the native-direct modules,
// firing their init() registrations. corefns -> clojure.core/seq.
import (
	_ "github.com/nooga/let-go/pkg/rt/corefns"
)
