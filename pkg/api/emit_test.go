package api_test

import (
	"testing"

	"github.com/nooga/let-go/pkg/api"
)

// TestWithEmit proves the emit dual of TestWithStdout: a Go callback passed
// via api.WithEmit receives the (name, dataJSON) pair a (js/emit ...) call
// dispatches, with the data JSON-marshaled the same way the WASM bundle
// hands it to LetGoHost.onEmit.
func TestWithEmit(t *testing.T) {
	var gotName, gotJSON string
	calls := 0
	lg, err := api.NewLetGo("withemit-test", api.WithEmit(func(name, dataJSON string) {
		calls++
		gotName, gotJSON = name, dataJSON
	}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lg.Run(`(js/emit :stats {:hp 7})`); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 emit, got %d", calls)
	}
	if gotName != "stats" {
		t.Errorf("event name = %q, want %q", gotName, "stats")
	}
	if gotJSON != `{"hp":7}` {
		t.Errorf("event data = %q, want %q", gotJSON, `{"hp":7}`)
	}
}

// TestEmitNoopWithoutOption proves (js/emit ...) is harmless — validates args
// and returns nil — when no emitter is installed (the default nopEmitter
// root). A bad event-name type still errors, so native dev catches type bugs.
func TestEmitNoopWithoutOption(t *testing.T) {
	lg, err := api.NewLetGo("emit-noop-test")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := lg.Run(`(js/emit :ping {:x 1})`); err != nil {
		t.Fatalf("emit with no host should no-op, got %v", err)
	}
	if _, err := lg.Run(`(js/emit 42 {:x 1})`); err == nil {
		t.Error("expected arg-validation error for non-name event, got nil")
	}
}
