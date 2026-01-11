package ui

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// TextInput is a single-line input box (no label).
// Height and proportions come from Theme; external layout controls only width.
type TextInput struct {
	base  Base
	theme *Theme

	// Value / defaults
	text                 string
	DefaultText          string
	RestoreDefaultOnBlur bool

	placeholder string

	// Caret config is exclusive to TextInput, as requested.
	CaretWidthPx  int
	CaretBlinkMs  int
	CaretMarginPx int

	caretTick int
}

func NewTextInput(placeholder string) *TextInput {
	return &TextInput{
		base:          NewBase(),
		placeholder:   placeholder,
		CaretWidthPx:  2,
		CaretBlinkMs:  600,
		CaretMarginPx: 0,
	}
}

func (t *TextInput) Base() *Base     { return &t.base }
func (t *TextInput) Focusable() bool { return true }
func (t *TextInput) WantsIME() bool  { return true }

func (t *TextInput) SetRectByWidth(x, y, w int) {
	if t.theme != nil {
		t.base.SetRectByWidth(t.theme, x, y, w)
		return
	}
	t.base.Rect = Rect{X: x, Y: y, W: w, H: 0}
}

func (t *TextInput) SetEnabled(v bool) { t.base.SetEnabled(v) }
func (t *TextInput) SetVisible(v bool) { t.base.SetVisible(v) }

func (t *TextInput) Text() string { return t.text }

// SetText sets the current value (does not change DefaultText).
func (t *TextInput) SetText(s string) { t.text = s }

// SetDefault sets DefaultText and also resets the current value to it.
func (t *TextInput) SetDefault(s string) {
	t.DefaultText = s
	t.text = s
}

func (t *TextInput) Reset() { t.text = t.DefaultText }

// HandleEvent allows focus transitions to apply policies (default restore).
func (t *TextInput) HandleEvent(ctx *Context, e Event) {
	if e.Type == EventFocusLost && t.RestoreDefaultOnBlur && t.text == "" {
		t.text = t.DefaultText
	}
}

func (t *TextInput) Update(ctx *Context) {
	t.theme = ctx.Theme
	if t.base.Rect.H == 0 {
		t.base.SetRectByWidth(ctx.Theme, t.base.Rect.X, t.base.Rect.Y, t.base.Rect.W)
	}

	// Tick caret
	if t.base.focused && t.base.Enabled {
		t.caretTick++
	} else {
		t.caretTick = 0
	}

	// If not focused, ignore input
	if !t.base.focused || !t.base.Enabled {
		return
	}

	// Typed characters (desktop fallback).
	for _, r := range ebiten.AppendInputChars(nil) {
		// Handle common backspace representations defensively (desktop vs some IMEs)
		if r == '\b' || r == 0x7f {
			t.backspace()
			continue
		}
		// Skip control chars
		if r < 0x20 {
			continue
		}
		t.text += string(r)
	}

	// Backspace via key (desktop)
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		t.backspace()
	}
	// Delete (treat as backspace for end-of-text caret)
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		t.backspace()
	}

	// Enter: blur
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyKPEnter) {
		ctx.SetFocus(nil)
	}
}

func (t *TextInput) backspace() {
	if t.text == "" {
		return
	}
	rs := []rune(t.text)
	if len(rs) == 0 {
		t.text = ""
		return
	}
	t.text = string(rs[:len(rs)-1])
}

func (t *TextInput) Draw(ctx *Context, dst *ebiten.Image) {
	t.theme = ctx.Theme
	if t.base.Rect.H == 0 {
		t.base.SetRectByWidth(ctx.Theme, t.base.Rect.X, t.base.Rect.Y, t.base.Rect.W)
	}

	r := t.base.Rect

	// Surface
	bg := ctx.Theme.Surface
	if !t.base.Enabled {
		bg = ctx.Theme.SurfacePressed
	} else if t.base.pressed {
		bg = ctx.Theme.SurfacePressed
	} else if t.base.hovered {
		bg = ctx.Theme.SurfaceHover
	}
	drawRoundedRect(dst, r, ctx.Theme.Radius, bg)

	// Border
	drawRoundedBorder(dst, r, ctx.Theme.Radius, ctx.Theme.BorderW, ctx.Theme.Border)

	// Focus ring
	if t.base.focused && t.base.Enabled {
		drawFocusRing(dst, r, ctx.Theme.Radius, ctx.Theme.FocusRingGap, ctx.Theme.FocusRingW, ctx.Theme.Focus)
	}

	content := r.Inset(ctx.Theme.PadX, ctx.Theme.PadY)

	// Baseline
	met, _ := MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	baselineY := r.Y + (r.H-met.Height)/2 + met.Ascent

	// Text / placeholder
	drawStr := t.text
	textCol := ctx.Theme.Text
	if drawStr == "" && !t.base.focused {
		drawStr = t.placeholder
		textCol = ctx.Theme.MutedText
	}

	// Horizontal scroll so the end is visible (caret is always at end).
	textW := MeasureStringPx(ctx.Theme.Font, ctx.Theme.FontPx, drawStr)
	shiftX := 0
	if textW > content.W {
		shiftX = content.W - textW
	}

	ctx.Text.SetAlign(0) // Left
	ctx.Text.SetColor(textCol)
	ctx.Text.Draw(dst, drawStr, content.X+shiftX, baselineY)

	// Caret (rect, no extra space)
	if t.base.focused && t.base.Enabled && t.CaretWidthPx > 0 {
		blinkFrames := int(math.Max(1, float64(t.CaretBlinkMs)/1000.0*60.0))
		if (t.caretTick/blinkFrames)%2 == 0 {
			wBefore := MeasureStringPx(ctx.Theme.Font, ctx.Theme.FontPx, t.text)
			cx := content.X + shiftX + wBefore + t.CaretMarginPx
			cy := baselineY - met.Ascent
			caretH := met.Height

			// Clamp caret inside the content area
			if cx < content.X {
				cx = content.X
			}
			if cx > content.X+content.W {
				cx = content.X + content.W
			}
			vector.DrawFilledRect(dst, float32(cx), float32(cy), float32(t.CaretWidthPx), float32(caretH), ctx.Theme.Caret, false)
		}
	}
}
