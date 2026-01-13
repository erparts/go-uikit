package widget

import (
	"github.com/erparts/go-uikit/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

// Container is an empty widget that lets you render custom content inside a themed box.
// It still participates in focus/invalid layout like any other widget.
type Container struct {
	base  ui.Base
	theme *ui.Theme

	OnUpdate func(ctx *ui.Context, content ui.Rect)
	OnDraw   func(ctx *ui.Context, dst *ebiten.Image, content ui.Rect)
}

func NewContainer() *Container {
	cfg := ui.NewWidgetBaseConfig()

	return &Container{
		base: ui.NewBase(cfg),
	}
}

func (c *Container) Base() *ui.Base  { return &c.base }
func (c *Container) Focusable() bool { return false }

func (c *Container) SetFrame(x, y, w int) {
	if c.theme != nil {
		c.base.SetFrame(c.theme, x, y, w)
		return
	}

	c.base.Rect = ui.Rect{X: x, Y: y, W: w, H: 0}
}

func (c *Container) Measure() ui.Rect { return c.base.Rect }

func (c *Container) Update(ctx *ui.Context) {
	c.theme = ctx.Theme
	if c.base.Rect.H == 0 {
		c.base.SetFrame(ctx.Theme, c.base.Rect.X, c.base.Rect.Y, c.base.Rect.W)
	}

	if c.OnUpdate != nil {
		c.OnUpdate(ctx, c.base.ControlRect(ctx.Theme).Inset(ctx.Theme.PadX, ctx.Theme.PadY))
	}
}

func (c *Container) Draw(ctx *ui.Context, dst *ebiten.Image) {
	r := c.base.Draw(ctx, dst)

	content := r.Inset(ctx.Theme.PadX, ctx.Theme.PadY)
	if c.OnDraw != nil {
		c.OnDraw(ctx, dst, content)
	}
}

// SetTheme allows layouts to provide Theme before SetFrame is called.
func (c *Container) SetTheme(theme *ui.Theme) { c.theme = theme }
