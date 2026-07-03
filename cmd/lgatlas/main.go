/*
 * lgatlas — font atlas baking tool for glplat
 *
 * Rasterizes a font into a texture atlas with metrics.
 */

package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func main() {
	fontPath := flag.String("font", "", "path to TTF font file")
	cellW := flag.Int("cell-w", 8, "cell width (half-width cells)")
	cellH := flag.Int("cell-h", 16, "cell height")
	runesPath := flag.String("runes", "", "path to UTF-8 file with runes to include")
	outPng := flag.String("out-png", "atlas.png", "output PNG file")
	outEdn := flag.String("out-edn", "atlas.edn", "output EDN metrics file")
	flag.Parse()

	if *fontPath == "" {
		log.Fatal("missing required flag: -font")
	}
	if *runesPath == "" {
		log.Fatal("missing required flag: -runes")
	}

	// Load font
	fontData, err := os.ReadFile(*fontPath)
	if err != nil {
		log.Fatalf("read font: %v", err)
	}

	fnt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatalf("parse font: %v", err)
	}

	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size: float64(*cellH),
		DPI:  72,
	})
	if err != nil {
		log.Fatalf("create face: %v", err)
	}
	defer face.Close()

	// Load runes to include
	runesData, err := os.ReadFile(*runesPath)
	if err != nil {
		log.Fatalf("read runes file: %v", err)
	}

	runes := []rune(string(runesData))
	fmt.Printf("Baking atlas for %d runes\n", len(runes))

	// Bake atlas
	metrics, img, err := bakeAtlas(face, runes, *cellW, *cellH)
	if err != nil {
		log.Fatalf("bake atlas: %v", err)
	}

	// Save PNG
	pngFile, err := os.Create(*outPng)
	if err != nil {
		log.Fatalf("create PNG: %v", err)
	}
	defer pngFile.Close()

	if err := png.Encode(pngFile, img); err != nil {
		log.Fatalf("encode PNG: %v", err)
	}
	fmt.Printf("Saved atlas to %s (%dx%d)\n", *outPng, img.Bounds().Dx(), img.Bounds().Dy())

	// Save EDN metrics
	ednFile, err := os.Create(*outEdn)
	if err != nil {
		log.Fatalf("create EDN: %v", err)
	}
	defer ednFile.Close()

	if err := writeMetricsEDN(ednFile, metrics); err != nil {
		log.Fatalf("write EDN: %v", err)
	}
	fmt.Printf("Saved metrics to %s\n", *outEdn)
}

type GlyphMetrics struct {
	Char string
	U0   float64
	V0   float64
	U1   float64
	V1   float64
}

func bakeAtlas(face font.Face, runes []rune, cellW, cellH int) ([]GlyphMetrics, image.Image, error) {
	// Determine grid dimensions (keep near-square, power-of-2 preferred)
	gridW := 16
	gridH := 1
	for gridW*gridH < len(runes) {
		if gridW > gridH {
			gridH *= 2
		} else {
			gridW *= 2
		}
	}

	texW := gridW * cellW
	texH := gridH * cellH

	// Make sure dimensions are power of 2
	texW = nextPowerOf2(texW)
	texH = nextPowerOf2(texH)

	fmt.Printf("Atlas grid: %dx%d cells = %dx%d pixels\n", gridW, gridH, texW, texH)

	// Create atlas image (white on transparent)
	atlasImg := image.NewRGBA(image.Rect(0, 0, texW, texH))
	// Image is initialized to transparent (all zeros) by default

	metrics := make([]GlyphMetrics, 0, len(runes))

	for idx, r := range runes {
		gridX := idx % gridW
		gridY := idx / gridW
		if gridY >= gridH {
			fmt.Printf("Warning: too many runes for grid (%d > %d*%d), skipping\n", len(runes), gridW, gridH)
			break
		}

		cellX := gridX * cellW
		cellY := gridY * cellH

		// Rasterize glyph into the cell
		drawer := &font.Drawer{
			Dst:  atlasImg,
			Src:  image.White,
			Face: face,
			Dot:  fixed.Point26_6{X: fixed.Int26_6(cellX * 64), Y: fixed.Int26_6((cellY + cellH) * 64)},
		}
		drawer.DrawString(string(r))

		// Calculate UV coordinates (v increases downward as per contract)
		u0 := float64(cellX) / float64(texW)
		v0 := float64(cellY) / float64(texH)
		u1 := float64(cellX+cellW) / float64(texW)
		v1 := float64(cellY+cellH) / float64(texH)

		metrics = append(metrics, GlyphMetrics{
			Char: string(r),
			U0:   u0,
			V0:   v0,
			U1:   u1,
			V1:   v1,
		})
	}

	return metrics, atlasImg, nil
}

func writeMetricsEDN(w io.Writer, metrics []GlyphMetrics) error {
	fmt.Fprint(w, "{")
	for i, m := range metrics {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, "%q [%.4f %.4f %.4f %.4f]",
			m.Char, m.U0, m.V0, m.U1, m.V1)
	}
	fmt.Fprint(w, "}\n")
	return nil
}

func nextPowerOf2(n int) int {
	p := 1
	for p < n {
		p *= 2
	}
	return p
}
