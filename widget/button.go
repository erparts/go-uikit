package widget

import (
	"github.com/erparts/go-uikit"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

type Button struct {
	uikit.Base
	label   string
	OnClick func()
}

func NewButton(theme *uikit.Theme, label string) *Button {
	cfg := uikit.NewWidgetBaseConfig(theme)

	return &Button{
		Base:  uikit.NewBase(cfg),
		label: label,
	}
}

func (b *Button) Focusable() bool {
	return true
}

func (b *Button) SetLabel(s string) { b.label = s }

func (b *Button) HandleEvent(ctx *uikit.Context, e uikit.Event) {
	if !b.IsEnabled() {
		return
	}

	if e.Type == uikit.EventClick {
		if b.OnClick != nil {
			b.OnClick()
		}
	}
}

func (b *Button) Update(ctx *uikit.Context) {
	if !b.IsEnabled() {
		return
	}

	// Keyboard click when focused
	if b.IsFocused() && (inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace)) {
		if b.OnClick != nil {
			b.OnClick()
		}
	}
}

func (b *Button) Draw(ctx *uikit.Context, dst *ebiten.Image) {
	r := b.Base.Draw(ctx, dst)

	// Centered label
	met, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	textW := uikit.MeasureStringPx(ctx.Theme.Font, ctx.Theme.FontPx, b.label)

	tx := r.Min.X + (r.Dx()-textW)/2
	baselineY := r.Min.Y + (r.Dy()-met.Height)/2 + met.Ascent

	col := ctx.Theme.Text
	if !b.Base.IsEnabled() {
		col = ctx.Theme.Disabled
	}

	ctx.Text.SetColor(col)
	ctx.Text.SetAlign(etxt.Left)
	ctx.Text.Draw(dst, b.label, tx, baselineY)
}
