package uikit

import (
	"image"
	"image/color"
	"math"

	"github.com/erparts/go-shapes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var shapesRenderer = shapes.NewRenderer()

func clampInt(v, lo, hi int) int {
	return int(math.Max(float64(lo), math.Min(float64(hi), float64(v))))
}

func drawPathOptionsForColor(col color.RGBA) *vector.DrawPathOptions {
	op := &vector.DrawPathOptions{
		AntiAlias: true,
	}
	op.ColorScale.ScaleWithColor(col)
	return op
}

func drawRoundedRect(dst *ebiten.Image, r image.Rectangle, radius int, col color.RGBA) {
	if r.Dy() <= 0 || r.Dx() <= 0 {
		return
	}

	rr := clampInt(radius, 0, min(r.Dx()/2, r.Dy()/2))
	shapesRenderer.SetColor(col)
	shapesRenderer.DrawRect(dst, r, float32(rr))
}

func drawRoundedBorder(dst *ebiten.Image, r image.Rectangle, radius int, borderW int, col color.RGBA) {
	if borderW <= 0 || r.Dx() <= 0 || r.Dy() <= 0 {
		return
	}

	rr := clampInt(radius, 0, min(r.Dx()/2, r.Dy()/2))
	bw := float32(borderW)

	shapesRenderer.SetColor(col)
	shapesRenderer.StrokeRect(dst, r, 0, bw, float32(rr))
}

func drawErrorText(ctx *Context, dst *ebiten.Image, r image.Rectangle, msg string) {
	if msg == "" || r.Dx() <= 0 || r.Dy() <= 0 {
		return
	}

	met, _ := MetricsPx(ctx.Theme.Font, ctx.Theme.ErrorFontPx)
	baselineY := r.Min.Y + met.Ascent

	ctx.Text.SetAlign(0)
	ctx.Text.SetSize(float64(ctx.Theme.ErrorFontPx))
	ctx.Text.SetColor(ctx.Theme.ErrorText)
	ctx.Text.Draw(dst, msg, r.Min.X, baselineY)

	ctx.Text.SetSize(float64(ctx.Theme.FontPx))
	ctx.Text.SetColor(ctx.Theme.Text)
}
