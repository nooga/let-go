//go:build !glplat_ebiten

package main

// bootMain runs the interpreter directly; the glplat_ebiten build replaces
// this with a trampoline that hands the main thread to ebiten.
func bootMain(run func() int) int {
	return run()
}
