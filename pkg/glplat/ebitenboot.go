//go:build glplat_ebiten

package glplat

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nooga/let-go/pkg/glplat/internal/ebitenbackend"
)

// EbitenBootMain hands the main thread to ebiten: the lg program runs on a
// goroutine, and if it calls glplat/Init the backend delivers a start request
// here, where RunGame takes over until the program terminates the platform
// (or finishes, which force-stops the window). The root package's main must
// route through this in glplat_ebiten builds — ebiten cannot run off the
// main thread on macOS.
func EbitenBootMain(run func() int) int {
	done := make(chan int, 1)
	go func() {
		code := run()
		ebitenbackend.Current.RequestStop()
		done <- code
	}()

	select {
	case code := <-done:
		return code
	case req := <-ebitenbackend.StartChan():
		ebiten.SetWindowSize(req.Width, req.Height)
		ebiten.SetWindowTitle(req.Title)
		ebiten.SetWindowClosingHandled(true)
		err := ebiten.RunGame(req.Game)
		ebitenbackend.Current.Stopped()
		if err != nil {
			fmt.Fprintf(os.Stderr, "glplat/ebiten: %v\n", err)
		}
		return <-done
	}
}
