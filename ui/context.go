package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

// Context holds shared state for all widgets.
type Context struct {
	Theme *Theme
	Text  *etxt.Renderer
	IME   IMEBridge // optional (nil on desktop)

	widgets []Widget
	focus   int // -1 means none

	// Mouse state (screen pixels)
	mouseX, mouseY int
	mouseDown      bool
	mouseJustDown  bool
	mouseJustUp    bool
}

func NewContext(theme *Theme, renderer *etxt.Renderer, ime IMEBridge) *Context {
	// Ensure renderer style is consistent with the theme.
	renderer.SetFont(theme.Font)
	renderer.SetSize(float64(theme.FontPx))
	return &Context{
		Theme: theme,
		Text:  renderer,
		IME:   ime,
		focus: -1,
	}
}

// SetIMEBridge sets/updates the IME bridge at runtime.
// It also synchronizes the IME visibility with the currently focused widget.
func (c *Context) SetIMEBridge(b IMEBridge) {
	c.IME = b
	c.updateIMEForce(c.Focused())
}

func (c *Context) Add(w Widget) { c.widgets = append(c.widgets, w) }
func (c *Context) Widgets() []Widget {
	return c.widgets
}

func (c *Context) Focused() Widget {
	if c.focus < 0 || c.focus >= len(c.widgets) {
		return nil
	}
	return c.widgets[c.focus]
}

func (c *Context) dispatch(w Widget, e Event) {
	if w == nil {
		return
	}
	if h, ok := any(w).(EventHandler); ok {
		h.HandleEvent(c, e)
	}
}

func (c *Context) SetFocus(w Widget) {
	old := c.Focused()

	// Resolve new focus index (or -1).
	newIdx := -1
	if w != nil {
		for i, ww := range c.widgets {
			if ww == w {
				newIdx = i
				break
			}
		}
	}

	// Emit focus events if changed
	if old != nil && (newIdx != c.focus) {
		c.dispatch(old, Event{Type: EventFocusLost})
	}
	c.focus = newIdx
	newW := c.Focused()
	if newW != nil && newW != old {
		c.dispatch(newW, Event{Type: EventFocusGained})
	}

	// IME show/hide based on focused widget.
	c.updateIME(old, newW)
}

func (c *Context) updateIME(oldW, newW Widget) {
	if c.IME == nil {
		return
	}

	oldWants := false
	if oldW != nil {
		if wi, ok := any(oldW).(WantsIME); ok && wi.WantsIME() {
			oldWants = true
		}
	}
	newWants := false
	if newW != nil {
		if wi, ok := any(newW).(WantsIME); ok && wi.WantsIME() {
			newWants = true
		}
	}

	// Only issue calls on state transitions.
	if oldWants && !newWants {
		c.IME.Hide()
	}
	if !oldWants && newWants {
		c.IME.Show()
	}
}

func (c *Context) updateIMEForce(focused Widget) {
	if c.IME == nil {
		return
	}
	wants := false
	if focused != nil {
		if wi, ok := any(focused).(WantsIME); ok && wi.WantsIME() {
			wants = true
		}
	}
	if wants {
		c.IME.Show()
	} else {
		c.IME.Hide()
	}
}

func (c *Context) focusNext() {
	if len(c.widgets) == 0 {
		c.SetFocus(nil)
		return
	}
	start := c.focus
	for i := 0; i < len(c.widgets); i++ {
		idx := (start + 1 + i) % len(c.widgets)
		if c.widgets[idx].Base().Visible && c.widgets[idx].Base().Enabled && c.widgets[idx].Focusable() {
			c.SetFocus(c.widgets[idx])
			return
		}
	}
}

func (c *Context) focusPrev() {
	if len(c.widgets) == 0 {
		c.SetFocus(nil)
		return
	}
	start := c.focus
	for i := 0; i < len(c.widgets); i++ {
		idx := start - 1 - i
		for idx < 0 {
			idx += len(c.widgets)
		}
		if c.widgets[idx].Base().Visible && c.widgets[idx].Base().Enabled && c.widgets[idx].Focusable() {
			c.SetFocus(c.widgets[idx])
			return
		}
	}
}

func (c *Context) Update() {
	// Snapshot mouse state
	c.mouseX, c.mouseY = ebiten.CursorPosition()
	c.mouseDown = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	c.mouseJustDown = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	c.mouseJustUp = inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)

	// Keyboard focus traversal
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			c.focusPrev()
		} else {
			c.focusNext()
		}
	}

	// Widget updates + basic event routing
	for _, w := range c.widgets {
		b := w.Base()
		if !b.Visible {
			continue
		}

		inside := b.Rect.Contains(c.mouseX, c.mouseY)
		b.hovered = inside && b.Enabled

		// Pointer down
		if c.mouseJustDown && inside && b.Enabled {
			b.pressed = true
			c.dispatch(w, Event{Type: EventPointerDown, X: c.mouseX, Y: c.mouseY})
			if w.Focusable() {
				c.SetFocus(w)
			}
		}

		// Focus flag
		b.focused = (c.Focused() == w) && b.Enabled && w.Focusable()

		// Let widgets update (legacy path)
		w.Update(c)

		// Pointer up + click
		if c.mouseJustUp {
			if b.pressed {
				c.dispatch(w, Event{Type: EventPointerUp, X: c.mouseX, Y: c.mouseY})
				if inside && b.Enabled {
					c.dispatch(w, Event{Type: EventClick, X: c.mouseX, Y: c.mouseY})
				}
			}
			b.pressed = false
		}

		// Recompute focus flag
		b.focused = (c.Focused() == w) && b.Enabled && w.Focusable()
	}
}

func (c *Context) Draw(dst *ebiten.Image) {
	for _, w := range c.widgets {
		if !w.Base().Visible {
			continue
		}
		w.Draw(c, dst)
	}
}
