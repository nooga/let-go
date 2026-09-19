// The native GLFW/OpenGL backend is opt-in: build with -tags glplat.
// Without the tag no backend registers and the pure-Go API returns
// "no backend registered" errors, so plain `go build ./...`, headless CI,
// and the GOOS=js wasm target never need cgo/GLFW.
//go:build glplat

package glplat

import (
	// Blank import to register the native GLFW/OpenGL backend.
	// This ensures runtime.LockOSThread() is called during init.
	_ "github.com/nooga/let-go/pkg/glplat/internal/native"
)
