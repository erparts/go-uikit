package widget

import (
	"image"
	"math"

	"github.com/erparts/go-uikit"
	"github.com/erparts/go-uikit/common"
	"github.com/hajimehoshi/ebiten/v2"
)

// Select is a simple dropdown selector.
// The dropdown is rendered as an overlay (does NOT change layout of other widgets).
type Select struct {
	uikit.Base

	options []string
	index   int

	open bool

	// List scroll offset (in options, not pixels)
	scroll int

	// MaxVisible controls how many options are shown when open.
	MaxVisible int
}

func NewSelect(theme *uikit.Theme, options []string) *Select {
	cfg := uikit.NewWidgetBaseConfig(theme)

	return &Select{
		Base:       uikit.NewBase(cfg),
		options:    options,
		index:      0,
		MaxVisible: 5,
	}
}

func (s *Select) Focusable() bool { return true }

func (s *Select) OverlayActive() bool { return s.open }

func (s *Select) SetOptions(opts []string) {
	s.options = opts
	if s.index >= len(opts) {
		s.index = 0
	}

}

func (s *Select) Index() int { return s.index }
func (s *Select) Value() string {
	if s.index < 0 || s.index >= len(s.options) {
		return ""
	}
	return s.options[s.index]
}

func (s *Select) SetIndex(i int) {
	if len(s.options) == 0 {
		s.index = 0
		return
	}
	if i < 0 {
		i = 0
	}
	if i >= len(s.options) {
		i = len(s.options) - 1
	}
	s.index = i
}

func (s *Select) listRect(ctx *uikit.Context) image.Rectangle {
	ctrl := s.Measure(false)
	n := len(s.options)
	max := s.MaxVisible
	if max <= 0 {
		max = 5
	}
	if n > max {
		n = max
	}

	listY := ctrl.Max.Y + ctx.Theme.SpaceS
	return image.Rect(ctrl.Min.X, listY, ctrl.Max.X, listY+(n*ctx.Theme.ControlH))
}

func (s *Select) HitTest(ctx *uikit.Context, x, y int) bool {
	ctrl := s.Measure(false)
	if common.Contains(ctrl, x, y) {
		return true
	}

	if s.open {
		return common.Contains(s.listRect(ctx), x, y)
	}

	return false
}

func (s *Select) Update(ctx *uikit.Context) {
	r := s.Measure(false)
	if r.Dx() > 0 && r.Dy() == 0 {
		s.SetFrame(r.Min.X, r.Min.Y, r.Dx())
	}

	if !s.IsEnabled() {
		return
	}

	list := s.listRect(ctx)

	//func (c *Context) Pointer() (x, y int, down, justDown, justUp, isTouch bool) {

	ptr := ctx.Pointer()

	// Toggle open on click in control.
	if ptr.IsJustDown && common.Contains(r, ptr.X, ptr.Y) {
		s.open = !s.open
	}

	// When open, select option on click; close on click outside.
	if s.open && ptr.IsJustDown {
		if common.Contains(list, ptr.X, ptr.Y) {
			row := (ptr.Y - list.Min.Y) / ctx.Theme.ControlH
			idx := s.scroll + row
			if idx >= 0 && idx < len(s.options) {
				s.index = idx
			}
			s.open = false
		} else if !common.Contains(r, ptr.X, ptr.Y) {
			s.open = false
		}
	}

	// Wheel scroll when open (desktop)
	if s.open {
		_, wy := ebiten.Wheel()
		if wy != 0 {
			step := int(math.Copysign(1, wy))
			s.scroll -= step
			maxScroll := len(s.options) - s.MaxVisible
			if maxScroll < 0 {
				maxScroll = 0
			}
			if s.scroll < 0 {
				s.scroll = 0
			}
			if s.scroll > maxScroll {
				s.scroll = maxScroll
			}
		}
	}
}

func (s *Select) Draw(ctx *uikit.Context, dst *ebiten.Image) {
	s.Base.Draw(ctx, dst)

	r := s.Base.Measure(false)

	met, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Min.Y + (r.Dy()-met.Height)/2 + met.Ascent

	val := s.Value()
	if val == "" {
		val = "—"
	}
	ctx.Text.SetColor(ctx.Theme.Text)
	ctx.Text.SetAlign(0)
	ctx.Text.Draw(dst, val, r.Min.X+ctx.Theme.PadX, baselineY)

	// Chevron
	chev := "▾"
	cw := uikit.MeasureStringPx(ctx.Theme.Font, ctx.Theme.FontPx, chev)
	ctx.Text.Draw(dst, chev, r.Max.X-ctx.Theme.PadX-cw, baselineY)
}

func (s *Select) DrawOverlay(ctx *uikit.Context, dst *ebiten.Image) {
	if !s.open {
		return
	}

	list := s.listRect(ctx)
	s.Base.DrawRoundedRect(dst, list, ctx.Theme.Radius, ctx.Theme.Surface)
	s.Base.DrawRoundedBorder(dst, list, ctx.Theme.Radius, ctx.Theme.BorderW, ctx.Theme.Border)

	met, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)

	n := list.Dy() / ctx.Theme.ControlH
	for i := 0; i < n; i++ {
		idx := s.scroll + i
		if idx >= len(s.options) {
			break
		}

		y := list.Min.Y + i*ctx.Theme.ControlH
		row := image.Rect(list.Min.X, y, list.Max.X, y+ctx.Theme.ControlH)

		if idx == s.index {
			s.Base.DrawRoundedRect(dst, row, 0, ctx.Theme.SurfaceHover)
		}

		bY := row.Min.Y + (row.Dy()-met.Height)/2 + met.Ascent
		ctx.Text.SetColor(ctx.Theme.Text)
		ctx.Text.SetAlign(0)
		ctx.Text.Draw(dst, s.options[idx], row.Min.X+ctx.Theme.PadX, bY)
	}
}
