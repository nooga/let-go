package api_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/api"
	"github.com/stretchr/testify/assert"
)

// TestNSDef proves the organized-registration story: host functions
// defined into a dedicated namespace via lg.NS are callable through
// their qualified names, with the same boxing LetGo.Def performs.
func TestNSDef(t *testing.T) {
	lg, err := api.NewLetGo("nsdef-test")
	assert.NoError(t, err)

	music := lg.NS("nsdef-test.music")
	assert.NoError(t, music.Def("add2", func(a, b int) int { return a + b }))
	assert.NoError(t, music.Def("answer", 42))

	v, err := lg.Run("(nsdef-test.music/add2 20 2)")
	assert.NoError(t, err)
	assert.Equal(t, "22", v.String())

	v, err = lg.Run("nsdef-test.music/answer")
	assert.NoError(t, err)
	assert.Equal(t, "42", v.String())
}

// TestNSDefShadowing proves the intentional-collision story: a host DSL
// name that collides with clojure.core (here `filter`) is defined via
// DefShadowing, and unqualified references in that namespace resolve to
// the host definition rather than the core var.
func TestNSDefShadowing(t *testing.T) {
	lg, err := api.NewLetGo("nsshadow-test")
	assert.NoError(t, err)

	ns := lg.NS("nsshadow-test")
	assert.NoError(t, ns.DefShadowing("filter", func(freq int) string {
		return "lowpass"
	}))

	// Unqualified `filter` in the instance namespace is the host's, not
	// clojure.core's lazy-seq filter.
	v, err := lg.Run("(filter 220)")
	assert.NoError(t, err)
	assert.Equal(t, `"lowpass"`, v.String())

	// Core's filter is still reachable qualified.
	v, err = lg.Run("(clojure.core/filter odd? [1 2 3])")
	assert.NoError(t, err)
	assert.Equal(t, "(1 3)", v.String())
}

// TestNSCreatesNamespace pins resolve-or-create: NS on a name that was
// never mentioned before returns a usable handle rather than nil.
func TestNSCreatesNamespace(t *testing.T) {
	lg, err := api.NewLetGo("nscreate-test")
	assert.NoError(t, err)

	fresh := lg.NS("nscreate-test.brand.new")
	assert.NotNil(t, fresh)
	assert.NoError(t, fresh.Def("marker", "here"))

	v, err := lg.Run("nscreate-test.brand.new/marker")
	assert.NoError(t, err)
	assert.Equal(t, `"here"`, v.String())
}
