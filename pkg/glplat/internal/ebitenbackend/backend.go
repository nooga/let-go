//go:build glplat_ebiten

/*
 * glplat ebiten backend — Ebitengine implementation of registry.Backend
 *
 * Inverts glplat's caller-owned loop: ebiten.RunGame owns the main thread
 * (handed over by the boot trampoline in the root package), the lg program
 * runs on a goroutine, and EndFrame hands each frame's batches to Update
 * over an unbuffered channel — that rendezvous is also the frame pacing.
 *
 * Known divergence from the GL backend: no depth buffer. Batches draw in
 * submission order (painter's algorithm), which matches the 2D top-down
 * frontend but mis-sorts overlapping 3D first-person geometry.
 */

package ebitenbackend

import (
	"fmt"
	"image"
	stdcolor "image/color"
	"image/png"
	"math"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/nooga/let-go/pkg/glplat/internal/registry"
)

// Current is the singleton backend; the boot trampoline in the root package
// needs it to signal RunGame teardown (Stopped, RequestStop).
var Current = newBackend()

func init() {
	registry.SetBackend(Current)
}

type batch struct {
	tex    int
	matrix [16]float64
	verts  []float64
}

type frame struct {
	clearR, clearG, clearB float64
	batches                []batch
}

type shotRequest struct {
	batches                []batch
	clearR, clearG, clearB float64
	path                   string
	done                   chan error
}

type texture struct {
	img  *ebiten.Image
	w, h int
}

type Backend struct {
	// lg-goroutine-only state
	cur       frame
	curMatrix [16]float64
	width     int
	height    int
	startTime time.Time

	texturesMu sync.RWMutex
	textures   map[int]*texture
	nextTexID  int

	frames  chan frame
	shots   chan shotRequest
	ready   chan struct{}
	stopped chan struct{}

	stateMu        sync.Mutex
	closeRequested bool
	terminated     bool

	eventMu    sync.Mutex
	eventQueue []string
}

func newBackend() *Backend {
	return &Backend{
		textures:   make(map[int]*texture),
		frames:     make(chan frame),
		shots:      make(chan shotRequest),
		ready:      make(chan struct{}),
		stopped:    make(chan struct{}),
		eventQueue: make([]string, 0, 64),
	}
}

// StartRequest is delivered to the boot trampoline when the lg program calls
// glplat/Init: the trampoline runs ebiten.RunGame(Game) on the main goroutine
// and calls Current.Stopped() when it returns.
type StartRequest struct {
	Game          ebiten.Game
	Width, Height int
	Title         string
}

var startCh = make(chan StartRequest, 1)

// StartChan is consumed by the boot trampoline in the root package.
func StartChan() <-chan StartRequest { return startCh }

// Stopped unblocks any lg-side call waiting on the rendezvous after RunGame
// has returned; without it Terminate/EndFrame would deadlock on teardown.
func (b *Backend) Stopped() {
	select {
	case <-b.stopped:
	default:
		close(b.stopped)
	}
}

// RequestStop makes the next Update return ebiten.Termination. The trampoline
// calls it when the lg program finishes without calling glplat/Terminate.
func (b *Backend) RequestStop() {
	b.stateMu.Lock()
	b.terminated = true
	b.stateMu.Unlock()
}

// --- registry.Backend: lg-goroutine side ---

func (b *Backend) Init(width, height int, title string) error {
	b.width = width
	b.height = height
	b.startTime = time.Now()

	// Texture 0 is the untextured white fallback, same as the GL backend.
	if _, err := b.LoadTextureRGBA([]byte{255, 255, 255, 255}, 1, 1); err != nil {
		return err
	}

	startCh <- StartRequest{Game: &game{b: b}, Width: width, Height: height, Title: title}
	select {
	case <-b.ready:
		return nil
	case <-b.stopped:
		return fmt.Errorf("glplat/ebiten: window closed during Init")
	}
}

func (b *Backend) ShouldClose() bool {
	b.stateMu.Lock()
	defer b.stateMu.Unlock()
	return b.closeRequested
}

func (b *Backend) PollEventsWindow() {
	// Window events are processed inside ebiten's Update; nothing to pump.
}

func (b *Backend) BeginFrame(r, g, blue float64) {
	b.cur = frame{clearR: r, clearG: g, clearB: blue}
}

func (b *Backend) EndFrame() {
	f := b.cur
	b.cur = frame{}
	select {
	case b.frames <- f:
	case <-b.stopped:
	}
}

func (b *Backend) Time() float64 {
	return time.Since(b.startTime).Seconds()
}

func (b *Backend) Terminate() error {
	b.RequestStop()
	// Wait for RunGame to unwind so teardown is complete before returning.
	<-b.stopped
	return nil
}

func (b *Backend) LoadTextureFile(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return 0, fmt.Errorf("decode PNG: %w", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return b.LoadTextureRGBA(rgba.Pix, rgba.Bounds().Dx(), rgba.Bounds().Dy())
}

func (b *Backend) LoadTextureRGBA(pixels []byte, w, h int) (int, error) {
	if len(pixels) != w*h*4 {
		return 0, fmt.Errorf("LoadTextureRGBA: expected %d bytes, got %d", w*h*4, len(pixels))
	}
	// glplat's contract is straight alpha (what GL's TexImage2D upload takes);
	// ebiten's WritePixels expects premultiplied.
	premul := make([]byte, len(pixels))
	for i := 0; i < len(pixels); i += 4 {
		a := uint32(pixels[i+3])
		premul[i] = byte(uint32(pixels[i]) * a / 255)
		premul[i+1] = byte(uint32(pixels[i+1]) * a / 255)
		premul[i+2] = byte(uint32(pixels[i+2]) * a / 255)
		premul[i+3] = byte(a)
	}
	img := ebiten.NewImage(w, h)
	img.WritePixels(premul)

	b.texturesMu.Lock()
	id := b.nextTexID
	b.nextTexID++
	b.textures[id] = &texture{img: img, w: w, h: h}
	b.texturesMu.Unlock()
	return id, nil
}

func (b *Backend) TextureSize(id int) (w, h int) {
	b.texturesMu.RLock()
	defer b.texturesMu.RUnlock()
	if tex, ok := b.textures[id]; ok {
		return tex.w, tex.h
	}
	return 0, 0
}

func (b *Backend) SetMatrix(m []float64) error {
	if len(m) != 16 {
		return fmt.Errorf("SetMatrix: matrix must be 16 floats")
	}
	// Captured per batch at submit time, like the GL backend's uniform upload.
	copy(b.curMatrix[:], m)
	return nil
}

func (b *Backend) SubmitTriangles(tex int, verts []float64) error {
	if len(verts)%9 != 0 {
		return fmt.Errorf("SubmitTriangles: verts length must be divisible by 9")
	}
	v := make([]float64, len(verts))
	copy(v, verts)
	b.cur.batches = append(b.cur.batches, batch{tex: tex, matrix: b.curMatrix, verts: v})
	return nil
}

func (b *Backend) PollInputEvents() []string {
	b.eventMu.Lock()
	events := b.eventQueue
	b.eventQueue = make([]string, 0, 64)
	b.eventMu.Unlock()
	return events
}

func (b *Backend) WindowSize() (w, h int) {
	return b.width, b.height
}

// Screenshot renders the batches submitted so far this frame (glplat contract:
// call after submits, before EndFrame) and blocks until ebiten's next Draw has
// rasterized and encoded them.
func (b *Backend) Screenshot(path string) error {
	req := shotRequest{
		batches: append([]batch(nil), b.cur.batches...),
		clearR:  b.cur.clearR, clearG: b.cur.clearG, clearB: b.cur.clearB,
		path: path,
		done: make(chan error, 1),
	}
	select {
	case b.shots <- req:
	case <-b.stopped:
		return fmt.Errorf("Screenshot: window closed")
	}
	select {
	case err := <-req.done:
		return err
	case <-b.stopped:
		return fmt.Errorf("Screenshot: window closed")
	}
}

// --- ebiten side ---

type game struct {
	b           *Backend
	current     frame
	hasFrame    bool
	readied     bool
	pendingShot *shotRequest
}

func (g *game) Update() error {
	b := g.b

	if !g.readied {
		g.readied = true
		close(b.ready)
	}

	b.stateMu.Lock()
	terminated := b.terminated
	b.stateMu.Unlock()
	if terminated {
		return ebiten.Termination
	}
	if ebiten.IsWindowBeingClosed() {
		b.stateMu.Lock()
		b.closeRequested = true
		b.stateMu.Unlock()
	}

	g.pollInput()

	select {
	case f := <-b.frames:
		g.current = f
		g.hasFrame = true
	default:
	}

	// Screenshot requests need a render target, so stash for Draw this tick.
	if g.pendingShot == nil {
		select {
		case req := <-b.shots:
			g.pendingShot = &req
		default:
		}
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.pendingShot != nil {
		req := g.pendingShot
		g.pendingShot = nil
		req.done <- g.b.renderShot(req)
	}
	if !g.hasFrame {
		return
	}
	f := &g.current
	screen.Fill(colorFromFloats(f.clearR, f.clearG, f.clearB))
	g.b.drawBatches(screen, f.batches)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.b.width, g.b.height
}

func (b *Backend) renderShot(req *shotRequest) error {
	img := ebiten.NewImage(b.width, b.height)
	img.Fill(colorFromFloats(req.clearR, req.clearG, req.clearB))
	b.drawBatches(img, req.batches)

	pix := make([]byte, b.width*b.height*4)
	img.ReadPixels(pix)
	img.Deallocate()

	out := image.NewRGBA(image.Rect(0, 0, b.width, b.height))
	copy(out.Pix, pix)
	for i := 3; i < len(out.Pix); i += 4 {
		out.Pix[i] = 255
	}
	f, err := os.Create(req.path)
	if err != nil {
		return fmt.Errorf("Screenshot: %w", err)
	}
	defer f.Close()
	return png.Encode(f, out)
}

// drawBatches applies each batch's MVP on the CPU (clip space → NDC → screen)
// and issues DrawTriangles per batch, chunked under the uint16 index limit.
// Triangles with any vertex at w<=0 (behind the camera) are dropped rather
// than clipped — crude, but only the 3D path can hit it.
func (b *Backend) drawBatches(dst *ebiten.Image, batches []batch) {
	sw, sh := float64(b.width), float64(b.height)
	opts := &ebiten.DrawTrianglesOptions{
		ColorScaleMode: ebiten.ColorScaleModeStraightAlpha,
		Filter:         ebiten.FilterNearest,
	}

	for _, bt := range batches {
		b.texturesMu.RLock()
		tex := b.textures[0]
		if bt.tex > 0 {
			if t, ok := b.textures[bt.tex]; ok {
				tex = t
			}
		}
		b.texturesMu.RUnlock()
		if tex == nil {
			continue
		}

		m := &bt.matrix
		nVerts := len(bt.verts) / 9
		vertices := make([]ebiten.Vertex, 0, nVerts)
		indices := make([]uint16, 0, nVerts)

		flush := func() {
			if len(indices) > 0 {
				dst.DrawTriangles(vertices, indices, tex.img, opts)
			}
			vertices = vertices[:0]
			indices = indices[:0]
		}

		for t := 0; t+27 <= len(bt.verts); t += 27 {
			var sxs, sys [3]float64
			ok := true
			for k := 0; k < 3; k++ {
				o := t + k*9
				x, y, z := bt.verts[o], bt.verts[o+1], bt.verts[o+2]
				cx := m[0]*x + m[4]*y + m[8]*z + m[12]
				cy := m[1]*x + m[5]*y + m[9]*z + m[13]
				cw := m[3]*x + m[7]*y + m[11]*z + m[15]
				if cw <= 0 || math.IsNaN(cw) {
					ok = false
					break
				}
				sxs[k] = (cx/cw + 1) / 2 * sw
				sys[k] = (1 - cy/cw) / 2 * sh
			}
			if !ok {
				continue
			}
			if len(vertices) > math.MaxUint16-3 {
				flush()
			}
			base := uint16(len(vertices))
			for k := 0; k < 3; k++ {
				o := t + k*9
				u, v := bt.verts[o+3], bt.verts[o+4]
				r, gg, bb, a := bt.verts[o+5], bt.verts[o+6], bt.verts[o+7], bt.verts[o+8]
				vertices = append(vertices, ebiten.Vertex{
					DstX: float32(sxs[k]), DstY: float32(sys[k]),
					SrcX: float32(u * float64(tex.w)), SrcY: float32(v * float64(tex.h)),
					ColorR: float32(r), ColorG: float32(gg), ColorB: float32(bb), ColorA: float32(a),
				})
			}
			indices = append(indices, base, base+1, base+2)
		}
		flush()
	}
}

func colorFromFloats(r, g, b float64) stdcolor.Color {
	return stdcolor.RGBA{
		R: uint8(clamp01(r) * 255),
		G: uint8(clamp01(g) * 255),
		B: uint8(clamp01(b) * 255),
		A: 255,
	}
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// --- input ---

// repeat timing in 60Hz ticks, approximating GLFW/macOS key repeat.
const repeatDelayTicks = 30
const repeatIntervalTicks = 3

func (g *game) pollInput() {
	b := g.b
	var events []string

	for _, k := range inpututil.AppendJustPressedKeys(nil) {
		if name, ok := ebitenKeyNames[k]; ok {
			events = append(events, "key:"+name)
		}
	}
	for k, name := range ebitenKeyNames {
		d := inpututil.KeyPressDuration(k)
		if d > repeatDelayTicks && (d-repeatDelayTicks)%repeatIntervalTicks == 0 {
			events = append(events, "key-repeat:"+name)
		}
	}
	for _, r := range ebiten.AppendInputChars(nil) {
		events = append(events, "char:"+string(r))
	}

	if len(events) > 0 {
		b.eventMu.Lock()
		b.eventQueue = append(b.eventQueue, events...)
		b.eventMu.Unlock()
	}
}

var ebitenKeyNames = map[ebiten.Key]string{
	ebiten.KeySpace: "space", ebiten.KeyEscape: "escape", ebiten.KeyEnter: "enter",
	ebiten.KeyTab: "tab", ebiten.KeyBackspace: "backspace", ebiten.KeyHome: "home",
	ebiten.KeyEnd: "end", ebiten.KeyDelete: "delete", ebiten.KeyArrowLeft: "left",
	ebiten.KeyArrowRight: "right", ebiten.KeyArrowUp: "up", ebiten.KeyArrowDown: "down",
	ebiten.KeyPageUp: "pageup", ebiten.KeyPageDown: "pagedown",
	ebiten.KeyDigit0: "0", ebiten.KeyDigit1: "1", ebiten.KeyDigit2: "2", ebiten.KeyDigit3: "3",
	ebiten.KeyDigit4: "4", ebiten.KeyDigit5: "5", ebiten.KeyDigit6: "6", ebiten.KeyDigit7: "7",
	ebiten.KeyDigit8: "8", ebiten.KeyDigit9: "9",
	ebiten.KeyA: "a", ebiten.KeyB: "b", ebiten.KeyC: "c", ebiten.KeyD: "d", ebiten.KeyE: "e",
	ebiten.KeyF: "f", ebiten.KeyG: "g", ebiten.KeyH: "h", ebiten.KeyI: "i", ebiten.KeyJ: "j",
	ebiten.KeyK: "k", ebiten.KeyL: "l", ebiten.KeyM: "m", ebiten.KeyN: "n", ebiten.KeyO: "o",
	ebiten.KeyP: "p", ebiten.KeyQ: "q", ebiten.KeyR: "r", ebiten.KeyS: "s", ebiten.KeyT: "t",
	ebiten.KeyU: "u", ebiten.KeyV: "v", ebiten.KeyW: "w", ebiten.KeyX: "x", ebiten.KeyY: "y",
	ebiten.KeyZ: "z", ebiten.KeyShiftLeft: "lshift", ebiten.KeyControlLeft: "lctrl",
	ebiten.KeyAltLeft: "lalt", ebiten.KeyShiftRight: "rshift", ebiten.KeyControlRight: "rctrl",
	ebiten.KeyAltRight: "ralt", ebiten.KeyF1: "f1", ebiten.KeyF2: "f2", ebiten.KeyF3: "f3",
	ebiten.KeyF4: "f4", ebiten.KeyF5: "f5", ebiten.KeyF6: "f6", ebiten.KeyF7: "f7",
	ebiten.KeyF8: "f8", ebiten.KeyF9: "f9", ebiten.KeyF10: "f10", ebiten.KeyF11: "f11",
	ebiten.KeyF12: "f12",
}
