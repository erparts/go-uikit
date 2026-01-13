package widget

import (
	"github.com/erparts/go-uikit/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Label struct {
	base  ui.Base
	text  string
	theme *ui.Theme
}

func NewLabel(text string) *Label {
	cfg := ui.NewWidgetBaseConfig()
	cfg.DrawSurface = false
	cfg.DrawBorder = false

	base := ui.NewBase(cfg)

	return &Label{
		base: base,
		text: text,
	}
}

func (l *Label) Base() *ui.Base    { return &l.base }
func (l *Label) Focusable() bool   { return false }
func (l *Label) SetText(s string)  { l.text = s }
func (l *Label) SetEnabled(v bool) { l.base.SetEnabled(v) }
func (l *Label) SetVisible(v bool) { l.base.SetVisible(v) }
func (l *Label) SetFrame(x, y, w int) {
	// Height is derived from theme; if theme not known yet, keep H=0 and
	// it will be fixed on first Update/Draw.
	if l.theme != nil {
		l.base.SetFrame(l.theme, x, y, w)
		return
	}

	l.base.Rect = ui.Rect{X: x, Y: y, W: w, H: 0}
}

func (l *Label) Measure() ui.Rect { return l.base.Rect }

func (l *Label) Update(ctx *ui.Context) {
	l.theme = ctx.Theme
	if l.base.Rect.H == 0 {
		l.base.SetFrame(ctx.Theme, l.base.Rect.X, l.base.Rect.Y, l.base.Rect.W)
	}
}

func (l *Label) Draw(ctx *ui.Context, dst *ebiten.Image) {
	r := l.base.Draw(ctx, dst)

	met, _ := ui.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Y + (r.H-met.Height)/2 + met.Ascent

	ctx.Text.SetColor(ctx.Theme.MutedText)
	ctx.Text.SetAlign(0) // Left
	ctx.Text.Draw(dst, l.text, r.X, baselineY)

}

// SetTheme allows layouts to provide Theme before SetFrame is called.
func (l *Label) SetTheme(theme *ui.Theme) { l.theme = theme }
