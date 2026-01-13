package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// StackLayout places children vertically. If height > 0 it becomes scrollable and clips via SubImage.
type StackLayout struct {
	base     Base
	children []Widget

	PadX int
	PadY int
	Gap  int

	Scroll Scroller

	contentH int

	scratch *ebiten.Image
}

func NewStackLayout(theme *Theme) *StackLayout {
	l := &StackLayout{base: NewBase(&WidgetBaseConfig{})}
	l.base.SetEnabled(true)
	//l.PadX = theme.SpaceM
	//l.PadY = theme.SpaceM
	l.Gap = theme.SpaceS
	l.Scroll = NewScroller()
	return l
}

func (l *StackLayout) Base() *Base     { return &l.base }
func (l *StackLayout) Focusable() bool { return false }

func (l *StackLayout) SetFrame(x, y, w int) {
	// Keep current height (0 => unlimited). Only external code sets height via SetHeight/SetRect.
	l.base.Rect = Rect{X: x, Y: y, W: w, H: l.base.Rect.H}
}

func (l *StackLayout) Measure() Rect { return l.base.Rect }

// SetHeight sets the viewport height. Use 0 for unlimited (no scroll, no clipping).
func (l *StackLayout) SetHeight(h int) {
	if h < 0 {
		h = 0
	}
	r := l.base.Rect
	r.H = h
	l.base.Rect = r
}

// Children management
func (l *StackLayout) Children() []Widget      { return l.children }
func (l *StackLayout) SetChildren(ws []Widget) { l.children = ws }
func (l *StackLayout) Add(ws ...Widget)        { l.children = append(l.children, ws...) }
func (l *StackLayout) Clear()                  { l.children = nil }

func (l *StackLayout) Update(ctx *Context) {
	// Layout pass (compute frames and content height)
	l.doLayout(ctx)

	// Scroll input only when height is limited
	if l.base.Rect.H > 0 {
		l.Scroll.Update(ctx, l.base.Rect, l.contentH)
		// Re-layout with updated scroll offset
		l.doLayout(ctx)
	}

	// Forward update
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

func (l *StackLayout) doLayout(ctx *Context) {
	vp := l.base.Rect
	x0 := vp.X + l.PadX
	y0 := vp.Y + l.PadY
	w0 := vp.W - l.PadX*2
	if w0 < 0 {
		w0 = 0
	}

	y := y0
	if vp.H > 0 {
		y -= l.Scroll.ScrollY
	}

	contentH := l.PadY * 2
	for i, ch := range l.children {
		if th, ok := any(ch).(Themeable); ok {
			th.SetTheme(ctx.Theme)
		}
		if !ch.Base().visible {
			continue
		}
		ch.SetFrame(x0, y, w0)
		r := ch.Measure()
		contentH += r.H
		if i != len(l.children)-1 {
			contentH += l.Gap
		}
		y += r.H + l.Gap
	}
	// At least viewport height so scrollbar math is stable
	if vp.H > 0 && contentH < vp.H {
		contentH = vp.H
	}
	l.contentH = contentH
}

func (l *StackLayout) Draw(ctx *Context, dst *ebiten.Image) {
	if !l.base.visible {
		return
	}

	vp := l.base.Rect
	if vp.H <= 0 {
		// Unlimited: draw directly
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

	// Scrollbar inside viewport (draw on clipped dst region)
	sub := dst.SubImage(vp.ImageRect()).(*ebiten.Image)
	l.Scroll.DrawBar(sub, ctx.Theme, vp.W, vp.H, l.contentH)
}

func (l *StackLayout) DrawOverlay(ctx *Context, dst *ebiten.Image) {
	if !l.base.visible {
		return
	}
	// Overlay should escape clipping -> draw on dst (not on subimage)
	for _, ch := range l.children {
		if th, ok := any(ch).(Themeable); ok {
			th.SetTheme(ctx.Theme)
		}
		if ow, ok := any(ch).(OverlayWidget); ok && ow.OverlayActive() {
			ow.DrawOverlay(ctx, dst)
		}
		// Nested layouts will propagate overlays naturally.
		if ll, ok := any(ch).(interface{ DrawOverlay(*Context, *ebiten.Image) }); ok {
			ll.DrawOverlay(ctx, dst)
		}
	}
}
