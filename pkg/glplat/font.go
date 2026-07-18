/*
 * font — font loading and rasterization for atlas baking
 *
 * Provides pure Go functions for font operations; no GL/window dependencies.
 */

package glplat

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// fontRegistry holds loaded font data, keyed by handle ID.
var (
	fontRegistry  = make(map[int]*fontEntry)
	fontIDCounter = 1
	fontMutex     sync.Mutex
)

type fontEntry struct {
	path  string
	fnt   *sfnt.Font
	buf   sfnt.Buffer
	faces map[int]font.Face // keyed by size
}

// FontLoad parses a TTF/OTF/TTC font and returns a handle ID.
// The font is stored in a registry for later use with FontHasGlyph and FontRasterizeCell.
func FontLoad(path string) (int, error) {
	fontData, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read font file: %w", err)
	}

	var fnt *sfnt.Font
	if strings.HasSuffix(strings.ToLower(path), ".ttc") {
		coll, err := opentype.ParseCollection(fontData)
		if err != nil {
			return 0, fmt.Errorf("parse font collection: %w", err)
		}
		fnt, err = coll.Font(0)
		if err != nil {
			return 0, fmt.Errorf("get font 0 from collection: %w", err)
		}
	} else {
		fnt, err = opentype.Parse(fontData)
		if err != nil {
			return 0, fmt.Errorf("parse font: %w", err)
		}
	}

	fontMutex.Lock()
	defer fontMutex.Unlock()

	id := fontIDCounter
	fontIDCounter++

	fontRegistry[id] = &fontEntry{
		path:  path,
		fnt:   fnt,
		faces: make(map[int]font.Face),
	}

	return id, nil
}

// FontHasGlyph returns true if the font can render the first rune of ch.
func FontHasGlyph(fontID int, ch string) bool {
	fontMutex.Lock()
	entry, ok := fontRegistry[fontID]
	fontMutex.Unlock()

	if !ok {
		return false
	}

	if ch == "" {
		return false
	}

	r := []rune(ch)[0]
	idx, err := entry.fnt.GlyphIndex(&entry.buf, r)
	return err == nil && idx != 0
}

// FontRasterizeCell rasterizes a single glyph into a cellW×cellH alpha grid.
// Returns row-major float64 array with values 0..255 (alpha channel).
// If the glyph is too wide, scales the face size down to fit.
func FontRasterizeCell(fontID int, ch string, cellW, cellH int) ([]float64, error) {
	fontMutex.Lock()
	entry, ok := fontRegistry[fontID]
	fontMutex.Unlock()

	if !ok {
		return nil, fmt.Errorf("font ID %d not found", fontID)
	}

	if ch == "" {
		return nil, fmt.Errorf("empty character")
	}

	if cellW <= 0 || cellH <= 0 {
		return nil, fmt.Errorf("invalid cell dimensions: %dx%d", cellW, cellH)
	}

	r := []rune(ch)[0]

	// Check if font has the glyph
	idx, err := entry.fnt.GlyphIndex(&entry.buf, r)
	if err != nil || idx == 0 {
		// Return blank cell (all zeros)
		return make([]float64, cellW*cellH), nil
	}

	// Try to create/reuse a face at cellH size
	face := getFace(entry, cellH)

	// Measure the glyph at cellH
	d := &font.Drawer{Face: face}
	advance := d.MeasureString(ch)
	advancePx := int(advance >> 6) // fixed.Int26_6 to pixels

	// If glyph is too wide, scale down
	if advancePx > cellW && advancePx > 0 {
		newSize := float64(cellH) * float64(cellW) / float64(advancePx)
		face = getFace(entry, int(newSize))
	}

	// Create cell image
	cellImg := image.NewRGBA(image.Rect(0, 0, cellW, cellH))
	// Fill with transparent black (zeros) - default is already zero-initialized

	// Draw glyph
	d = &font.Drawer{
		Dst:  cellImg,
		Src:  image.White,
		Face: face,
		// Baseline ~3px above the cell bottom
		Dot: fixed.Point26_6{
			X: fixed.Int26_6(0),
			Y: fixed.Int26_6((cellH - 3) * 64),
		},
	}
	d.DrawString(ch)

	// Extract alpha channel to float64 array
	alphas := make([]float64, cellW*cellH)
	for y := 0; y < cellH; y++ {
		for x := 0; x < cellW; x++ {
			_, _, _, a := cellImg.At(x, y).RGBA()
			// RGBA returns 0..65535 range; convert to 0..255
			alphas[y*cellW+x] = float64(a >> 8)
		}
	}

	return alphas, nil
}

// SaveGlyphAtlasPNG saves a glyph atlas as a PNG file.
// alphas: row-major array of alpha values (0..255), must have length w*h.
func SaveGlyphAtlasPNG(path string, w, h int, alphas []float64) error {
	if len(alphas) != w*h {
		return fmt.Errorf("alphas length %d does not match w*h=%d", len(alphas), w*h)
	}

	if w <= 0 || h <= 0 {
		return fmt.Errorf("invalid atlas dimensions: %dx%d", w, h)
	}

	// Create RGBA image (white with alpha from alphas)
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for i, alpha := range alphas {
		y := i / w
		x := i % w
		a := uint8(alpha)
		// White with alpha
		img.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: a})
	}

	// Save to PNG
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create PNG file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("encode PNG: %w", err)
	}

	return nil
}

// getFace gets or creates a face at the given size.
func getFace(entry *fontEntry, size int) font.Face {
	fontMutex.Lock()
	defer fontMutex.Unlock()

	if f, ok := entry.faces[size]; ok {
		return f
	}

	f, _ := opentype.NewFace(entry.fnt, &opentype.FaceOptions{
		Size: float64(size),
		DPI:  72,
	})

	entry.faces[size] = f
	return f
}
