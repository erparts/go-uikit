package ui

import (
	"image/color"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// Theme is the single source of truth for widget proportions.
// Only the font and font size are inputs; everything else derives from them.
type Theme struct {
	Font     *sfnt.Font
	FontPx   int
	ControlH int

	PadX int
	PadY int

	Radius       int
	BorderW      int
	FocusRingW   int
	FocusRingGap int

	SpaceS int
	SpaceM int
	SpaceL int

	CheckSize int // checkbox square size

	// Validation
	ErrorFontPx int
	ErrorGap    int

	// Colors
	Text           color.RGBA
	MutedText      color.RGBA
	Bg             color.RGBA
	Surface        color.RGBA
	SurfaceHover   color.RGBA
	SurfacePressed color.RGBA
	Border         color.RGBA
	Focus          color.RGBA
	Disabled       color.RGBA
	ErrorText      color.RGBA
	ErrorBorder    color.RGBA
	Scrollbar      color.RGBA

	// Scrollbar
	ScrollbarRadius int
	ScrollbarMinH   int
	ScrollbarW      int

	Caret color.RGBA
}

type FontMetricsPx struct {
	Ascent  int
	Descent int
	Height  int
}

// MetricsPx returns font metrics in pixels for a given size (ppem).
func MetricsPx(f *sfnt.Font, sizePx int) (FontMetricsPx, error) {
	var buf sfnt.Buffer
	m, err := f.Metrics(&buf, fixed.I(sizePx), font.HintingNone)
	if err != nil {
		return FontMetricsPx{}, err
	}
	ascent := int(m.Ascent.Round())
	descent := int(m.Descent.Round())
	height := int(m.Height.Round())
	if height <= 0 {
		height = ascent + descent
	}
	return FontMetricsPx{Ascent: ascent, Descent: descent, Height: height}, nil
}

func NewTheme(font *sfnt.Font, fontPx int) *Theme {
	if fontPx < 10 {
		fontPx = 10
	}

	met, err := MetricsPx(font, fontPx)
	if err != nil {
		panic(err)
	}

	// Control height derived from font height.
	// k=1.8 gives a "web input" feel without being too tall.
	controlH := int(math.Round(float64(met.Height) * 1.8))
	if controlH < met.Height+6 {
		controlH = met.Height + 6
	}

	padY := (controlH - met.Height) / 2
	if padY < 2 {
		padY = 2
	}
	padX := int(math.Round(float64(padY) * 1.6))
	if padX < 6 {
		padX = 6
	}

	radius := int(math.Round(float64(controlH) * 0.22))
	if radius < 4 {
		radius = 4
	}

	borderW := int(math.Round(float64(controlH) * 0.06))
	if borderW < 1 {
		borderW = 1
	}

	focusW := borderW
	focusGap := int(math.Round(float64(borderW) * 1.5))
	if focusGap < 2 {
		focusGap = 2
	}

	spaceS := int(math.Round(float64(controlH) * 0.20))
	spaceM := int(math.Round(float64(controlH) * 0.35))
	spaceL := int(math.Round(float64(controlH) * 0.60))
	if spaceS < 4 {
		spaceS = 4
	}
	if spaceM < 8 {
		spaceM = 8
	}
	if spaceL < 12 {
		spaceL = 12
	}

	checkSize := met.Height
	if checkSize < controlH-padY*2 {
		checkSize = controlH - padY*2
	}
	// Keep it visually balanced
	checkSize = int(math.Round(float64(checkSize) * 0.92))

	// Validation
	errorFontPx := int(math.Round(float64(fontPx) * 0.85))
	if errorFontPx < 10 {
		errorFontPx = 10
	}
	errorGap := int(math.Round(float64(controlH) * 0.15))
	if errorGap < 4 {
		errorGap = 4
	}

	return &Theme{
		Font:     font,
		FontPx:   fontPx,
		ControlH: controlH,
		PadX:     padX,
		PadY:     padY,

		Radius:       radius,
		BorderW:      borderW,
		FocusRingW:   focusW,
		FocusRingGap: focusGap,

		SpaceS: spaceS,
		SpaceM: spaceM,
		SpaceL: spaceL,

		CheckSize: checkSize,

		ErrorFontPx: errorFontPx,
		ErrorGap:    errorGap,

		Text:           color.RGBA{235, 238, 242, 255},
		MutedText:      color.RGBA{170, 176, 186, 255},
		Bg:             color.RGBA{20, 22, 26, 255},
		Surface:        color.RGBA{34, 38, 46, 255},
		SurfaceHover:   color.RGBA{42, 48, 58, 255},
		SurfacePressed: color.RGBA{28, 32, 40, 255},
		Border:         color.RGBA{76, 84, 98, 255},
		Focus:          color.RGBA{120, 170, 255, 255},
		Disabled:       color.RGBA{90, 96, 106, 255},
		ErrorText:      color.RGBA{235, 110, 110, 255},
		ErrorBorder:    color.RGBA{235, 110, 110, 255},

		Caret: color.RGBA{235, 238, 242, 255},
	}
}
