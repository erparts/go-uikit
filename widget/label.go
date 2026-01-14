package widget

import (
	"github.com/erparts/go-uikit"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type Label struct {
	uikit.Base
	text string
}

func NewLabel(theme *uikit.Theme, text string) *Label {
	cfg := uikit.NewWidgetBaseConfig(theme)
	cfg.DrawSurface = false
	cfg.DrawBorder = false

	base := uikit.NewBase(cfg)

	return &Label{
		Base: base,
		text: text,
	}
}

func (l *Label) Focusable() bool  { return false }
func (l *Label) SetText(s string) { l.text = s }

func (l *Label) Update(ctx *uikit.Context) {
	r := l.Measure(false)
	if r.Dy() == 0 {
		l.SetFrame(r.Min.X, r.Min.Y, r.Dx())
	}
}

func (l *Label) Draw(ctx *uikit.Context, dst *ebiten.Image) {
	r := l.Base.Draw(ctx, dst)

	met, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Min.Y + (r.Dy()-met.Height)/2 + met.Ascent

	ctx.Text.SetColor(ctx.Theme.MutedText)
	ctx.Text.SetAlign(etxt.Left)
	ctx.Text.Draw(dst, l.text, r.Min.X, baselineY)
}
