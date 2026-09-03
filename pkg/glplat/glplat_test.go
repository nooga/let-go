package glplat

import "testing"

// LoadTextureRGBA must reject a pixel buffer shorter than w*h*4 before the
// data reaches a backend, since the native backend hands the slice to a C
// function that reads w*h*4 bytes. The guard runs with no backend registered,
// so it is exercised here without GLFW.
func TestLoadTextureRGBARejectsShortBuffer(t *testing.T) {
	// A 2x2 RGBA image needs 16 bytes; give it 8.
	if _, err := LoadTextureRGBA(make([]byte, 8), 2, 2); err == nil {
		t.Fatal("short pixel buffer: want an error, got nil")
	}
	if _, err := LoadTextureRGBA(make([]byte, 64), -1, 2); err == nil {
		t.Fatal("negative dimensions: want an error, got nil")
	}

	// A correctly sized buffer passes the guard and reaches the (absent)
	// backend, so it fails with the no-backend error, not the arg error.
	_, err := LoadTextureRGBA(make([]byte, 16), 2, 2)
	if err == nil {
		t.Fatal("valid buffer with no backend: want an error, got nil")
	}
	if _, isArg := err.(*invalidArgError); isArg {
		t.Fatalf("a correctly sized buffer should pass the guard; got %v", err)
	}
}
