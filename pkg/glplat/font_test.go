package glplat

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

// loadTestFont writes the embedded Go font to a temp file and loads it. The
// public FontLoad takes a path, so the fixture lands on disk under t.TempDir.
func loadTestFont(t *testing.T) int {
	t.Helper()
	path := filepath.Join(t.TempDir(), "go.ttf")
	if err := os.WriteFile(path, goregular.TTF, 0o644); err != nil {
		t.Fatalf("write test font: %v", err)
	}
	id, err := FontLoad(path)
	if err != nil {
		t.Fatalf("FontLoad: %v", err)
	}
	return id
}

// The glyph-lookup paths shared a per-font sfnt.Buffer that GlyphIndex
// mutates, so concurrent callers raced on it. Run them in parallel; the race
// detector is what makes this a regression test.
func TestFontConcurrentGlyphAccess(t *testing.T) {
	id := loadTestFont(t)
	const workers = 8
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, ch := range []string{"A", "g", "W", "."} {
				FontHasGlyph(id, ch)
				if _, err := FontRasterizeCell(id, ch, 12, 16); err != nil {
					t.Errorf("FontRasterizeCell(%q): %v", ch, err)
				}
			}
		}()
	}
	wg.Wait()
}

// getFace discarded NewFace's error and cached the nil face, so a later
// DrawString panicked. Non-positive sizes are the trigger the width-fit
// rescale can reach when int(newSize) rounds to zero.
func TestGetFaceRejectsNonPositiveSize(t *testing.T) {
	id := loadTestFont(t)
	entry := fontRegistry[id]
	entry.mu.Lock()
	defer entry.mu.Unlock()

	for _, size := range []int{0, -1, -100} {
		f, err := getFaceLocked(entry, size)
		if err == nil || f != nil {
			t.Errorf("getFaceLocked(size=%d) = (%v, %v); want (nil, error)", size, f, err)
		}
	}
	if _, ok := entry.faces[0]; ok {
		t.Error("getFaceLocked cached a face for size 0")
	}

	if f, err := getFaceLocked(entry, 12); err != nil || f == nil {
		t.Errorf("getFaceLocked(size=12) = (%v, %v); want a face and nil error", f, err)
	}
}

// A 1x1 cell drives the width-fit rescale toward a zero-size face. It must
// return a cell rather than panic, falling back to the unscaled face.
func TestFontRasterizeNarrowCellNoPanic(t *testing.T) {
	id := loadTestFont(t)
	got, err := FontRasterizeCell(id, "W", 1, 1)
	if err != nil {
		t.Fatalf("FontRasterizeCell(narrow): %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("cell length = %d; want 1", len(got))
	}
}
