package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GridLayout places children in a fixed column grid. If height > 0 it becomes scrollable and clips via SubImage.
type GridLayout struct {
	base     Base
	children []Widget

	Columns int

	PadX int
	PadY int
	GapX int
	GapY int

	Scroll Scroller

	contentH int

	scratch *ebiten.Image
}

func NewGridLayout(theme *Theme) *GridLayout {
	l := &GridLayout{base: NewBase(&WidgetBaseConfig{})}
	l.Columns = 2
	//l.PadX = theme.SpaceM
	//l.PadY = theme.SpaceM
	l.GapX = theme.SpaceS
	l.GapY = theme.SpaceS
	l.Scroll = NewScroller()
	return l
}

func (l *GridLayout) Base() *Base     { return &l.base }
func (l *GridLayout) Focusable() bool { return false }

func (l *GridLayout) SetFrame(x, y, w int) { l.base.Rect = Rect{X: x, Y: y, W: w, H: l.base.Rect.H} }
func (l *GridLayout) Measure() Rect        { return l.base.Rect }
func (l *GridLayout) SetHeight(h int) {
	if h < 0 {
		h = 0
	}
	r := l.base.Rect
	r.H = h
	l.base.Rect = r
}

func (l *GridLayout) Children() []Widget      { return l.children }
func (l *GridLayout) SetChildren(ws []Widget) { l.children = ws }
func (l *GridLayout) Add(ws ...Widget)        { l.children = append(l.children, ws...) }
func (l *GridLayout) Clear()                  { l.children = nil }

func (l *GridLayout) Update(ctx *Context) {
	l.doLayout(ctx)

	if l.base.Rect.H > 0 {
		l.Scroll.Update(ctx, l.base.Rect, l.contentH)
		l.doLayout(ctx)
	}

	for _, ch := range l.children {
		if th, ok := any(ch).(Themeable); ok {
			th.SetTheme(ctx.Theme)
		}
		if !ch.Base().visible {
			continue
		}
		ch.Update(ctx)
	}
}

func (l *GridLayout) doLayout(ctx *Context) {
	vp := l.base.Rect
	cols := l.Columns
	if cols <= 0 {
		cols = 2
	}

	innerW := vp.W - l.PadX*2
	if innerW < 0 {
		innerW = 0
	}
	cellW := innerW
	if cols > 0 {
		cellW = (innerW - (cols-1)*l.GapX) / cols
		if cellW < 0 {
			cellW = 0
		}
	}

	x0 := vp.X + l.PadX
	y0 := vp.Y + l.PadY
	x := x0
	y := y0
	if vp.H > 0 {
		y -= l.Scroll.ScrollY
	}

	contentH := l.PadY * 2
	rowMaxH := 0
	col := 0

	for i, ch := range l.children {
		if th, ok := any(ch).(Themeable); ok {
			th.SetTheme(ctx.Theme)
		}
		if !ch.Base().visible {
			continue
		}
		ch.SetFrame(x, y, cellW)
		r := ch.Measure()
		if r.H > rowMaxH {
			rowMaxH = r.H
		}

		col++
		last := i == len(l.children)-1
		if col >= cols || last {
			contentH += rowMaxH
			if !last {
				contentH += l.GapY
			}
			y += rowMaxH + l.GapY
			x = x0
			col = 0
			rowMaxH = 0
		} else {
			x += cellW + l.GapX
		}
	}

	if vp.H > 0 && contentH < vp.H {
		contentH = vp.H
	}
	l.contentH = contentH
}

func (l *GridLayout) Draw(ctx *Context, dst *ebiten.Image) {
	if !l.base.visible {
		return
	}
	vp := l.base.Rect
	if vp.H <= 0 {
		for _, ch := range l.children {
			if th, ok := any(ch).(Themeable); ok {
				th.SetTheme(ctx.Theme)
			}
			if !ch.Base().visible {
				continue
			}
			ch.Draw(ctx, dst)
		}
		return
	}

	// Scrollable: render to a full-screen scratch (no coordinate shifting),
	// then copy only the viewport region back to dst using SubImage.
	sw, sh := dst.Bounds().Dx(), dst.Bounds().Dy()
	if l.scratch == nil || l.scratch.Bounds().Dx() != sw || l.scratch.Bounds().Dy() != sh {
		l.scratch = ebiten.NewImage(sw, sh)
	}
	l.scratch.Clear()

	for _, ch := range l.children {
		if th, ok := any(ch).(Themeable); ok {
			th.SetTheme(ctx.Theme)
		}
		if !ch.Base().visible {
			continue
		}
		ch.Draw(ctx, l.scratch)
	}

	part := l.scratch.SubImage(vp.ImageRect()).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(vp.X), float64(vp.Y))
	dst.DrawImage(part, op)

	sub := dst.SubImage(vp.ImageRect()).(*ebiten.Image)
	l.Scroll.DrawBar(sub, ctx.Theme, vp.W, vp.H, l.contentH)
}

func (l *GridLayout) DrawOverlay(ctx *Context, dst *ebiten.Image) {
	if !l.base.visible {
		return
	}
	for _, ch := range l.children {
		if th, ok := any(ch).(Themeable); ok {
			th.SetTheme(ctx.Theme)
		}
		if ow, ok := any(ch).(OverlayWidget); ok && ow.OverlayActive() {
			ow.DrawOverlay(ctx, dst)
		}
		if ll, ok := any(ch).(interface{ DrawOverlay(*Context, *ebiten.Image) }); ok {
			ll.DrawOverlay(ctx, dst)
		}
	}
}
