package api_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/api"
	"github.com/stretchr/testify/assert"
)

// TestIsIncomplete pins the REPL-loop contract: input that ends inside an
// open form is incomplete (keep reading), a hard syntax error is not
// (report now), and a successful evaluation returns no error at all.
func TestIsIncomplete(t *testing.T) {
	lg, err := api.NewLetGo("incomplete-test")
	assert.NoError(t, err)

	// Ends inside an open form: accumulate more input.
	_, err = lg.Run("(defn foo [x]")
	assert.Error(t, err)
	assert.True(t, api.IsIncomplete(err))

	// Ends inside an open string literal: same story.
	_, err = lg.Run(`(println "unterminated`)
	assert.Error(t, err)
	assert.True(t, api.IsIncomplete(err))

	// Unmatched delimiter: no amount of further input completes this.
	_, err = lg.Run("(]")
	assert.Error(t, err)
	assert.False(t, api.IsIncomplete(err))

	// Empty and whitespace-only input read as incomplete: a REPL keeps
	// prompting rather than reporting an error.
	_, err = lg.Run("")
	assert.Error(t, err)
	assert.True(t, api.IsIncomplete(err))
	_, err = lg.Run("   \n\n")
	assert.Error(t, err)
	assert.True(t, api.IsIncomplete(err))

	// Comment-only input evaluates cleanly (no error), so a REPL never
	// consults IsIncomplete for it.
	_, err = lg.Run("; just a comment")
	assert.NoError(t, err)

	// The accumulate-and-retry loop converges: the completed form runs.
	v, err := lg.Run("(defn foo [x]\n  (+ x 1))\n")
	assert.NoError(t, err)
	assert.NotNil(t, v)
}
