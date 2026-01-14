package widget

import (
	"image"
	"math"
	"strings"

	"github.com/erparts/go-uikit"
	"github.com/erparts/go-uikit/common"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
)

// TextArea is a multi-line text editor with internal vertical scrolling.
type TextArea struct {
	base uikit.Base

	text        string
	placeholder string

	lines int

	Scroll uikit.Scroller

	// Caret config (end-of-text caret)
	CaretWidthPx  int
	CaretBlinkMs  int
	CaretMarginPx int
	caretTick     int
}

func NewTextArea(theme *uikit.Theme, placeholder string) *TextArea {
	cfg := uikit.NewWidgetBaseConfig(theme)

	t := &TextArea{
		base:          uikit.NewBase(cfg),
		placeholder:   placeholder,
		lines:         5,
		CaretWidthPx:  2,
		CaretBlinkMs:  600,
		CaretMarginPx: 0,
	}

	t.Scroll = uikit.NewScroller()
	t.Scroll.Scrollbar = uikit.ScrollbarAlways // textarea typically shows it (content-dependent)
	return t
}

func (t *TextArea) Base() *uikit.Base { return &t.base }
func (t *TextArea) Focusable() bool   { return true }
func (t *TextArea) WantsIME() bool    { return true }

func (t *TextArea) Text() string     { return t.text }
func (t *TextArea) SetText(s string) { t.text = s }

func (t *TextArea) SetLines(n int) {
	if n < 1 {
		n = 1
	}
	t.lines = n
}

func (t *TextArea) Measure() image.Rectangle { return t.base.Rect }

func (t *TextArea) SetFrame(x, y, w int) {
	met, _ := uikit.MetricsPx(t.base.Theme().Font, t.base.Theme().FontPx)
	lines := t.lines
	if lines <= 0 {
		lines = 5
	}

	controlH := t.base.Theme().PadY*2 + lines*met.Height
	if controlH < t.base.Theme().ControlH {
		controlH = t.base.Theme().ControlH
	}

	totalH := controlH
	if ok, errTxt := t.base.IsInvalid(); ok && errTxt != "" {
		em, _ := uikit.MetricsPx(t.base.Theme().Font, t.base.Theme().ErrorFontPx)
		totalH += t.base.Theme().ErrorGap + em.Height
	}

	if w < 0 {
		w = 0
	}

	t.base.Rect = image.Rect(x, y, x+w, y+totalH)
}

func (t *TextArea) Update(ctx *uikit.Context) {
	if t.base.Rect.Dx() > 0 && t.base.Rect.Dy() == 0 {
		t.SetFrame(t.base.Rect.Min.X, t.base.Rect.Min.Y, t.base.Rect.Dx())
	}

	if t.base.IsFocused() && t.base.IsEnabled() {
		t.caretTick++
	} else {
		t.caretTick = 0
	}

	// Compute content viewport and metrics
	_, content := t.controlAndContentRects(ctx)
	met, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)

	// Compute content height (scrollable)
	lineCount := 1
	if t.text != "" {
		lineCount = 1 + strings.Count(t.text, "\n")
	}
	contentH := lineCount * met.Height
	if contentH < content.Dy() {
		contentH = content.Dy()
	}

	// Update scroller using the content rect as viewport
	t.Scroll.Update(ctx, content, contentH)

	// Editing only when focused
	if !t.base.IsFocused() || !t.base.IsEnabled() {
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

		minScroll := caretBottom - content.Dy()
		if minScroll < 0 {
			minScroll = 0
		}
		if t.Scroll.ScrollY < minScroll {
			t.Scroll.ScrollY = minScroll
		}

		// clamp just in case
		maxScroll := contentH - content.Dy()
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

func (t *TextArea) controlAndContentRects(ctx *uikit.Context) (ctrl image.Rectangle, content image.Rectangle) {
	r := t.base.Rect

	errorH := 0
	if ok, errTxt := t.base.IsInvalid(); ok && errTxt != "" {
		em, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.ErrorFontPx)
		errorH = ctx.Theme.ErrorGap + em.Height
	}

	ctrlH := r.Dy() - errorH
	if ctrlH < 0 {
		ctrlH = 0
	}

	ctrl = image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+ctrlH)
	content = common.Inset(ctrl, ctx.Theme.PadX, ctx.Theme.PadY)
	return ctrl, content
}

func (t *TextArea) errorRect(ctx *uikit.Context) image.Rectangle {
	if ok, _ := t.base.IsInvalid(); !ok {
		return image.Rectangle{}
	}

	ctrl, _ := t.controlAndContentRects(ctx)

	em, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.ErrorFontPx)

	top := ctrl.Max.Y + ctx.Theme.ErrorGap
	return image.Rect(ctrl.Min.X, top, ctrl.Max.X, top+em.Height)
}

func (t *TextArea) Draw(ctx *uikit.Context, dst *ebiten.Image) {
	if t.base.Rect.Dx() > 0 && t.base.Rect.Dy() == 0 {
		t.SetFrame(t.base.Rect.Min.X, t.base.Rect.Min.Y, t.base.Rect.Dx())
	}

	r, content := t.controlAndContentRects(ctx)

	t.base.DrawSurfece(ctx, dst, r)
	t.base.DrawBoder(ctx, dst, r)
	t.base.DrawFocus(ctx, dst, r)

	// Clip to content area
	sub := dst.SubImage(content).(*ebiten.Image)
	ox, oy := sub.Bounds().Min.X, sub.Bounds().Min.Y

	// Ensure renderer state for content
	ctx.Text.SetFont(ctx.Theme.Font)
	ctx.Text.SetSize(float64(ctx.Theme.FontPx))
	ctx.Text.SetAlign(etxt.Left)

	met, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
	startY := -t.Scroll.ScrollY

	// Placeholder
	drawStr := t.text
	col := ctx.Theme.Text
	if drawStr == "" && !t.base.IsFocused() {
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
	if contentH < content.Dy() {
		contentH = content.Dy()
	}

	t.Scroll.DrawBar(sub, ctx.Theme, content.Dx(), content.Dy(), contentH)

	// Caret at end (approx). Draw on sub so it's clipped.
	if t.base.IsFocused() && t.base.IsEnabled() && t.CaretWidthPx > 0 {
		blinkFrames := int(math.Max(1, float64(t.CaretBlinkMs)/1000.0*60.0))
		if (t.caretTick/blinkFrames)%2 == 0 {
			lastLine := ""
			lastIdx := 0
			if t.text != "" {
				parts := strings.Split(t.text, "\n")
				lastIdx = len(parts) - 1
				lastLine = parts[lastIdx]
			}

			wBefore := uikit.MeasureStringPx(ctx.Theme.Font, ctx.Theme.FontPx, lastLine)
			cx := wBefore + t.CaretMarginPx
			cy := (lastIdx * met.Height) - t.Scroll.ScrollY

			if cx < 0 {
				cx = 0
			}
			maxX := content.Dx() - t.CaretWidthPx
			if maxX < 0 {
				maxX = 0
			}
			if cx > maxX {
				cx = maxX
			}

			if cy+met.Height > 0 && cy <= content.Dy() {
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

	t.base.DrawInvalid(ctx, dst, r)
}

func (t *TextArea) HitTest(ctx *uikit.Context, x, y int) bool {
	ctrl, _ := t.controlAndContentRects(ctx)
	if common.Contains(ctrl, x, y) {
		return true
	}

	return false
}
