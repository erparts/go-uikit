package ui

// Base contains shared widget state.
// Height is always Theme.ControlH; external layout controls only X/Y/Width.
type Base struct {
	Rect    Rect
	Visible bool
	Enabled bool

	hovered bool
	pressed bool
	focused bool
}

func NewBase() Base {
	return Base{
		Visible: true,
		Enabled: true,
	}
}

func (b *Base) SetRectByWidth(theme *Theme, x, y, w int) {
	if w < 0 {
		w = 0
	}
	b.Rect = Rect{X: x, Y: y, W: w, H: theme.ControlH}
}

func (b *Base) Hovered() bool { return b.hovered }
func (b *Base) Pressed() bool { return b.pressed }
func (b *Base) Focused() bool { return b.focused }

func (b *Base) SetEnabled(v bool) { b.Enabled = v }
func (b *Base) SetVisible(v bool) { b.Visible = v }
