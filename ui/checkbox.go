package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Checkbox struct {
	base  Base
	theme *Theme

	label   string
	checked bool
}

func NewCheckbox(label string) *Checkbox {
	return &Checkbox{
		base:  NewBase(),
		label: label,
	}
}

func (c *Checkbox) Base() *Base     { return &c.base }
func (c *Checkbox) Focusable() bool { return true }

func (c *Checkbox) SetRectByWidth(x, y, w int) {
	if c.theme != nil {
		c.base.SetRectByWidth(c.theme, x, y, w)
		return
	}
	c.base.Rect = Rect{X: x, Y: y, W: w, H: 0}
}

func (c *Checkbox) SetEnabled(v bool) { c.base.SetEnabled(v) }
func (c *Checkbox) SetVisible(v bool) { c.base.SetVisible(v) }

func (c *Checkbox) SetChecked(v bool) { c.checked = v }
func (c *Checkbox) Checked() bool     { return c.checked }

func (c *Checkbox) HandleEvent(ctx *Context, e Event) {
	if !c.base.Enabled {
		return
	}
	if e.Type == EventClick {
		c.checked = !c.checked
	}
}

func (c *Checkbox) Update(ctx *Context) {
	c.theme = ctx.Theme
	if c.base.Rect.H == 0 {
		c.base.SetRectByWidth(ctx.Theme, c.base.Rect.X, c.base.Rect.Y, c.base.Rect.W)
	}

	if !c.base.Enabled {
		return
	}

	// Keyboard toggle
	if c.base.focused && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		c.checked = !c.checked
	}
}

func (c *Checkbox) Draw(ctx *Context, dst *ebiten.Image) {
	c.theme = ctx.Theme
	if c.base.Rect.H == 0 {
		c.base.SetRectByWidth(ctx.Theme, c.base.Rect.X, c.base.Rect.Y, c.base.Rect.W)
	}

	r := c.base.Rect

	// Surface
	bg := ctx.Theme.Surface
	if !c.base.Enabled {
		bg = ctx.Theme.SurfacePressed
	} else if c.base.pressed {
		bg = ctx.Theme.SurfacePressed
	} else if c.base.hovered {
		bg = ctx.Theme.SurfaceHover
	}
	drawRoundedRect(dst, r, ctx.Theme.Radius, bg)
	drawRoundedBorder(dst, r, ctx.Theme.Radius, ctx.Theme.BorderW, ctx.Theme.Border)

	if c.base.focused && c.base.Enabled {
		drawFocusRing(dst, r, ctx.Theme.Radius, ctx.Theme.FocusRingGap, ctx.Theme.FocusRingW, ctx.Theme.Focus)
	}

	// Checkbox box (left)
	content := r.Inset(ctx.Theme.PadX, ctx.Theme.PadY)
	boxSize := ctx.Theme.CheckSize
	if boxSize < 10 {
		boxSize = 10
	}
	boxY := r.Y + (r.H-boxSize)/2
	box := Rect{X: content.X, Y: boxY, W: boxSize, H: boxSize}

	drawRoundedRect(dst, box, int(float64(boxSize)*0.25), ctx.Theme.Bg)
	drawRoundedBorder(dst, box, int(float64(boxSize)*0.25), ctx.Theme.BorderW, ctx.Theme.Border)

	if c.checked {
		// Draw a clean checkmark (two strokes), proportional.
		x1 := float32(box.X) + float32(box.W)*0.22
		y1 := float32(box.Y) + float32(box.H)*0.55
		x2 := float32(box.X) + float32(box.W)*0.43
		y2 := float32(box.Y) + float32(box.H)*0.73
		x3 := float32(box.X) + float32(box.W)*0.78
		y3 := float32(box.Y) + float32(box.H)*0.28

		w := float32(ctx.Theme.BorderW)
		if w < 2 {
			w = 2
		}
		vector.StrokeLine(dst, x1, y1, x2, y2, w, ctx.Theme.Focus, false)
		vector.StrokeLine(dst, x2, y2, x3, y3, w, ctx.Theme.Focus, false)
	}

	// Label
	met, _ := MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Y + (r.H-met.Height)/2 + met.Ascent
	tx := box.Right() + ctx.Theme.SpaceS

	col := ctx.Theme.Text
	if !c.base.Enabled {
		col = ctx.Theme.Disabled
	}

	ctx.Text.SetColor(col)
	ctx.Text.SetAlign(0) // Left
	ctx.Text.Draw(dst, c.label, tx, baselineY)
}
