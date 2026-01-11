package ui

import (
	"unicode/utf8"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// MeasureStringPx measures the advance width of a string in pixels using sfnt.
// It is simple (no shaping), but it is deterministic and matches our UI rules.
func MeasureStringPx(f *sfnt.Font, sizePx int, s string) int {
	if s == "" {
		return 0
	}
	var buf sfnt.Buffer
	ppem := fixed.I(sizePx)

	var adv fixed.Int26_6
	for len(s) > 0 {
		r, n := utf8.DecodeRuneInString(s)
		s = s[n:]
		if r == utf8.RuneError && n == 1 {
			continue
		}

		gid, err := f.GlyphIndex(&buf, r)
		if err != nil {
			continue
		}
		a, err := f.GlyphAdvance(&buf, gid, ppem, font.HintingNone)
		if err != nil {
			continue
		}
		adv += a
	}

	return int(adv.Round())
}
