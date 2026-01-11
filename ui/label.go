package ui

import "github.com/hajimehoshi/ebiten/v2"

type Label struct {
	base  Base
	text  string
	theme *Theme
}

func NewLabel(text string) *Label {
	return &Label{base: NewBase(), text: text}
}

func (l *Label) Base() *Base       { return &l.base }
func (l *Label) Focusable() bool   { return false }
func (l *Label) SetText(s string)  { l.text = s }
func (l *Label) SetEnabled(v bool) { l.base.SetEnabled(v) }
func (l *Label) SetVisible(v bool) { l.base.SetVisible(v) }
func (l *Label) SetRectByWidth(x, y, w int) {
	// Height is derived from theme; if theme not known yet, keep H=0 and
	// it will be fixed on first Update/Draw.
	if l.theme != nil {
		l.base.SetRectByWidth(l.theme, x, y, w)
		return
	}
	l.base.Rect = Rect{X: x, Y: y, W: w, H: 0}
}

func (l *Label) Update(ctx *Context) {
	l.theme = ctx.Theme
	if l.base.Rect.H == 0 {
		l.base.SetRectByWidth(ctx.Theme, l.base.Rect.X, l.base.Rect.Y, l.base.Rect.W)
	}
}

func (l *Label) Draw(ctx *Context, dst *ebiten.Image) {
	l.theme = ctx.Theme
	if l.base.Rect.H == 0 {
		l.base.SetRectByWidth(ctx.Theme, l.base.Rect.X, l.base.Rect.Y, l.base.Rect.W)
	}

	r := l.base.Rect
	met, _ := MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Y + (r.H-met.Height)/2 + met.Ascent

	ctx.Text.SetColor(ctx.Theme.MutedText)
	ctx.Text.SetAlign(0) // Left
	ctx.Text.Draw(dst, l.text, r.X, baselineY)
}
