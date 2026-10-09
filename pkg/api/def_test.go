package api_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/api"
	"github.com/nooga/let-go/pkg/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefRedefinitionReachesRetainedFn proves a host redefinition through
// LetGo.Def updates the Var a previously compiled function references, the
// same as a source-level (def x 2) does.
func TestDefRedefinitionReachesRetainedFn(t *testing.T) {
	for _, tc := range []struct {
		name  string
		redef func(*api.LetGo) error
	}{
		{"api", func(lg *api.LetGo) error { return lg.Def("x", 2) }},
		{"source", func(lg *api.LetGo) error { _, err := lg.Run("(def x 2)"); return err }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lg, err := api.NewLetGo("def-redef-" + tc.name)
			require.NoError(t, err)
			require.NoError(t, lg.Def("x", 1))
			retained, err := lg.Run("(fn [] x)")
			require.NoError(t, err)

			require.NoError(t, tc.redef(lg))

			current, err := lg.Run("x")
			require.NoError(t, err)
			assert.Equal(t, "2", current.String())
			got, err := retained.(vm.Fn).Invoke(nil)
			require.NoError(t, err)
			assert.Equal(t, "2", got.String())
		})
	}
}

// TestNSDefRedefinitionReachesRetainedFn is the NS.Def counterpart: a
// caller compiled against a qualified host var sees the redefined value.
func TestNSDefRedefinitionReachesRetainedFn(t *testing.T) {
	lg, err := api.NewLetGo("nsdef-redef-test")
	require.NoError(t, err)
	ns := lg.NS("nsdef-redef-test.host")
	require.NoError(t, ns.Def("f", func() string { return "old" }))
	retained, err := lg.Run("(fn [] (nsdef-redef-test.host/f))")
	require.NoError(t, err)

	require.NoError(t, ns.Def("f", func() string { return "new" }))

	got, err := retained.(vm.Fn).Invoke(nil)
	require.NoError(t, err)
	assert.Equal(t, `"new"`, got.String())
}

// TestNSDefShadowingRedefinitionReachesRetainedFn covers the shadowing path:
// a caller compiled against the host's shadowing var sees the redefinition,
// and clojure.core's var is left alone.
func TestNSDefShadowingRedefinitionReachesRetainedFn(t *testing.T) {
	lg, err := api.NewLetGo("nsshadow-redef-test")
	require.NoError(t, err)
	ns := lg.NS("nsshadow-redef-test")
	require.NoError(t, ns.DefShadowing("filter", func() string { return "old" }))
	retained, err := lg.Run("(fn [] (filter))")
	require.NoError(t, err)

	require.NoError(t, ns.DefShadowing("filter", func() string { return "new" }))

	got, err := retained.(vm.Fn).Invoke(nil)
	require.NoError(t, err)
	assert.Equal(t, `"new"`, got.String())
	v, err := lg.Run("(clojure.core/filter odd? [1 2 3])")
	require.NoError(t, err)
	assert.Equal(t, "(1 3)", v.String())
}
