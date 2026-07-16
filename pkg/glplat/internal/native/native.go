//go:build glplat

/*
 * glplat native backend — GLFW + OpenGL (macOS CoreProfile)
 *
 * Requires: runtime.LockOSThread() on macOS for GL/GLFW on main thread.
 */

package native

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/nooga/let-go/pkg/glplat/internal/registry"
)

func init() {
	runtime.LockOSThread()
	registry.SetBackend(&Backend{
		textures:   make(map[int]*Texture),
		eventQueue: make([]string, 0, 64),
		glfwToName: buildKeyNameMap(),
	})
}

type Backend struct {
	window        *glfw.Window
	startTime     time.Time
	texturesMu    sync.RWMutex
	textures      map[int]*Texture
	nextTexID     int
	vao           uint32
	vbo           uint32
	shaderProgram uint32
	mvpUniformLoc int32
	matrix        [16]float64
	eventQueue    []string
	eventMu       sync.Mutex
	glfwToName    map[glfw.Key]string
	windowWidth   int
	windowHeight  int
}

type Texture struct {
	ID     uint32
	Width  int
	Height int
}

func (b *Backend) Init(width, height int, title string) error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("glfw.Init: %w", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		glfw.Terminate()
		return fmt.Errorf("glfw.CreateWindow: %w", err)
	}

	b.window = window
	b.windowWidth = width
	b.windowHeight = height
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		window.Destroy()
		glfw.Terminate()
		return fmt.Errorf("gl.Init: %w", err)
	}

	glfw.SwapInterval(1)

	if err := b.setupGL(); err != nil {
		window.Destroy()
		glfw.Terminate()
		return err
	}

	window.SetKeyCallback(b.keyCallback)
	window.SetCharCallback(b.charCallback)

	b.startTime = time.Now()
	return nil
}

func (b *Backend) setupGL() error {
	vertSrc := `#version 410 core
layout(location = 0) in vec3 pos;
layout(location = 1) in vec2 uv;
layout(location = 2) in vec4 color;

out vec2 fragUV;
out vec4 fragColor;

uniform mat4 MVP;

void main() {
	gl_Position = MVP * vec4(pos, 1.0);
	fragUV = uv;
	fragColor = color;
}
`
	fragSrc := `#version 410 core
in vec2 fragUV;
in vec4 fragColor;

uniform sampler2D tex;

out vec4 outColor;

void main() {
	vec4 texel = texture(tex, fragUV);
	outColor = texel * fragColor;
}
`

	var err error
	b.shaderProgram, err = compileShaderProgram(vertSrc, fragSrc)
	if err != nil {
		return err
	}

	b.mvpUniformLoc = gl.GetUniformLocation(b.shaderProgram, gl.Str("MVP\x00"))

	gl.GenVertexArrays(1, &b.vao)
	gl.GenBuffers(1, &b.vbo)

	gl.BindVertexArray(b.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 9*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 9*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, 9*4, gl.PtrOffset(5*4))
	gl.EnableVertexAttribArray(2)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	whitePixel := []byte{255, 255, 255, 255}
	_, err = b.LoadTextureRGBA(whitePixel, 1, 1)
	if err != nil {
		return err
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)

	return nil
}

func (b *Backend) ShouldClose() bool {
	if b.window == nil {
		return true
	}
	return b.window.ShouldClose()
}

func (b *Backend) PollEventsWindow() {
	glfw.PollEvents()
}

func (b *Backend) BeginFrame(r, g, blue float64) {
	// Track the framebuffer size every frame: on Retina displays the
	// framebuffer is larger than the logical window, and it changes when
	// the window moves between displays or resizes.
	if b.window != nil {
		fbW, fbH := b.window.GetFramebufferSize()
		gl.Viewport(0, 0, int32(fbW), int32(fbH))
	}
	gl.ClearColor(float32(r), float32(g), float32(blue), 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (b *Backend) EndFrame() {
	if b.window != nil {
		b.window.SwapBuffers()
	}
	if err := gl.GetError(); err != 0 {
		fmt.Printf("GL error: 0x%x\n", err)
	}
}

func (b *Backend) Time() float64 {
	return time.Since(b.startTime).Seconds()
}

func (b *Backend) Terminate() error {
	if b.vao != 0 {
		gl.DeleteVertexArrays(1, &b.vao)
	}
	if b.vbo != 0 {
		gl.DeleteBuffers(1, &b.vbo)
	}
	if b.shaderProgram != 0 {
		gl.DeleteProgram(b.shaderProgram)
	}

	b.texturesMu.Lock()
	for _, tex := range b.textures {
		gl.DeleteTextures(1, &tex.ID)
	}
	b.texturesMu.Unlock()

	if b.window != nil {
		b.window.Destroy()
	}
	glfw.Terminate()
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

	w := rgba.Bounds().Dx()
	h := rgba.Bounds().Dy()
	return b.LoadTextureRGBA(rgba.Pix, w, h)
}

func (b *Backend) LoadTextureRGBA(pixels []byte, w, h int) (int, error) {
	var texID uint32
	gl.GenTextures(1, &texID)
	gl.BindTexture(gl.TEXTURE_2D, texID)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	b.texturesMu.Lock()
	id := b.nextTexID
	b.nextTexID++
	b.textures[id] = &Texture{ID: texID, Width: w, Height: h}
	b.texturesMu.Unlock()

	return id, nil
}

func (b *Backend) TextureSize(id int) (w, h int) {
	b.texturesMu.RLock()
	defer b.texturesMu.RUnlock()
	if tex, ok := b.textures[id]; ok {
		return tex.Width, tex.Height
	}
	return 0, 0
}

func (b *Backend) SetMatrix(m []float64) error {
	if len(m) != 16 {
		return fmt.Errorf("SetMatrix: matrix must be 16 floats")
	}
	copy(b.matrix[:], m)
	return nil
}

func (b *Backend) SubmitTriangles(texID int, verts []float64) error {
	if len(verts)%9 != 0 {
		return fmt.Errorf("SubmitTriangles: verts length must be divisible by 9")
	}

	gl.UseProgram(b.shaderProgram)
	m32 := make([]float32, 16)
	for i := range m32 {
		m32[i] = float32(b.matrix[i])
	}
	// Matrices arrive column-major (the API contract); transpose must be
	// false or the GPU applies the transpose of every MVP.
	gl.UniformMatrix4fv(b.mvpUniformLoc, 1, false, &m32[0])

	b.texturesMu.RLock()
	glTexID := uint32(0)
	if texID > 0 {
		if tex, ok := b.textures[texID]; ok {
			glTexID = tex.ID
		}
	} else {
		if tex, ok := b.textures[0]; ok {
			glTexID = tex.ID
		}
	}
	b.texturesMu.RUnlock()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, glTexID)

	gl.BindVertexArray(b.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)

	vertCount := len(verts) / 9
	verts32 := make([]float32, len(verts))
	for i, v := range verts {
		verts32[i] = float32(v)
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(verts32)*4, gl.Ptr(verts32), gl.STREAM_DRAW)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(vertCount))

	return nil
}

// Screenshot reads the current framebuffer (back buffer — call after
// rendering, BEFORE EndFrame swaps) and writes it as PNG.
func (b *Backend) Screenshot(path string) error {
	if b.window == nil {
		return fmt.Errorf("Screenshot: no window")
	}
	fbW, fbH := b.window.GetFramebufferSize()
	pixels := make([]byte, fbW*fbH*4)
	gl.ReadPixels(0, 0, int32(fbW), int32(fbH), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	img := image.NewRGBA(image.Rect(0, 0, fbW, fbH))
	// GL rows are bottom-up; flip into the image's top-down layout.
	rowLen := fbW * 4
	for y := 0; y < fbH; y++ {
		src := pixels[(fbH-1-y)*rowLen : (fbH-y)*rowLen]
		dst := img.Pix[y*img.Stride : y*img.Stride+rowLen]
		copy(dst, src)
	}
	// Force opaque alpha (the framebuffer's alpha channel is meaningless here).
	for i := 3; i < len(img.Pix); i += 4 {
		img.Pix[i] = 255
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Screenshot: %w", err)
	}
	defer f.Close()
	return png.Encode(f, img)
}

func (b *Backend) PollInputEvents() []string {
	b.eventMu.Lock()
	events := b.eventQueue
	b.eventQueue = make([]string, 0, 64)
	b.eventMu.Unlock()
	return events
}

func (b *Backend) WindowSize() (w, h int) {
	return b.windowWidth, b.windowHeight
}

func (b *Backend) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if name, ok := b.glfwToName[key]; ok {
		var event string
		if action == glfw.Press {
			event = "key:" + name
		} else if action == glfw.Repeat {
			event = "key-repeat:" + name
		}
		if event != "" {
			b.eventMu.Lock()
			b.eventQueue = append(b.eventQueue, event)
			b.eventMu.Unlock()
		}
	}
}

func (b *Backend) charCallback(w *glfw.Window, codepoint rune) {
	event := "char:" + string(codepoint)
	b.eventMu.Lock()
	b.eventQueue = append(b.eventQueue, event)
	b.eventMu.Unlock()
}

func compileShaderProgram(vertSrc, fragSrc string) (uint32, error) {
	vertShader, err := compileShader(vertSrc, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertShader)

	fragShader, err := compileShader(fragSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLen)
		log := make([]byte, logLen)
		gl.GetProgramInfoLog(program, logLen, nil, &log[0])
		return 0, fmt.Errorf("link shader program: %s", string(log))
	}

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	defer free()
	gl.ShaderSource(shader, 1, csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)
		log := make([]byte, logLen)
		gl.GetShaderInfoLog(shader, logLen, nil, &log[0])
		gl.DeleteShader(shader)
		return 0, fmt.Errorf("compile shader: %s", string(log))
	}

	return shader, nil
}

func buildKeyNameMap() map[glfw.Key]string {
	return map[glfw.Key]string{
		glfw.KeySpace: "space", glfw.KeyEscape: "escape", glfw.KeyEnter: "enter",
		glfw.KeyTab: "tab", glfw.KeyBackspace: "backspace", glfw.KeyHome: "home",
		glfw.KeyEnd: "end", glfw.KeyDelete: "delete", glfw.KeyLeft: "left",
		glfw.KeyRight: "right", glfw.KeyUp: "up", glfw.KeyDown: "down",
		glfw.KeyPageUp: "pageup", glfw.KeyPageDown: "pagedown",
		glfw.Key0: "0", glfw.Key1: "1", glfw.Key2: "2", glfw.Key3: "3", glfw.Key4: "4",
		glfw.Key5: "5", glfw.Key6: "6", glfw.Key7: "7", glfw.Key8: "8", glfw.Key9: "9",
		glfw.KeyA: "a", glfw.KeyB: "b", glfw.KeyC: "c", glfw.KeyD: "d", glfw.KeyE: "e",
		glfw.KeyF: "f", glfw.KeyG: "g", glfw.KeyH: "h", glfw.KeyI: "i", glfw.KeyJ: "j",
		glfw.KeyK: "k", glfw.KeyL: "l", glfw.KeyM: "m", glfw.KeyN: "n", glfw.KeyO: "o",
		glfw.KeyP: "p", glfw.KeyQ: "q", glfw.KeyR: "r", glfw.KeyS: "s", glfw.KeyT: "t",
		glfw.KeyU: "u", glfw.KeyV: "v", glfw.KeyW: "w", glfw.KeyX: "x", glfw.KeyY: "y",
		glfw.KeyZ: "z", glfw.KeyLeftShift: "lshift", glfw.KeyLeftControl: "lctrl",
		glfw.KeyLeftAlt: "lalt", glfw.KeyRightShift: "rshift", glfw.KeyRightControl: "rctrl",
		glfw.KeyRightAlt: "ralt", glfw.KeyF1: "f1", glfw.KeyF2: "f2", glfw.KeyF3: "f3",
		glfw.KeyF4: "f4", glfw.KeyF5: "f5", glfw.KeyF6: "f6", glfw.KeyF7: "f7",
		glfw.KeyF8: "f8", glfw.KeyF9: "f9", glfw.KeyF10: "f10", glfw.KeyF11: "f11",
		glfw.KeyF12: "f12",
	}
}
