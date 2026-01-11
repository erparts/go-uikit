package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// whitePixel is a 1x1 white image used to draw colored triangles.
var whitePixel = func() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}()

func setVertexColor(vs []ebiten.Vertex, col color.RGBA) {
	r := float32(col.R) / 255.0
	g := float32(col.G) / 255.0
	b := float32(col.B) / 255.0
	a := float32(col.A) / 255.0
	for i := range vs {
		vs[i].ColorR = r
		vs[i].ColorG = g
		vs[i].ColorB = b
		vs[i].ColorA = a
	}
}

// roundedRectPath builds a rounded-rectangle path using quadratic curves.
func roundedRectPath(r Rect, radius int) vector.Path {
	var p vector.Path
	if r.W <= 0 || r.H <= 0 {
		return p
	}

	rr := float32(radius)
	if rr < 0 {
		rr = 0
	}
	if rr > float32(r.W)/2 {
		rr = float32(r.W) / 2
	}
	if rr > float32(r.H)/2 {
		rr = float32(r.H) / 2
	}

	x0 := float32(r.X)
	y0 := float32(r.Y)
	x1 := float32(r.Right())
	y1 := float32(r.Bottom())

	// Start at top-left corner (after radius).
	p.MoveTo(x0+rr, y0)
	p.LineTo(x1-rr, y0)
	p.QuadTo(x1, y0, x1, y0+rr)
	p.LineTo(x1, y1-rr)
	p.QuadTo(x1, y1, x1-rr, y1)
	p.LineTo(x0+rr, y1)
	p.QuadTo(x0, y1, x0, y1-rr)
	p.LineTo(x0, y0+rr)
	p.QuadTo(x0, y0, x0+rr, y0)
	p.Close()

	return p
}

// drawRoundedRect fills a rounded rectangle.
func drawRoundedRect(dst *ebiten.Image, r Rect, radius int, col color.RGBA) {
	if r.W <= 0 || r.H <= 0 {
		return
	}
	if radius <= 0 {
		vector.DrawFilledRect(dst, float32(r.X), float32(r.Y), float32(r.W), float32(r.H), col, false)
		return
	}

	p := roundedRectPath(r, radius)
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	setVertexColor(vs, col)
	dst.DrawTriangles(vs, is, whitePixel, nil)
}

func drawRoundedBorder(dst *ebiten.Image, r Rect, radius int, borderW int, col color.RGBA) {
	if borderW <= 0 || r.W <= 0 || r.H <= 0 {
		return
	}
	p := roundedRectPath(r, radius)

	// Stroke geometry
	op := &vector.StrokeOptions{
		Width: float32(borderW),
	}
	vs, is := p.AppendVerticesAndIndicesForStroke(nil, nil, op)
	setVertexColor(vs, col)
	dst.DrawTriangles(vs, is, whitePixel, nil)
}

func drawFocusRing(dst *ebiten.Image, r Rect, radius int, gap int, w int, col color.RGBA) {
	if w <= 0 {
		return
	}
	rr := Rect{
		X: r.X - gap - w,
		Y: r.Y - gap - w,
		W: r.W + (gap+w)*2,
		H: r.H + (gap+w)*2,
	}
	drawRoundedBorder(dst, rr, radius+gap+w, w, col)
}

func clampInt(v, lo, hi int) int {
	return int(math.Max(float64(lo), math.Min(float64(hi), float64(v))))
}
