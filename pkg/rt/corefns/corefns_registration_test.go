package corefns_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/rt"
	_ "github.com/nooga/let-go/pkg/rt/corefns" // import triggers self-registration
)

func TestCorefnsSelfRegistersDirectCallMetadata(t *testing.T) {
	for _, name := range []string{"seq", "first", "second", "next", "rest", "count"} {
		if rt.LookupNativeDirect("clojure.core", name) == nil {
			t.Errorf("clojure.core/%s not registered as a native-direct fn", name)
		}
	}
}
