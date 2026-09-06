//go:build glplat_ebiten

package main

import (
	"github.com/nooga/let-go/pkg/glplat"
)

// bootMain hands the main thread to ebiten; see glplat.EbitenBootMain.
func bootMain(run func() int) int {
	return glplat.EbitenBootMain(run)
}
