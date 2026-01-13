package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Container is an empty widget that lets you render custom content inside a themed box.
// It still participates in focus/invalid layout like any other widget.
type Container struct {
	base  Base
	theme *Theme

	// If nil, no surface is drawn (only content callback).
	DrawSurface bool

	OnUpdate func(ctx *Context, content Rect)
	OnDraw   func(ctx *Context, dst *ebiten.Image, content Rect)
}

func NewContainer() *Container {
	return &Container{
		base:        NewBase(&WidgetBaseConfig{}),
		DrawSurface: true,
	}
}

func (c *Container) Base() *Base     { return &c.base }
func (c *Container) Focusable() bool { return false }

func (c *Container) SetFrame(x, y, w int) {
	if c.theme != nil {
		c.base.SetFrame(c.theme, x, y, w)
		return
	}

	c.base.Rect = Rect{X: x, Y: y, W: w, H: 0}
}

func (c *Container) Measure() Rect { return c.base.Rect }

func (c *Container) Update(ctx *Context) {
	c.theme = ctx.Theme
	if c.base.Rect.H == 0 {
		c.base.SetFrame(ctx.Theme, c.base.Rect.X, c.base.Rect.Y, c.base.Rect.W)
	}

	if c.OnUpdate != nil {
		c.OnUpdate(ctx, c.base.ControlRect(ctx.Theme).Inset(ctx.Theme.PadX, ctx.Theme.PadY))
	}
}

func (c *Container) Draw(ctx *Context, dst *ebiten.Image) {
	c.theme = ctx.Theme
	if c.base.Rect.H == 0 {
		c.base.SetFrame(ctx.Theme, c.base.Rect.X, c.base.Rect.Y, c.base.Rect.W)
	}

	r := c.base.ControlRect(ctx.Theme)

	if c.DrawSurface {
		drawRoundedRect(dst, r, ctx.Theme.Radius, ctx.Theme.Surface)
		borderCol := ctx.Theme.Border
		if c.base.invalid {
			borderCol = ctx.Theme.ErrorBorder
		}

		drawRoundedBorder(dst, r, ctx.Theme.Radius, ctx.Theme.BorderW, borderCol)
	}

	content := r.Inset(ctx.Theme.PadX, ctx.Theme.PadY)
	if c.OnDraw != nil {
		c.OnDraw(ctx, dst, content)
	}

	err := c.base.ErrorRect(ctx.Theme)
	if c.base.invalid {
		drawErrorText(ctx, dst, err, c.base.errorText)
	}
}

// SetTheme allows layouts to provide Theme before SetFrame is called.
func (c *Container) SetTheme(theme *Theme) { c.theme = theme }
