package uikit

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Base() *Base
	Focusable() bool
	// SetFrame sets the widget position and width. Height is derived from the theme and widget state.
	SetFrame(x, y, w int)
	// Measure returns the current widget rectangle (including any extra height such as validation errors).
	Measure() image.Rectangle
	Update(ctx *Context)
	Draw(ctx *Context, dst *ebiten.Image)
}

// OverlayWidget can draw an overlay above all other widgets (e.g. Select dropdown).
type OverlayWidget interface {
	OverlayActive() bool
	DrawOverlay(ctx *Context, dst *ebiten.Image)
}

// TextWidget is implemented by widgets that want to control the platform IME (e.g., TextInput, TextArea).
type TextWidget interface {
	Widget
	WantsIME() bool
}

// Hittable allows a widget to extend its clickable area beyond Base.ControlRect.
type Hittable interface {
	HitTest(ctx *Context, x, y int) bool
}

// Layout is a Widget that owns children.
type Layout interface {
	Widget

	Children() []Widget
	SetChildren([]Widget)
	Add(...Widget)
	Clear()
}

type PointerStatus struct {
	X, Y       int
	IsDown     bool
	IsJustDown bool
	IsJustUp   bool
	IsTouch    bool
	TouchID    ebiten.TouchID
}
