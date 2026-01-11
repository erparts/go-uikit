package ui

// EventType represents high-level UI events dispatched by Context.
type EventType int

const (
	EventNone EventType = iota
	EventFocusGained
	EventFocusLost
	EventPointerDown
	EventPointerUp
	EventClick
	EventKeyDown
	EventKeyUp
)

// Event is a UI event routed to a widget.
type Event struct {
	Type EventType

	// Pointer position in pixels
	X int
	Y int

	// For key events
	Key int // ebiten.Key as int (kept int to avoid importing ebiten in ui core)
}

// EventHandler can be implemented by widgets that want explicit event delivery.
type EventHandler interface {
	HandleEvent(ctx *Context, e Event)
}
