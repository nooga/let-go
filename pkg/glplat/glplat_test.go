package glplat

import (
	"testing"

	"github.com/nooga/let-go/pkg/glplat/internal/registry"
)

// LoadTextureRGBA must reject a pixel buffer shorter than w*h*4 before the
// data reaches a backend, since the native backend hands the slice to a C
// function that reads w*h*4 bytes. The guard itself needs no backend, so it is
// exercised here without GLFW.
func TestLoadTextureRGBARejectsShortBuffer(t *testing.T) {
	// A 2x2 RGBA image needs 16 bytes; give it 8.
	if _, err := LoadTextureRGBA(make([]byte, 8), 2, 2); err == nil {
		t.Fatal("short pixel buffer: want an error, got nil")
	}
	if _, err := LoadTextureRGBA(make([]byte, 64), -1, 2); err == nil {
		t.Fatal("negative dimensions: want an error, got nil")
	}

	// Everything below calls with a correctly sized buffer, which by design
	// passes the guard and reaches the backend. A tagged build has one
	// registered, and calling it here would run backend code with no live
	// context: the native backend goes straight to glGenTextures, which
	// segfaults under -tags glplat. The guard is what this test is about, and
	// it is covered above, so stop here.
	if registry.Get() != nil {
		return
	}

	// Untagged, no backend registers, so a correctly sized buffer surfaces the
	// no-backend error rather than the arg error.
	_, err := LoadTextureRGBA(make([]byte, 16), 2, 2)
	if err == nil {
		t.Fatal("valid buffer with no backend: want an error, got nil")
	}
	if _, isArg := err.(*invalidArgError); isArg {
		t.Fatalf("a correctly sized buffer should pass the guard; got %v", err)
	}
}
