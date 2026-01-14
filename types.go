package uikit

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	SetFrame(x, y, w int)

	IsHovered() bool
	SetHovered(bool)
	IsPressed() bool
	SetPressed(bool)
	IsFocused() bool
	SetFocused(bool)
	IsEnabled() bool
	SetEnabled(bool)
	IsVisible() bool
	SetVisible(bool)

	Focusable() bool

	Measure(bool) image.Rectangle
	Update(ctx *Context)
	Draw(ctx *Context, dst *ebiten.Image)
}

// OverlayWidget can draw an overlay above all other widgets (e.g. Select dropdown).
type OverlayWidget interface {
	OverlayActive() bool
	DrawOverlay(ctx *Context, dst *ebiten.Image)
}

type ValidableWidget interface {
	IsValidable() bool
	IsInvalid() (bool, string)
	SetInvalid(err string)
	ClearInvalid()
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
