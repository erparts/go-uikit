package widget

import (
	"math"
	"strings"

	"github.com/erparts/go-uikit/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
)

// TextArea is a multi-line text editor with internal vertical scrolling.
type TextArea struct {
	base  ui.Base
	theme *ui.Theme

	text        string
	placeholder string

	lines int

	Scroll ui.Scroller

	// Caret config (end-of-text caret)
	CaretWidthPx  int
	CaretBlinkMs  int
	CaretMarginPx int
	caretTick     int
}

func NewTextArea(placeholder string) *TextArea {
	cfg := ui.NewWidgetBaseConfig()

	t := &TextArea{
		base:          ui.NewBase(cfg),
		placeholder:   placeholder,
		lines:         5,
		CaretWidthPx:  2,
		CaretBlinkMs:  600,
		CaretMarginPx: 0,
	}

	t.Scroll = ui.NewScroller()
	t.Scroll.Scrollbar = ui.ScrollbarAlways // textarea typically shows it (content-dependent)
	return t
}

func (t *TextArea) Base() *ui.Base  { return &t.base }
func (t *TextArea) Focusable() bool { return true }
func (t *TextArea) WantsIME() bool  { return true }

func (t *TextArea) Text() string     { return t.text }
func (t *TextArea) SetText(s string) { t.text = s }

func (t *TextArea) SetLines(n int) {
	if n < 1 {
		n = 1
	}
	t.lines = n
}

func (t *TextArea) Measure() ui.Rect { return t.base.Rect }

func (t *TextArea) SetFrame(x, y, w int) {
	if t.theme == nil {
		t.base.Rect = ui.Rect{X: x, Y: y, W: w, H: t.base.Rect.H}
		return
	}

	met, _ := ui.MetricsPx(t.theme.Font, t.theme.FontPx)
	lines := t.lines
	if lines <= 0 {
		lines = 5
	}

	controlH := t.theme.PadY*2 + lines*met.Height
	if controlH < t.theme.ControlH {
		controlH = t.theme.ControlH
	}

	totalH := controlH
	if ok, errTxt := t.base.IsInvalid(); ok && errTxt != "" {
		em, _ := ui.MetricsPx(t.theme.Font, t.theme.ErrorFontPx)
		totalH += t.theme.ErrorGap + em.Height
	}

	t.base.Rect = ui.Rect{X: x, Y: y, W: w, H: totalH}
}

func (t *TextArea) Update(ctx *ui.Context) {
	t.theme = ctx.Theme
	if t.base.Rect.W > 0 && t.base.Rect.H == 0 {
		t.SetFrame(t.base.Rect.X, t.base.Rect.Y, t.base.Rect.W)
	}

	if t.base.Focused() && t.base.IsEnabled() {
		t.caretTick++
	} else {
		t.caretTick = 0
	}

	// Compute content viewport and metrics
	_, content := t.controlAndContentRects(ctx)
	met, _ := ui.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)

	// Compute content height (scrollable)
	lineCount := 1
	if t.text != "" {
		lineCount = 1 + strings.Count(t.text, "\n")
	}
	contentH := lineCount * met.Height
	if contentH < content.H {
		contentH = content.H
	}

	// Update scroller using the content rect as viewport
	t.Scroll.Update(ctx, content, contentH)

	// Editing only when focused
	if !t.base.Focused() || !t.base.IsEnabled() {
		return
	}

	changed := false

	for _, r := range ebiten.AppendInputChars(nil) {
		if r == '\b' || r == 0x7f {
			t.backspace()
			changed = true
			continue
		}
		if r == '\n' || r == '\r' {
			t.text += "\n"
			changed = true
			continue
		}
		if r < 0x20 {
			continue
		}
		t.text += string(r)
		changed = true
	}

	// Optional: keep these for platforms that don't send backspace/enter via chars
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		t.backspace()
		changed = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyKPEnter) {
		t.text += "\n"
		changed = true
	}

	// Scroll-to-caret (caret at end of text)
	if changed {
		parts := strings.Split(t.text, "\n")
		lastIdx := len(parts) - 1
		caretBottom := (lastIdx + 1) * met.Height // bottom of last line cell

		minScroll := caretBottom - content.H
		if minScroll < 0 {
			minScroll = 0
		}
		if t.Scroll.ScrollY < minScroll {
			t.Scroll.ScrollY = minScroll
		}

		// clamp just in case
		maxScroll := contentH - content.H
		if maxScroll < 0 {
			maxScroll = 0
		}
		if t.Scroll.ScrollY > maxScroll {
			t.Scroll.ScrollY = maxScroll
		}
	}
}

func (t *TextArea) backspace() {
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

func (t *TextArea) controlAndContentRects(ctx *ui.Context) (ctrl ui.Rect, content ui.Rect) {
	r := t.base.Rect

	errorH := 0
	if ok, errTxt := t.base.IsInvalid(); ok && errTxt != "" {
		em, _ := ui.MetricsPx(ctx.Theme.Font, ctx.Theme.ErrorFontPx)
		errorH = ctx.Theme.ErrorGap + em.Height
	}

	ctrlH := r.H - errorH
	if ctrlH < 0 {
		ctrlH = 0
	}

	ctrl = ui.Rect{X: r.X, Y: r.Y, W: r.W, H: ctrlH}
	content = ctrl.Inset(ctx.Theme.PadX, ctx.Theme.PadY)
	return ctrl, content
}

func (t *TextArea) errorRect(ctx *ui.Context) ui.Rect {
	if ok, _ := t.base.IsInvalid(); !ok {
		return ui.Rect{}
	}

	ctrl, _ := t.controlAndContentRects(ctx)

	em, _ := ui.MetricsPx(ctx.Theme.Font, ctx.Theme.ErrorFontPx)
	return ui.Rect{
		X: ctrl.X,
		Y: ctrl.Bottom() + ctx.Theme.ErrorGap,
		W: ctrl.W,
		H: em.Height,
	}
}

func (t *TextArea) Draw(ctx *ui.Context, dst *ebiten.Image) {
	t.theme = ctx.Theme
	if t.base.Rect.W > 0 && t.base.Rect.H == 0 {
		t.SetFrame(t.base.Rect.X, t.base.Rect.Y, t.base.Rect.W)
	}

	ctrl, content := t.controlAndContentRects(ctx)

	// Surface
	bg := ctx.Theme.Surface
	if !t.base.IsEnabled() {
		bg = ctx.Theme.SurfacePressed
	} else if t.base.Pressed() {
		bg = ctx.Theme.SurfacePressed
	} else if t.base.Hovered() {
		bg = ctx.Theme.SurfaceHover
	}

	t.base.DrawRoundedRect(dst, ctrl, ctx.Theme.Radius, bg)

	// Border (red when invalid)
	borderCol := ctx.Theme.Border
	if ok, _ := t.base.IsInvalid(); ok {
		borderCol = ctx.Theme.ErrorBorder
	}

	t.base.DrawRoundedBorder(dst, ctrl, ctx.Theme.Radius, ctx.Theme.BorderW, borderCol)

	// Focus ring
	if t.base.Focused() && t.base.IsEnabled() {
		//drawFocusRing(dst, ctrl, ctx.Theme.Radius, ctx.Theme.FocusRingGap, ctx.Theme.FocusRingW, ctx.Theme.Focus)
	}

	// Clip to content area
	sub := dst.SubImage(content.ImageRect()).(*ebiten.Image)
	ox, oy := sub.Bounds().Min.X, sub.Bounds().Min.Y

	// Ensure renderer state for content
	ctx.Text.SetFont(ctx.Theme.Font)
	ctx.Text.SetSize(float64(ctx.Theme.FontPx))
	ctx.Text.SetAlign(etxt.Left)

	met, _ := ui.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	startY := -t.Scroll.ScrollY

	// Placeholder
	drawStr := t.text
	col := ctx.Theme.Text
	if drawStr == "" && !t.base.Focused() {
		drawStr = t.placeholder
		col = ctx.Theme.MutedText
	}
	ctx.Text.SetColor(col)

	lines := strings.Split(drawStr, "\n")
	for i, line := range lines {
		y := startY + i*met.Height + met.Ascent
		// IMPORTANT: SubImage keeps absolute coordinates -> use (ox, oy)
		ctx.Text.Draw(sub, line, ox, oy+y)
	}

	// Scrollbar (SubImage-safe after patch below)
	lineCount := 1
	if t.text != "" {
		lineCount = 1 + strings.Count(t.text, "\n")
	}
	contentH := lineCount * met.Height
	if contentH < content.H {
		contentH = content.H
	}
	t.Scroll.DrawBar(sub, ctx.Theme, content.W, content.H, contentH)

	// Caret at end (approx). Draw on sub so it's clipped.
	if t.base.Focused() && t.base.IsEnabled() && t.CaretWidthPx > 0 {
		blinkFrames := int(math.Max(1, float64(t.CaretBlinkMs)/1000.0*60.0))
		if (t.caretTick/blinkFrames)%2 == 0 {
			lastLine := ""
			lastIdx := 0
			if t.text != "" {
				parts := strings.Split(t.text, "\n")
				lastIdx = len(parts) - 1
				lastLine = parts[lastIdx]
			}

			wBefore := ui.MeasureStringPx(ctx.Theme.Font, ctx.Theme.FontPx, lastLine)
			cx := wBefore + t.CaretMarginPx
			cy := (lastIdx * met.Height) - t.Scroll.ScrollY

			if cx < 0 {
				cx = 0
			}
			maxX := content.W - t.CaretWidthPx
			if maxX < 0 {
				maxX = 0
			}
			if cx > maxX {
				cx = maxX
			}

			if cy+met.Height >= 0 && cy <= content.H {
				// IMPORTANT: SubImage keeps absolute coords -> add (ox,oy)
				vector.DrawFilledRect(
					sub,
					float32(ox+cx),
					float32(oy+cy),
					float32(t.CaretWidthPx),
					float32(met.Height),
					ctx.Theme.Caret,
					false,
				)
			}
		}
	}

	// Validation message
	//if t.base.Invalid && t.base.ErrorText != "" {
	//err := t.errorRect(ctx)
	//	drawErrorText(ctx, dst, err, t.base.ErrorText)
	//}
}

// SetTheme allows layouts to provide Theme before SetFrame is called.
func (t *TextArea) SetTheme(theme *ui.Theme) { t.theme = theme }
