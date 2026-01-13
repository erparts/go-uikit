package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

// Context holds shared state for all widgets.
type Context struct {
	root  Layout
	Theme *Theme
	Text  *etxt.Renderer
	IME   IMEBridge // optional (nil on desktop)
	Scale Scale     // kept for apps, but input + layout are in logical pixels

	widgets []Widget
	focus   int // -1 means none

	ptr         *PointerStatus
	hasTouch    bool
	prevTouches map[ebiten.TouchID]struct{}
}

func NewContext(theme *Theme, renderer *etxt.Renderer, ime IMEBridge) *Context {
	// Ensure renderer style is consistent with the theme.
	renderer.SetFont(theme.Font)
	renderer.SetSize(float64(theme.FontPx))

	root := NewStackLayout(theme)

	return &Context{
		Theme:       theme,
		Text:        renderer,
		IME:         ime,
		Scale:       Scale{Device: 1, UI: 1},
		focus:       -1,
		prevTouches: map[ebiten.TouchID]struct{}{},
		root:        root,
		widgets:     []Widget{root},
		ptr:         &PointerStatus{},
	}
}

// Root returns the root widget (typically a Layout).
func (c *Context) Root() Layout { return c.root }

// SetRoot replaces the root widget.
func (c *Context) SetRoot(l Layout) {
	c.root = l
}

// SetScale stores a scale value for the app. The default Context input space is logical pixels.
func (c *Context) SetScale(s Scale) { c.Scale = s }

// SetIMEBridge sets/updates the IME bridge at runtime.
// It also synchronizes the IME visibility with the currently focused widget.
func (c *Context) SetIMEBridge(b IMEBridge) {
	c.IME = b
	c.updateIMEForce(c.Focused())
}

func (c *Context) Add(w Widget) {
	c.root.Add(w)
}

func (c *Context) Focused() Widget {

	if c.focus < 0 || c.focus >= len(c.widgets) {
		return nil
	}
	return c.widgets[c.focus]
}

// Pointer returns the current pointer state in logical pixels.
// On desktop this is the mouse; on mobile this is the active touch.

func (c *Context) rebuildWidgets() {
	c.widgets = c.widgets[:0]
	var walk func(w Widget)
	walk = func(w Widget) {
		if w == nil {
			return
		}
		c.widgets = append(c.widgets, w)
		if l, ok := any(w).(Layout); ok {
			for _, ch := range l.Children() {
				walk(ch)
			}
			return
		}
		// Support nested layouts even if not typed as Layout (defensive)
		if hw, ok := any(w).(interface{ Children() []Widget }); ok {
			for _, ch := range hw.Children() {
				walk(ch)
			}
		}
	}
	walk(c.root)
}

func (c *Context) Pointer() PointerStatus {
	return *c.ptr
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
		if c.widgets[idx].Base().visible && c.widgets[idx].Base().IsEnabled() && c.widgets[idx].Focusable() {
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
		if c.widgets[idx].Base().visible && c.widgets[idx].Base().IsEnabled() && c.widgets[idx].Focusable() {
			c.SetFocus(c.widgets[idx])
			return
		}
	}
}

func (c *Context) readPointerSnapshot() {
	c.ptr.IsJustDown = false
	c.ptr.IsJustUp = false
	c.ptr.IsTouch = false

	// Touch tracking (prefer this on mobile; CursorPosition is always (0,0) there).
	touches := ebiten.TouchIDs()
	curr := map[ebiten.TouchID]struct{}{}
	for _, id := range touches {
		curr[id] = struct{}{}
	}

	// Determine transitions
	var justPressed []ebiten.TouchID
	var justReleased []ebiten.TouchID
	for id := range curr {
		if _, ok := c.prevTouches[id]; !ok {
			justPressed = append(justPressed, id)
		}
	}
	for id := range c.prevTouches {
		if _, ok := curr[id]; !ok {
			justReleased = append(justReleased, id)
		}
	}
	c.prevTouches = curr

	// Acquire an active touch when pressed
	if !c.hasTouch && len(justPressed) > 0 {
		c.ptr.TouchID = justPressed[0]
		c.hasTouch = true
		c.ptr.IsJustDown = true
	}

	if c.hasTouch {
		if _, ok := curr[c.ptr.TouchID]; ok {
			c.ptr.IsDown = true
			c.ptr.IsTouch = true
			c.ptr.X, c.ptr.Y = ebiten.TouchPosition(c.ptr.TouchID)
		} else {
			c.ptr.IsDown = false
			c.ptr.IsTouch = true
			c.ptr.IsJustUp = true
			c.hasTouch = false
		}

		return
	}

	c.ptr.X, c.ptr.Y = ebiten.CursorPosition()
	c.ptr.IsDown = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	c.ptr.IsJustDown = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	c.ptr.IsJustUp = inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

func (c *Context) widgetHit(w Widget, x, y int) bool {
	if h, ok := any(w).(Hittable); ok {
		return h.HitTest(c, x, y)
	}
	// Default: only the control area is interactive (not the error message line).
	return w.Base().ControlRect(c.Theme).Contains(x, y)
}

func (c *Context) topmostAt(x, y int) Widget {
	for i := len(c.widgets) - 1; i >= 0; i-- {
		w := c.widgets[i]
		b := w.Base()
		if !b.visible || !b.IsEnabled() {
			continue
		}
		if c.widgetHit(w, x, y) {
			return w
		}
	}
	return nil
}

func (c *Context) Update() {
	// Snapshot input once per frame
	c.readPointerSnapshot()
	c.root.Update(c)

	// Rebuild flat list for hit-testing, focus traversal, and event dispatch.
	c.rebuildWidgets()

	// Keyboard focus traversal (desktop)
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			c.focusPrev()
		} else {
			c.focusNext()
		}
	}

	// IME behavior: any tap/click outside a text input closes the IME.
	// We decide focus target on pointer down using a "topmost hit" strategy.
	if c.ptr.IsJustDown {
		w := c.topmostAt(c.ptr.X, c.ptr.Y)
		if w != nil {
			if tw, ok := w.(TextWidget); ok && tw.WantsIME() {
				c.SetFocus(w)
			} else {
				// If clicked non-text widget, clear focus (closes IME).
				c.SetFocus(nil)
			}
		} else {
			c.SetFocus(nil)
		}
	}

	// Pointer down chooses a single target (topmost), prevents underlying widgets from also receiving presses.
	var target Widget
	if c.ptr.IsJustDown {
		target = c.topmostAt(c.ptr.X, c.ptr.Y)
	}

	// Hover only for mouse: topmost under pointer.
	var hoverTarget Widget
	if !c.ptr.IsTouch {
		hoverTarget = c.topmostAt(c.ptr.X, c.ptr.Y)
	}

	for _, w := range c.widgets {
		b := w.Base()
		if !b.visible {
			continue
		}

		b.hovered = (hoverTarget == w)

		// Pointer down routed to the chosen target.
		if c.ptr.IsJustDown && target == w && b.IsEnabled() {
			b.pressed = true
			c.dispatch(w, Event{Type: EventPointerDown, X: c.ptr.X, Y: c.ptr.Y})
		}

		// Pointer up: release + click if pointer ends inside widget.
		if c.ptr.IsJustUp {
			wasPressed := b.pressed
			if wasPressed {
				c.dispatch(w, Event{Type: EventPointerUp, X: c.ptr.X, Y: c.ptr.Y})
				if b.IsEnabled() && c.widgetHit(w, c.ptr.X, c.ptr.Y) {
					c.dispatch(w, Event{Type: EventClick, X: c.ptr.X, Y: c.ptr.Y})
				}
			}

			b.pressed = false
		}

		b.focused = (c.Focused() == w) && b.IsEnabled() && w.Focusable()
	}
}

func (c *Context) Draw(dst *ebiten.Image) {
	if c.root == nil {
		return
	}

	c.root.Draw(c, dst)
	// Overlay pass delegated to root/layout
	if ow, ok := any(c.root).(interface{ DrawOverlay(*Context, *ebiten.Image) }); ok {
		ow.DrawOverlay(c, dst)
	}
}
