package uikit

import (
	"image"

	"github.com/erparts/go-uikit/common"
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

func (l *GridLayout) SetFrame(x, y, w int) {
	l.base.Rect = image.Rect(x, y, x+w, y+l.base.Rect.Dy())
}

func (l *GridLayout) Measure() image.Rectangle {
	return l.base.Rect
}

func (l *GridLayout) SetHeight(h int) {
	if h < 0 {
		h = 0
	}

	l.base.Rect = common.ChangeRectangleHeight(l.base.Rect, h)
}

func (l *GridLayout) Children() []Widget      { return l.children }
func (l *GridLayout) SetChildren(ws []Widget) { l.children = ws }
func (l *GridLayout) Add(ws ...Widget)        { l.children = append(l.children, ws...) }
func (l *GridLayout) Clear()                  { l.children = nil }

func (l *GridLayout) Update(ctx *Context) {
	l.doLayout(ctx)

	if l.base.Rect.Dy() > 0 {
		l.Scroll.Update(ctx, l.base.Rect, l.contentH)
		l.doLayout(ctx)
	}

	for _, ch := range l.children {
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

	innerW := vp.Dx() - l.PadX*2
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

	x0 := vp.Min.X + l.PadX
	y0 := vp.Min.Y + l.PadY
	x := x0
	y := y0
	if vp.Dy() > 0 {
		y -= l.Scroll.ScrollY
	}

	contentH := l.PadY * 2
	rowMaxH := 0
	col := 0

	for i, ch := range l.children {
		if !ch.Base().visible {
			continue
		}
		ch.SetFrame(x, y, cellW)
		r := ch.Measure()
		if r.Dy() > rowMaxH {
			rowMaxH = r.Dy()
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

	if vp.Dy() > 0 && contentH < vp.Dy() {
		contentH = vp.Dy()
	}

	l.contentH = contentH
}

func (l *GridLayout) Draw(ctx *Context, dst *ebiten.Image) {
	if !l.base.visible {
		return
	}
	vp := l.base.Rect
	if vp.Dy() <= 0 {
		for _, ch := range l.children {
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
		if !ch.Base().visible {
			continue
		}
		ch.Draw(ctx, l.scratch)
	}

	part := l.scratch.SubImage(vp).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(vp.Min.X), float64(vp.Min.Y))
	dst.DrawImage(part, op)

	sub := dst.SubImage(vp).(*ebiten.Image)
	l.Scroll.DrawBar(sub, ctx.Theme, vp.Dx(), vp.Dy(), l.contentH)
}

func (l *GridLayout) DrawOverlay(ctx *Context, dst *ebiten.Image) {
	if !l.base.visible {
		return
	}
	for _, ch := range l.children {
		if ow, ok := any(ch).(OverlayWidget); ok && ow.OverlayActive() {
			ow.DrawOverlay(ctx, dst)
		}
		if ll, ok := any(ch).(interface{ DrawOverlay(*Context, *ebiten.Image) }); ok {
			ll.DrawOverlay(ctx, dst)
		}
	}
}
