package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Button struct {
	base  Base
	theme *Theme

	label   string
	OnClick func()
}

func NewButton(label string) *Button {
	return &Button{
		base:  NewBase(),
		label: label,
	}
}

func (b *Button) Base() *Base     { return &b.base }
func (b *Button) Focusable() bool { return true }

func (b *Button) SetRectByWidth(x, y, w int) {
	if b.theme != nil {
		b.base.SetRectByWidth(b.theme, x, y, w)
		return
	}
	b.base.Rect = Rect{X: x, Y: y, W: w, H: 0}
}

func (b *Button) SetEnabled(v bool) { b.base.SetEnabled(v) }
func (b *Button) SetVisible(v bool) { b.base.SetVisible(v) }
func (b *Button) SetLabel(s string) { b.label = s }

func (b *Button) HandleEvent(ctx *Context, e Event) {
	if !b.base.Enabled {
		return
	}
	if e.Type == EventClick {
		if b.OnClick != nil {
			b.OnClick()
		}
	}
}

func (b *Button) Update(ctx *Context) {
	b.theme = ctx.Theme
	if b.base.Rect.H == 0 {
		b.base.SetRectByWidth(ctx.Theme, b.base.Rect.X, b.base.Rect.Y, b.base.Rect.W)
	}

	if !b.base.Enabled {
		return
	}

	// Keyboard click when focused
	if b.base.focused && (inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)) {
		if b.OnClick != nil {
			b.OnClick()
		}
	}
}

func (b *Button) Draw(ctx *Context, dst *ebiten.Image) {
	b.theme = ctx.Theme
	if b.base.Rect.H == 0 {
		b.base.SetRectByWidth(ctx.Theme, b.base.Rect.X, b.base.Rect.Y, b.base.Rect.W)
	}

	r := b.base.Rect

	// Surface
	bg := ctx.Theme.Surface
	if !b.base.Enabled {
		bg = ctx.Theme.SurfacePressed
	} else if b.base.pressed {
		bg = ctx.Theme.SurfacePressed
	} else if b.base.hovered {
		bg = ctx.Theme.SurfaceHover
	}
	drawRoundedRect(dst, r, ctx.Theme.Radius, bg)

	// Border
	border := ctx.Theme.Border
	if !b.base.Enabled {
		border = ctx.Theme.Disabled
	}
	drawRoundedBorder(dst, r, ctx.Theme.Radius, ctx.Theme.BorderW, border)

	// Focus ring
	if b.base.focused && b.base.Enabled {
		drawFocusRing(dst, r, ctx.Theme.Radius, ctx.Theme.FocusRingGap, ctx.Theme.FocusRingW, ctx.Theme.Focus)
	}

	// Centered label
	met, _ := MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	textW := MeasureStringPx(ctx.Theme.Font, ctx.Theme.FontPx, b.label)

	tx := r.X + (r.W-textW)/2
	baselineY := r.Y + (r.H-met.Height)/2 + met.Ascent

	col := ctx.Theme.Text
	if !b.base.Enabled {
		col = ctx.Theme.Disabled
	}

	ctx.Text.SetColor(col)
	ctx.Text.SetAlign(0) // Left
	ctx.Text.Draw(dst, b.label, tx, baselineY)
}
