/*
 * glplat — minimal OpenGL platform for let-go spike
 *
 * Pure Go public API; native GLFW/OpenGL backend in internal/native.
 */

package glplat

import (
	"github.com/nooga/let-go/pkg/glplat/internal/registry"
)

// --- Public API ---

// Init initializes the platform with a window of the given dimensions and title.
// Must be called before any other platform functions.
func Init(width, height int, title string) error {
	backend := registry.Get()
	if backend == nil {
		return errNoBackend("Init")
	}
	return backend.Init(width, height, title)
}

// ShouldClose returns true if the window should close.
func ShouldClose() bool {
	backend := registry.Get()
	if backend == nil {
		return false
	}
	return backend.ShouldClose()
}

// PollEvents polls for window events (resize, close, etc.).
func PollEvents() {
	backend := registry.Get()
	if backend == nil {
		return
	}
	backend.PollEventsWindow()
}

// BeginFrame clears the framebuffer with the given color (0.0-1.0 range).
func BeginFrame(r, g, b float64) {
	backend := registry.Get()
	if backend == nil {
		return
	}
	backend.BeginFrame(r, g, b)
}

// EndFrame presents the framebuffer to the screen.
func EndFrame() {
	backend := registry.Get()
	if backend == nil {
		return
	}
	backend.EndFrame()
}

// Time returns elapsed time in seconds since Init() was called.
func Time() float64 {
	backend := registry.Get()
	if backend == nil {
		return 0
	}
	return backend.Time()
}

// Terminate closes the window and cleans up resources.
func Terminate() error {
	backend := registry.Get()
	if backend == nil {
		return nil
	}
	return backend.Terminate()
}

// LoadTextureFile loads a PNG texture from the given file path.
// Returns a texture ID (handle) for use with SubmitTriangles.
func LoadTextureFile(path string) (int, error) {
	backend := registry.Get()
	if backend == nil {
		return 0, errNoBackend("LoadTextureFile")
	}
	return backend.LoadTextureFile(path)
}

// LoadTextureRGBA loads a texture from raw RGBA pixel data.
// pixels: raw RGBA bytes in row-major order (width * height * 4 bytes)
func LoadTextureRGBA(pixels []byte, w, h int) (int, error) {
	backend := registry.Get()
	if backend == nil {
		return 0, errNoBackend("LoadTextureRGBA")
	}
	return backend.LoadTextureRGBA(pixels, w, h)
}

// TextureSize returns the width and height of a texture.
func TextureSize(id int) (w, h int) {
	backend := registry.Get()
	if backend == nil {
		return 0, 0
	}
	return backend.TextureSize(id)
}

// SetMatrix sets the MVP (model-view-projection) matrix for subsequent rendering.
// m must be 16 floats in column-major order (OpenGL convention).
func SetMatrix(m []float64) error {
	backend := registry.Get()
	if backend == nil {
		return errNoBackend("SetMatrix")
	}
	if len(m) != 16 {
		return &invalidArgError{fn: "SetMatrix", arg: "m", reason: "must be 16 floats"}
	}
	return backend.SetMatrix(m)
}

// SubmitTriangles submits a batch of vertices for rendering.
// tex: texture ID (0 = untextured, white 1x1 fallback)
// verts: flat array of floats, 9 per vertex (x, y, z, u, v, r, g, b, a)
func SubmitTriangles(tex int, verts []float64) error {
	backend := registry.Get()
	if backend == nil {
		return errNoBackend("SubmitTriangles")
	}
	if len(verts)%9 != 0 {
		return &invalidArgError{fn: "SubmitTriangles", arg: "verts", reason: "length must be divisible by 9"}
	}
	return backend.SubmitTriangles(tex, verts)
}

// PollInputEvents drains and returns a slice of recorded input events.
// Format: "key:<name>" for key press, "key-repeat:<name>" for repeat,
// "char:<c>" for character input, where <name> is lowercase GLFW key name.
func PollInputEvents() []string {
	backend := registry.Get()
	if backend == nil {
		return nil
	}
	return backend.PollInputEvents()
}

// WindowSize returns the current window width and height in pixels.
func WindowSize() (w, h int) {
	backend := registry.Get()
	if backend == nil {
		return 0, 0
	}
	return backend.WindowSize()
}

// errNoBackend returns an error when no backend is registered.
func errNoBackend(op string) error {
	return &backendError{op: op}
}

type backendError struct {
	op string
}

func (e *backendError) Error() string {
	return "glplat: no backend registered for " + e.op
}

type invalidArgError struct {
	fn     string
	arg    string
	reason string
}

func (e *invalidArgError) Error() string {
	return "glplat: " + e.fn + ": invalid " + e.arg + ": " + e.reason
}
