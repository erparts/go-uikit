package widget

import (
	"image"

	"github.com/erparts/go-uikit"
	"github.com/erparts/go-uikit/common"
	"github.com/hajimehoshi/ebiten/v2"
)

// Container is an empty widget that lets you render custom content inside a themed box.
// It still participates in focus/invalid layout like any other widget.
type Container struct {
	uikit.Base
	height int

	OnUpdate func(ctx *uikit.Context, content image.Rectangle)
	OnDraw   func(ctx *uikit.Context, dst *ebiten.Image, content image.Rectangle)
}

func NewContainer(theme *uikit.Theme) *Container {
	cfg := uikit.NewWidgetBaseConfig(theme)

	w := &Container{}
	w.Base = uikit.NewBase(cfg)
	w.Base.HeightCaculator = func() int {
		return w.height
	}

	return w
}

func (c *Container) SetHeight(h int) {
	c.height = h
}

func (c *Container) Focusable() bool { return false }

func (c *Container) Update(ctx *uikit.Context) {
	r := c.Measure(false)
	if r.Dy() == 0 {
		c.SetFrame(r.Min.X, r.Min.Y, r.Dx())
	}

	if c.OnUpdate != nil {
		c.OnUpdate(ctx, common.Inset(c.Measure(false), ctx.Theme.PadX, ctx.Theme.PadY))
	}
}

func (c *Container) Draw(ctx *uikit.Context, dst *ebiten.Image) {
	r := c.Base.Draw(ctx, dst)

	content := common.Inset(r, ctx.Theme.PadX, ctx.Theme.PadY)
	if c.OnDraw != nil {
		c.OnDraw(ctx, dst, content)
	}
}
