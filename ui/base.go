package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type WidgetBaseConfig struct {
	DrawSurface bool
	DrawBorder  bool
	DrawFocus   bool
	DrawInvalid bool
}

func NewWidgetBaseConfig() *WidgetBaseConfig {
	return &WidgetBaseConfig{
		DrawSurface: true,
		DrawBorder:  true,
		DrawFocus:   true,
		DrawInvalid: true,
	}
}

// Base contains shared widget state.
// Height is usually Theme.ControlH, but can be overridden per widget (e.g. TextArea). External layout controls only X/Y/Width.
type Base struct {
	cfg *WidgetBaseConfig

	Rect Rect

	// ControlH overrides Theme.ControlH when > 0 (for variable-height controls like TextArea).
	ControlH int

	// internal state (managed by Context)
	hovered   bool
	pressed   bool
	focused   bool
	visible   bool
	enabled   bool
	invalid   bool
	errorText string

	theme *Theme
}

func NewBase(cfg *WidgetBaseConfig) Base {
	return Base{
		cfg:     cfg,
		visible: true,
		enabled: true,
	}
}

func (b *Base) ControlHeight(theme *Theme) int {
	if b.ControlH > 0 {
		return b.ControlH
	}

	return theme.ControlH
}

func (b *Base) RequiredHeight(theme *Theme) int {
	h := b.ControlHeight(theme)
	if b.invalid && b.errorText != "" {
		// error line height uses theme.ErrorFontPx metrics
		met, _ := MetricsPx(theme.Font, theme.ErrorFontPx)
		h += theme.ErrorGap + met.Height
	}
	return h
}

func (b *Base) ControlRect(theme *Theme) Rect {
	return Rect{X: b.Rect.X, Y: b.Rect.Y, W: b.Rect.W, H: b.ControlHeight(theme)}
}

func (b *Base) ErrorRect(theme *Theme) Rect {
	if !(b.invalid && b.errorText != "") {
		return Rect{}
	}
	met, _ := MetricsPx(theme.Font, theme.ErrorFontPx)
	y := b.Rect.Y + b.ControlHeight(theme) + theme.ErrorGap
	return Rect{X: b.Rect.X, Y: y, W: b.Rect.W, H: met.Height}
}

func (b *Base) IsInvalid() (bool, string) {
	return b.invalid, b.errorText
}

func (b *Base) SetInvalid(err string) {
	b.invalid = true
	b.errorText = err
}

func (b *Base) ClearInvalid() {
	b.invalid = false
	b.errorText = ""
}

// SetFrame sets the widget position (x,y) and width (w).
// The final height is computed from the theme and the widget state (e.g. error line, variable-height controls).
func (b *Base) SetFrame(theme *Theme, x, y, w int) {
	if w < 0 {
		w = 0
	}
	b.Rect = Rect{X: x, Y: y, W: w, H: b.RequiredHeight(theme)}
}

func (b *Base) Hovered() bool {
	return b.hovered
}

func (b *Base) Pressed() bool {
	return b.pressed
}

func (b *Base) Focused() bool {
	return b.focused
}

func (b *Base) IsEnabled() bool {
	return b.enabled
}

func (b *Base) SetEnabled(v bool) {
	b.enabled = v
}

func (b *Base) IsVisible() bool {
	return b.visible
}

func (b *Base) SetVisible(v bool) {
	b.visible = v
}

func (c *Base) Draw(ctx *Context, dst *ebiten.Image) Rect {
	c.theme = ctx.Theme
	if c.Rect.H == 0 {
		c.SetFrame(ctx.Theme, c.Rect.X, c.Rect.Y, c.Rect.W)
	}

	c.DrawSurfece(ctx, dst)
	c.DrawBoder(ctx, dst)
	c.DrawFocus(ctx, dst)
	c.DrawInvalid(ctx, dst)
	return c.ControlRect(ctx.Theme)
}

func (c *Base) DrawSurfece(ctx *Context, dst *ebiten.Image) {
	if !c.cfg.DrawSurface {
		return
	}

	bg := ctx.Theme.Surface
	if !c.enabled {
		bg = ctx.Theme.SurfacePressed
	} else if c.pressed {
		bg = ctx.Theme.SurfacePressed
	} else if c.hovered {
		bg = ctx.Theme.SurfaceHover
	}

	r := c.ControlRect(ctx.Theme)
	drawRoundedRect(dst, r, ctx.Theme.Radius, bg)
}

func (c *Base) DrawBoder(ctx *Context, dst *ebiten.Image) {
	if !c.cfg.DrawBorder {
		return
	}

	border := ctx.Theme.Border
	if !c.enabled {
		border = ctx.Theme.Disabled
	}
	if c.invalid {
		border = ctx.Theme.ErrorBorder
	}

	r := c.ControlRect(ctx.Theme)
	drawRoundedBorder(dst, r, ctx.Theme.Radius, ctx.Theme.BorderW, border)
}

func (c *Base) DrawFocus(ctx *Context, dst *ebiten.Image) {
	if !c.cfg.DrawFocus {
		return
	}

	if !c.focused || !c.enabled {
		return
	}

	r := c.ControlRect(ctx.Theme)
	drawFocusRing(dst, r, ctx.Theme.Radius, ctx.Theme.FocusRingGap, ctx.Theme.FocusRingW, ctx.Theme.Focus)
}

func (c *Base) DrawInvalid(ctx *Context, dst *ebiten.Image) {
	if !c.cfg.DrawInvalid {
		return
	}

	if !c.invalid {
		return
	}

	err := c.ErrorRect(ctx.Theme)
	drawErrorText(ctx, dst, err, c.errorText)
}

// SetTheme allows layouts to provide Theme before SetFrame is called.
func (c *Base) SetTheme(theme *Theme) {
	c.theme = theme
}

func (c *Base) Theme() *Theme {
	return c.theme
}

func (c *Base) DrawRoundedRect(dst *ebiten.Image, r Rect, radius int, col color.RGBA) {
	drawRoundedRect(dst, r, radius, col)
}

func (c *Base) DrawRoundedBorder(dst *ebiten.Image, r Rect, radius int, borderW int, col color.RGBA) {
	drawRoundedBorder(dst, r, radius, borderW, col)
}
