package demo

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"

	"github.com/erparts/go-uikit"
	"github.com/erparts/go-uikit/widget"
)

type Game struct {
	stack    *uikit.StackLayout
	grid     *uikit.GridLayout
	useGrid  bool
	contentH int
	ime      uikit.IMEBridge

	winW, winH int

	scale    uikit.Scale
	renderer *etxt.Renderer
	theme    *uikit.Theme
	ctx      *uikit.Context

	title        *widget.Label
	txtA         *widget.TextInput
	txtB         *widget.TextInput
	txtDis       *widget.TextInput
	ta           *widget.TextArea
	sel          *widget.Select
	box          *widget.Container
	chkA         *widget.Checkbox
	chkDis       *widget.Checkbox
	chkGrid      *widget.Checkbox
	btnA         *widget.Button
	btnDis       *widget.Button
	footer       *widget.Label
	focusInfo    *widget.Label
	exampleLabel *widget.Label
}

func mustFont() *sfnt.Font {
	f, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	return f
}

func New() *Game { return &Game{} }

// SetIMEBridge can be called from mobile bindings to enable keyboard show/hide.
func (g *Game) SetIMEBridge(b uikit.IMEBridge) {
	g.ime = b
	if g.ctx != nil {
		g.ctx.SetIMEBridge(b)
	}
}

func (g *Game) initOnce() {
	if g.renderer != nil {
		return
	}

	g.renderer = etxt.NewRenderer()
	g.renderer.Utils().SetCache8MiB()
	g.renderer.SetAlign(etxt.Left)
	g.renderer.Glyph().SetMissHandler(etxt.OnMissNotdef)

	f := mustFont()
	g.renderer.SetFont(f)

	g.theme = uikit.NewTheme(f, 20)

	g.ctx = uikit.NewContext(g.theme, g.renderer, g.ime)
	g.stack = uikit.NewStackLayout(g.theme)
	g.grid = uikit.NewGridLayout(g.theme)

	g.title = widget.NewLabel(g.theme, "UI Kit Demo (TPS: %0.2f - FPS: %0.2f")
	g.focusInfo = widget.NewLabel(g.theme, "")
	g.exampleLabel = widget.NewLabel(g.theme, "Label example: static helper text")

	g.txtA = widget.NewTextInput(g.theme, "Type here…")
	g.txtA.SetDefault("Hello Ebiten UI")

	g.txtB = widget.NewTextInput(g.theme, "Search…")

	g.txtDis = widget.NewTextInput(g.theme, "Disabled input")
	g.txtDis.SetDefault("Disabled value")
	g.txtDis.SetEnabled(false)

	g.ta = widget.NewTextArea(g.theme, "Multi-line…")
	g.ta.SetLines(5)
	g.ta.SetText("Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7")

	g.sel = widget.NewSelect(g.theme, []string{"Option A", "Option B", "Option C", "Option D", "Option E", "Option F"})

	g.box = widget.NewContainer(g.theme)
	g.box.SetHeight(200)
	g.box.OnDraw = func(ctx *uikit.Context, dst *ebiten.Image, content image.Rectangle) {

		lines := []string{
			"Custom container (user content)",
			"",
			"Select value: " + g.sel.Value(),
			"Search text: " + g.txtB.Text(),
			"TextArea chars: " + fmt.Sprintf("%d", len([]rune(g.ta.Text()))),
		}

		dst = dst.SubImage(content).(*ebiten.Image)

		met, _ := uikit.MetricsPx(ctx.Theme.Font, ctx.Theme.FontPx)
		y := (content.Min.Y) + met.Ascent
		x := content.Min.X

		ctx.Text.SetColor(ctx.Theme.MutedText)
		for _, ln := range lines {
			ctx.Text.Draw(dst, ln, x, y)
			y += met.Height
		}
	}

	g.chkA = widget.NewCheckbox(g.theme, "Enable main button")
	g.chkA.SetChecked(true)

	g.chkDis = widget.NewCheckbox(g.theme, "Disabled checkbox")
	g.chkDis.SetChecked(true)
	g.chkDis.SetEnabled(false)

	g.chkGrid = widget.NewCheckbox(g.theme, "Use grid layout")

	g.btnA = widget.NewButton(g.theme, "Action (enabled)")
	g.btnA.OnClick = func() {
		g.footer.SetText("Button clicked!")
	}

	g.btnDis = widget.NewButton(g.theme, "Action (disabled)")
	g.btnDis.SetEnabled(false)

	g.footer = widget.NewLabel(g.theme, "")

	g.ctx.Add(g.title)
	g.ctx.Add(g.focusInfo)
	g.ctx.Add(g.chkGrid)

	g.ctx.Add(g.stack)
	g.ctx.Add(g.grid)

	contentWidgets := []uikit.Widget{
		g.exampleLabel,
		g.txtA,
		g.txtB,
		g.txtDis,
		g.ta,
		g.sel,
		g.box,
		g.chkA,
		g.chkDis,
		g.btnA,
		g.btnDis,
	}

	g.stack.SetChildren(contentWidgets)
	g.grid.SetChildren(contentWidgets)

}

func (g *Game) Layout(outW, outH int) (int, int) {
	g.initOnce()

	g.winW, g.winH = outW, outH

	dev := ebiten.DeviceScaleFactor()
	if dev <= 0 {
		dev = 1
	}

	// Optional: make UI a bit larger on small screens.
	uiScale := 1.0
	minSide := float64(outW)
	if float64(outH) < minSide {
		minSide = float64(outH)
	}
	if minSide <= 520 {
		uiScale = 1.5
	} else if minSide <= 720 {
		uiScale = 1.25
	}

	g.scale = uikit.Scale{Device: dev, UI: uiScale}
	g.ctx.SetScale(g.scale)

	// IMPORTANT: renderer scale must stay 1 (do not double scale).
	g.renderer.SetScale(1)

	return outW, outH
}

func (g *Game) Update() error {
	g.initOnce()

	title := fmt.Sprintf("UI Kit Demo (TPS: %0.2f - FPS: %0.2f)", ebiten.ActualTPS(), ebiten.ActualFPS())
	g.title.SetText(title)

	// Layout constants
	x := g.theme.SpaceL
	y := g.theme.SpaceL
	w := g.winW - g.theme.SpaceL*2
	if w < 0 {
		w = 0
	}

	// Header
	g.ctx.Root().SetFrame(x, y, w)
	g.title.SetFrame(x, y, w)
	y += g.title.Measure().Dy() + g.theme.SpaceS

	fw := g.ctx.Focused()
	if fw == nil {
		g.focusInfo.SetText("Focused: (none) — tap a widget or TAB")
	} else {
		g.focusInfo.SetText(fmt.Sprintf("Focused: %T", fw))
	}
	g.focusInfo.SetFrame(x, y, w)
	y += g.focusInfo.Measure().Dy() + g.theme.SpaceM

	// Scrollable viewport for the content widgets
	footerH := g.footer.Measure().Dy()
	viewportH := g.winH - y - g.theme.SpaceM - footerH - g.theme.SpaceM
	if viewportH < g.theme.ControlH {
		viewportH = g.theme.ControlH
	}

	g.stack.SetFrame(x, y, w)
	g.stack.SetHeight(viewportH)
	g.grid.SetFrame(x, y, w)
	g.grid.SetHeight(viewportH)

	// Demo validations
	if g.txtB.Text() == "" {
		g.txtB.Base().SetInvalid("Required")
	} else {
		g.txtB.Base().ClearInvalid()
	}

	if strings.TrimSpace(g.ta.Text()) == "" {
		g.txtB.Base().SetInvalid("Required")
	} else {
		g.txtB.Base().ClearInvalid()
	}

	selVal := g.sel.Value()
	if selVal == "Option A" || selVal == "" {
		g.sel.Base().SetInvalid("Option A is not allowed")
	} else {
		g.txtB.Base().ClearInvalid()
	}

	// Enable button based on checkbox state
	g.btnA.SetEnabled(g.chkA.Checked())
	g.useGrid = g.chkGrid.Checked()

	// Swap root layout (stack or grid)
	if g.useGrid {
		g.grid.Base().SetVisible(true)
		g.stack.Base().SetVisible(false)
	} else {
		g.stack.Base().SetVisible(true)
		g.grid.Base().SetVisible(false)
	}

	// Update widgets (events, focus, etc.)
	g.ctx.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 22, 26, 255})
	g.ctx.Draw(screen)
}
