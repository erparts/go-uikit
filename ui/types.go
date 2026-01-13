package ui

import "github.com/hajimehoshi/ebiten/v2"

// Layout is a Widget that owns children.
type Layout interface {
	Widget

	Children() []Widget
	SetChildren([]Widget)
	Add(...Widget)
	Clear()
}

// Hittable allows a widget to extend its clickable area beyond Base.ControlRect.
type Hittable interface {
	HitTest(ctx *Context, x, y int) bool
}

// OverlayWidget can draw an overlay above all other widgets (e.g. Select dropdown).
type OverlayWidget interface {
	OverlayActive() bool
	DrawOverlay(ctx *Context, dst *ebiten.Image)
}

// Themeable is implemented by widgets that need Theme for sizing (SetFrame) or drawing.
// Layouts should call SetTheme before calling SetFrame.
type Themeable interface {
	SetTheme(theme *Theme)
}

type PointerStatus struct {
	X, Y       int
	IsDown     bool
	IsJustDown bool
	IsJustUp   bool
	IsTouch    bool
	TouchID    ebiten.TouchID
}
