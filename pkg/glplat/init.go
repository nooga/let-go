package glplat

import (
	// Blank import to register the native GLFW/OpenGL backend.
	// This ensures runtime.LockOSThread() is called during init.
	_ "github.com/nooga/let-go/pkg/glplat/internal/native"
)
