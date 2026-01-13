package widget

import (
	"github.com/erparts/go-uikit/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Checkbox struct {
	base  ui.Base
	theme *ui.Theme

	label   string
	checked bool
}

func NewCheckbox(label string) *Checkbox {
	cfg := ui.NewWidgetBaseConfig()

	return &Checkbox{
		base:  ui.NewBase(cfg),
		label: label,
	}
}

func (c *Checkbox) Base() *ui.Base  { return &c.base }
func (c *Checkbox) Focusable() bool { return true }

func (c *Checkbox) SetFrame(x, y, w int) {
	if c.theme != nil {
		c.base.SetFrame(c.theme, x, y, w)
		return
	}

	c.base.Rect = ui.Rect{X: x, Y: y, W: w, H: 0}
}

func (c *Checkbox) Measure() ui.Rect { return c.base.Rect }

func (c *Checkbox) SetEnabled(v bool) { c.base.SetEnabled(v) }
func (c *Checkbox) SetVisible(v bool) { c.base.SetVisible(v) }

func (c *Checkbox) SetChecked(v bool) { c.checked = v }
func (c *Checkbox) Checked() bool     { return c.checked }

func (c *Checkbox) HandleEvent(ctx *ui.Context, e ui.Event) {
	if !c.base.IsEnabled() {
		return
	}
	if e.Type == ui.EventClick {
		c.checked = !c.checked
	}
}

func (c *Checkbox) Update(ctx *ui.Context) {
	c.theme = ctx.Theme
	if c.base.Rect.H == 0 {
		c.base.SetFrame(ctx.Theme, c.base.Rect.X, c.base.Rect.Y, c.base.Rect.W)
	}

	if !c.base.IsEnabled() {
		return
	}

	// Keyboard toggle
	if c.base.Focused() && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		c.checked = !c.checked
	}
}

func (c *Checkbox) Draw(ctx *ui.Context, dst *ebiten.Image) {
	c.base.Draw(ctx, dst)

	r := c.base.ControlRect(ctx.Theme)

	// Checkbox box (left)
	content := r.Inset(ctx.Theme.PadX, ctx.Theme.PadY)
	boxSize := ctx.Theme.CheckSize
	if boxSize < 10 {
		boxSize = 10
	}
	boxY := r.Y + (r.H-boxSize)/2
	box := ui.Rect{X: content.X, Y: boxY, W: boxSize, H: boxSize}

	c.base.DrawRoundedRect(dst, box, int(float64(boxSize)*0.25), ctx.Theme.Bg)
	c.base.DrawRoundedBorder(dst, box, int(float64(boxSize)*0.25), ctx.Theme.BorderW, ctx.Theme.Border)

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
	met, _ := ui.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Y + (r.H-met.Height)/2 + met.Ascent
	tx := box.Right() + ctx.Theme.SpaceS

	col := ctx.Theme.Text
	if !c.base.IsEnabled() {
		col = ctx.Theme.Disabled
	}

	ctx.Text.SetColor(col)
	ctx.Text.SetAlign(0) // Left
	ctx.Text.Draw(dst, c.label, tx, baselineY)
}

// SetTheme allows layouts to provide Theme before SetFrame is called.
func (c *Checkbox) SetTheme(theme *ui.Theme) { c.theme = theme }
