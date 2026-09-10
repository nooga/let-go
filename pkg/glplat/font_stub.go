//go:build !glplat && !glplat_ebiten

/*
 * font_stub — untagged stand-ins for the font primitives in font.go
 *
 * The real implementations import the x/image/font and x/text/encoding
 * closure, which check-default-deps keeps out of untagged builds (see the
 * header comment in font.go). Package rt links the generated glplat interop
 * unconditionally, so these entry points must still exist without a backend
 * tag; they report that the primitives were not built rather than silently
 * returning an empty atlas.
 */

package glplat

// FontLoad reports that the font primitives were not built.
func FontLoad(path string) (int, error) {
	return 0, errFontNotBuilt("FontLoad")
}

// FontHasGlyph reports false, matching the backend-less accessors in glplat.go
// that degrade to a zero value rather than an error.
func FontHasGlyph(fontID int, ch string) bool {
	return false
}

// FontRasterizeCell reports that the font primitives were not built.
func FontRasterizeCell(fontID int, ch string, cellW, cellH int) ([]float64, error) {
	return nil, errFontNotBuilt("FontRasterizeCell")
}

// SaveGlyphAtlasPNG reports that the font primitives were not built.
func SaveGlyphAtlasPNG(path string, w, h int, alphas []float64) error {
	return errFontNotBuilt("SaveGlyphAtlasPNG")
}

// errFontNotBuilt mirrors errNoBackend's shape, naming the build tags that
// bring the primitives in.
func errFontNotBuilt(op string) error {
	return &fontNotBuiltError{op: op}
}

type fontNotBuiltError struct {
	op string
}

func (e *fontNotBuiltError) Error() string {
	return "glplat: font primitives not built for " + e.op +
		"; build with -tags glplat or -tags glplat_ebiten"
}
