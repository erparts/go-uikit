package widget

import (
	"image"

	"github.com/erparts/go-uikit"
	"github.com/erparts/go-uikit/common"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Checkbox struct {
	uikit.Base

	label   string
	checked bool
}

func NewCheckbox(theme *uikit.Theme, label string) *Checkbox {
	cfg := uikit.NewWidgetBaseConfig(theme)

	return &Checkbox{
		Base:  uikit.NewBase(cfg),
		label: label,
	}
}

func (c *Checkbox) Focusable() bool { return true }

func (c *Checkbox) SetChecked(v bool) { c.checked = v }
func (c *Checkbox) Checked() bool     { return c.checked }

func (c *Checkbox) HandleEvent(ctx *uikit.Context, e uikit.Event) {
	if !c.IsEnabled() {
		return
	}

	if e.Type == uikit.EventClick {
		c.checked = !c.checked
	}
}

func (c *Checkbox) Update(ctx *uikit.Context) {
	r := c.Measure(false)
	if r.Dy() == 0 {
		c.SetFrame(r.Min.X, r.Min.Y, r.Dx())
	}

	if !c.IsEnabled() {
		return
	}

	if c.IsFocused() && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		c.checked = !c.checked
	}
}

func (c *Checkbox) Draw(ctx *uikit.Context, dst *ebiten.Image) {
	c.Base.Draw(ctx, dst)

	r := c.Measure(false)

	// Checkbox box (left)
	content := common.Inset(r, ctx.Theme.PadX, ctx.Theme.PadY)
	boxSize := ctx.Theme.CheckSize
	if boxSize < 10 {
		boxSize = 10
	}

	boxY := r.Min.Y + (r.Dy()-boxSize)/2
	box := image.Rect(content.Min.X, boxY, content.Min.X+boxSize, boxY+boxSize)

	c.Base.DrawRoundedRect(dst, box, int(float64(boxSize)*0.25), ctx.Theme.Bg)
	c.Base.DrawRoundedBorder(dst, box, int(float64(boxSize)*0.25), ctx.Theme.BorderW, ctx.Theme.Border)

	if c.checked {
		// Draw a clean checkmark (two strokes), proportional.
		x1 := float32(box.Min.X) + float32(box.Dx())*0.22
		y1 := float32(box.Min.Y) + float32(box.Dy())*0.55
		x2 := float32(box.Min.X) + float32(box.Dx())*0.43
		y2 := float32(box.Min.Y) + float32(box.Dy())*0.73
		x3 := float32(box.Min.X) + float32(box.Dx())*0.78
		y3 := float32(box.Min.Y) + float32(box.Dy())*0.28

		w := float32(ctx.Theme.BorderW)
		if w < 2 {
			w = 2
		}
		vector.StrokeLine(dst, x1, y1, x2, y2, w, ctx.Theme.Focus, false)
		vector.StrokeLine(dst, x2, y2, x3, y3, w, ctx.Theme.Focus, false)
	}

	// Label
	met, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Min.Y + (r.Dy()-met.Height)/2 + met.Ascent
	tx := box.Max.X + ctx.Theme.SpaceS

	col := ctx.Theme.Text
	if !c.IsEnabled() {
		col = ctx.Theme.Disabled
	}

	ctx.Text.SetColor(col)
	ctx.Text.SetAlign(0) // Left
	ctx.Text.Draw(dst, c.label, tx, baselineY)
}
