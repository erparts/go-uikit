package widget

import (
	"image"

	"github.com/erparts/go-uikit"
	"github.com/hajimehoshi/ebiten/v2"
)

type Label struct {
	base uikit.Base
	text string
}

func NewLabel(theme *uikit.Theme, text string) *Label {
	cfg := uikit.NewWidgetBaseConfig(theme)
	cfg.DrawSurface = false
	cfg.DrawBorder = false

	base := uikit.NewBase(cfg)

	return &Label{
		base: base,
		text: text,
	}
}

func (l *Label) Base() *uikit.Base { return &l.base }
func (l *Label) Focusable() bool   { return false }
func (l *Label) SetText(s string)  { l.text = s }
func (l *Label) SetEnabled(v bool) { l.base.SetEnabled(v) }
func (l *Label) SetVisible(v bool) { l.base.SetVisible(v) }
func (l *Label) SetFrame(x, y, w int) {
	l.base.SetFrame(x, y, w)
}

func (l *Label) Measure() image.Rectangle { return l.base.Rect }

func (l *Label) Update(ctx *uikit.Context) {
	if l.base.Rect.Dy() == 0 {
		l.base.SetFrame(l.base.Rect.Min.X, l.base.Rect.Min.Y, l.base.Rect.Dx())
	}
}

func (l *Label) Draw(ctx *uikit.Context, dst *ebiten.Image) {
	r := l.base.Draw(ctx, dst)

	met, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Min.Y + (r.Dy()-met.Height)/2 + met.Ascent

	ctx.Text.SetColor(ctx.Theme.MutedText)
	ctx.Text.SetAlign(0) // Left
	ctx.Text.Draw(dst, l.text, r.Min.X, baselineY)

}
