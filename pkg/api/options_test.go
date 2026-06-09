package api_test

import (
	"bytes"
	"testing"

	"github.com/nooga/let-go/pkg/api"
	"github.com/stretchr/testify/assert"
)

// TestWithStdout proves the basic embedder-capture story: a buffer
// passed via api.WithStdout receives (println ...) output as if it
// were the runtime's stdout.
func TestWithStdout(t *testing.T) {
	var buf bytes.Buffer
	lg, err := api.NewLetGo("withstdout-test", api.WithStdout(&buf))
	assert.NoError(t, err)

	_, err = lg.Run(`(println "hello-from-api")`)
	assert.NoError(t, err)
	assert.Equal(t, "hello-from-api\n", buf.String())
}

// TestWithStderr proves the same story for *err* — including the
// previously-broken (binding [*out* *err*] (println ...)) pattern.
func TestWithStderr(t *testing.T) {
	var out, err bytes.Buffer
	lg, errInit := api.NewLetGo("withstderr-test",
		api.WithStdout(&out),
		api.WithStderr(&err))
	assert.NoError(t, errInit)

	// Direct write to *err* via the IOHandle.
	_, errRun := lg.Run(`(write! *err* "direct-to-err")`)
	assert.NoError(t, errRun)
	assert.Equal(t, "", out.String())
	assert.Equal(t, "direct-to-err", err.String())

	// Reset, then test the binding-redirected println case.
	out.Reset()
	err.Reset()
	_, errRun = lg.Run(`(binding [*out* *err*] (println "redirected"))`)
	assert.NoError(t, errRun)
	assert.Equal(t, "", out.String())
	assert.Equal(t, "redirected\n", err.String())
}

// TestPrintFnsConsultOut covers all four print fns (print/pr/prn/println)
// and confirms the runtime correctly routes each through *out*.
func TestPrintFnsConsultOut(t *testing.T) {
	var buf bytes.Buffer
	lg, err := api.NewLetGo("printfns-test", api.WithStdout(&buf))
	assert.NoError(t, err)

	// Run only processes one form per call, so wrap in a do.
	_, err = lg.Run(`(do (print "a") (pr "b") (println "c") (prn "d"))`)
	assert.NoError(t, err)
	// print: a (no newline)
	// pr:    "b" (readable, no newline)
	// println: c (with newline)
	// prn: "d" (readable, with newline)
	assert.Equal(t, `a"b"c
"d"
`, buf.String())
}

// TestWithOutStrUnaffected proves the with-out-str rewrite preserves
// its semantic: nested captures don't leak past the inner binding.
func TestWithOutStrUnaffected(t *testing.T) {
	var buf bytes.Buffer
	lg, err := api.NewLetGo("withoutstr-test", api.WithStdout(&buf))
	assert.NoError(t, err)

	val, err := lg.Run(`(with-out-str (println "captured"))`)
	assert.NoError(t, err)
	// The with-out-str result IS the captured value.
	assert.Equal(t, `"captured\n"`, val.String())
	// And nothing leaked to the outer *out*.
	assert.Equal(t, "", buf.String())
}

// TestTwoRuntimesIsolation is the regression test for the bug codex
// caught: WithStdout used to mutate the process-global root binding of
// core/*out*, so creating a second LetGo would silently redirect the
// first one's output. The fix moves the binding push to Run scope, so
// each Run only sees its own configured stream.
func TestTwoRuntimesIsolation(t *testing.T) {
	var bufA, bufB bytes.Buffer
	lgA, err := api.NewLetGo("iso-a", api.WithStdout(&bufA))
	assert.NoError(t, err)
	lgB, err := api.NewLetGo("iso-b", api.WithStdout(&bufB))
	assert.NoError(t, err)

	_, err = lgA.Run(`(println "from-A")`)
	assert.NoError(t, err)
	_, err = lgB.Run(`(println "from-B")`)
	assert.NoError(t, err)

	assert.Equal(t, "from-A\n", bufA.String())
	assert.Equal(t, "from-B\n", bufB.String())

	// Interleaved: run A again. Should still hit bufA.
	_, err = lgA.Run(`(println "from-A-again")`)
	assert.NoError(t, err)
	assert.Equal(t, "from-A\nfrom-A-again\n", bufA.String())
	assert.Equal(t, "from-B\n", bufB.String())
}

// TestNoOptionRuntimeUnaffectedByConfiguredOne — a LetGo constructed
// without WithStdout must NOT inherit the writer of a previously-
// constructed configured LetGo. Pre-fix, this would have failed
// because WithStdout mutated the global root.
//
// We can't easily assert against the process os.Stdout from a test,
// but we CAN assert that the configured runtime's buffer doesn't
// receive output from the unconfigured one's Run.
func TestNoOptionRuntimeUnaffectedByConfiguredOne(t *testing.T) {
	var bufConfigured bytes.Buffer
	lgConfigured, err := api.NewLetGo("dflt-cfg", api.WithStdout(&bufConfigured))
	assert.NoError(t, err)
	lgDefault, err := api.NewLetGo("dflt-no-cfg")
	assert.NoError(t, err)

	// Configured runtime writes into its buffer.
	_, err = lgConfigured.Run(`(println "configured")`)
	assert.NoError(t, err)
	assert.Equal(t, "configured\n", bufConfigured.String())

	// Default runtime's output goes to os.Stdout (not asserted) — the
	// load-bearing claim is that it does NOT land in bufConfigured.
	_, err = lgDefault.Run(`(println "default")`)
	assert.NoError(t, err)
	assert.Equal(t, "configured\n", bufConfigured.String())
}
