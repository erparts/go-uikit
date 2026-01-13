package widget

import (
	"math"

	"github.com/erparts/go-uikit/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

// Select is a simple dropdown selector.
// The dropdown is rendered as an overlay (does NOT change layout of other widgets).
type Select struct {
	base  ui.Base
	theme *ui.Theme

	options []string
	index   int

	open bool

	// List scroll offset (in options, not pixels)
	scroll int

	// MaxVisible controls how many options are shown when open.
	MaxVisible int
}

func NewSelect(options []string) *Select {
	cfg := ui.NewWidgetBaseConfig()

	return &Select{
		base:       ui.NewBase(cfg),
		options:    options,
		index:      0,
		MaxVisible: 5,
	}
}

func (s *Select) Base() *ui.Base  { return &s.base }
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

func (s *Select) SetFrame(x, y, w int) {
	if s.theme != nil {
		s.base.SetFrame(s.theme, x, y, w)
		return
	}
	s.base.Rect = ui.Rect{X: x, Y: y, W: w, H: 0}
}

func (s *Select) Measure() ui.Rect { return s.base.Rect }

func (s *Select) listRect(ctx *ui.Context) ui.Rect {
	ctrl := s.base.ControlRect(ctx.Theme)
	n := len(s.options)
	max := s.MaxVisible
	if max <= 0 {
		max = 5
	}
	if n > max {
		n = max
	}
	listY := ctrl.Bottom() + ctx.Theme.SpaceS
	return ui.Rect{X: ctrl.X, Y: listY, W: ctrl.W, H: n * ctx.Theme.ControlH}
}

func (s *Select) HitTest(ctx *ui.Context, x, y int) bool {
	ctrl := s.base.ControlRect(ctx.Theme)
	if ctrl.Contains(x, y) {
		return true
	}
	if s.open {
		return s.listRect(ctx).Contains(x, y)
	}
	return false
}

func (s *Select) Update(ctx *ui.Context) {
	s.theme = ctx.Theme
	if s.base.Rect.H == 0 {
		s.base.SetFrame(ctx.Theme, s.base.Rect.X, s.base.Rect.Y, s.base.Rect.W)
	}
	if !s.base.IsEnabled() {
		return
	}

	ctrl := s.base.ControlRect(ctx.Theme)
	list := s.listRect(ctx)

	//func (c *Context) Pointer() (x, y int, down, justDown, justUp, isTouch bool) {

	ptr := ctx.Pointer()

	// Toggle open on click in control.
	if ptr.IsJustDown && ctrl.Contains(ptr.X, ptr.Y) {
		s.open = !s.open
	}

	// When open, select option on click; close on click outside.
	if s.open && ptr.IsJustDown {
		if list.Contains(ptr.X, ptr.Y) {
			row := (ptr.X - list.Y) / ctx.Theme.ControlH
			idx := s.scroll + row
			if idx >= 0 && idx < len(s.options) {
				s.index = idx
			}
			s.open = false
		} else if !ctrl.Contains(ptr.X, ptr.Y) {
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

func (s *Select) Draw(ctx *ui.Context, dst *ebiten.Image) {
	s.base.Draw(ctx, dst)

	r := s.base.ControlRect(ctx.Theme)

	met, _ := ui.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Y + (r.H-met.Height)/2 + met.Ascent

	val := s.Value()
	if val == "" {
		val = "—"
	}
	ctx.Text.SetColor(ctx.Theme.Text)
	ctx.Text.SetAlign(0)
	ctx.Text.Draw(dst, val, r.X+ctx.Theme.PadX, baselineY)

	// Chevron
	chev := "▾"
	cw := ui.MeasureStringPx(ctx.Theme.Font, ctx.Theme.FontPx, chev)
	ctx.Text.Draw(dst, chev, r.Right()-ctx.Theme.PadX-cw, baselineY)
}

func (s *Select) DrawOverlay(ctx *ui.Context, dst *ebiten.Image) {
	if !s.open {
		return
	}

	list := s.listRect(ctx)
	s.base.DrawRoundedRect(dst, list, ctx.Theme.Radius, ctx.Theme.Surface)
	s.base.DrawRoundedBorder(dst, list, ctx.Theme.Radius, ctx.Theme.BorderW, ctx.Theme.Border)

	met, _ := ui.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)

	n := list.H / ctx.Theme.ControlH
	for i := 0; i < n; i++ {
		idx := s.scroll + i
		if idx >= len(s.options) {
			break
		}
		row := ui.Rect{X: list.X, Y: list.Y + i*ctx.Theme.ControlH, W: list.W, H: ctx.Theme.ControlH}

		if idx == s.index {
			s.base.DrawRoundedRect(dst, row, 0, ctx.Theme.SurfaceHover)
		}

		bY := row.Y + (row.H-met.Height)/2 + met.Ascent
		ctx.Text.SetColor(ctx.Theme.Text)
		ctx.Text.SetAlign(0)
		ctx.Text.Draw(dst, s.options[idx], row.X+ctx.Theme.PadX, bY)
	}
}

// SetTheme allows layouts to provide Theme before SetFrame is called.
func (s *Select) SetTheme(theme *ui.Theme) { s.theme = theme }
