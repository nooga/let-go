/*
 * registry — shared backend registration for glplat
 */

package registry

// Backend is the interface for platform-specific implementation.
type Backend interface {
	Init(width, height int, title string) error
	ShouldClose() bool
	PollEventsWindow() // Process window events
	BeginFrame(r, g, b float64)
	EndFrame()
	Time() float64
	Terminate() error
	LoadTextureFile(path string) (int, error)
	LoadTextureRGBA(pixels []byte, w, h int) (int, error)
	TextureSize(id int) (w, h int)
	SetMatrix(m []float64) error
	SubmitTriangles(tex int, verts []float64) error
	PollInputEvents() []string
	WindowSize() (w, h int)
}

// backend holds the current platform backend instance.
var backend Backend

// SetBackend registers the platform backend.
func SetBackend(b Backend) {
	backend = b
}

// Get returns the registered backend, or nil if none.
func Get() Backend {
	return backend
}
