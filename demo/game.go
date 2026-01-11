package demo

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"

	"github.com/erparts/go-uikit/ui"
)

type Game struct {
	ime ui.IMEBridge
	// logical window size from Layout inputs
	winW, winH int

	// scale system (device * ui)
	scale ui.Scale

	renderer *etxt.Renderer
	theme    *ui.Theme
	ctx      *ui.Context

	// Widgets (showcase)
	title     *ui.Label
	txtA      *ui.TextInput
	txtB      *ui.TextInput
	txtDis    *ui.TextInput
	chkA      *ui.Checkbox
	chkDis    *ui.Checkbox
	btnA      *ui.Button
	btnDis    *ui.Button
	footer    *ui.Label
	focusInfo *ui.Label
}

func mustFont() *sfnt.Font {
	f, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	return f
}

// New returns an embeddable ebiten.Game.
func New() *Game { return &Game{} }

// SetIMEBridge can be called from mobile bindings to enable keyboard show/hide.
func (g *Game) SetIMEBridge(b ui.IMEBridge) {
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

	f := mustFont()
	g.renderer.SetFont(f)

	// Theme is defined in *physical pixels*. We render on a physical canvas for crisp UI.
	g.theme = ui.NewTheme(f, 20)

	// Desktop demo: IME bridge nil
	g.ctx = ui.NewContext(g.theme, g.renderer, g.ime)

	// Widgets
	g.title = ui.NewLabel("UI Kit Demo — consistent proportions (Theme-driven)")
	g.focusInfo = ui.NewLabel("")

	g.txtA = ui.NewTextInput("Type here…")
	g.txtA.SetDefault("Hello Ebiten UI")
	g.txtA.RestoreDefaultOnBlur = false

	g.txtB = ui.NewTextInput("Search…")
	g.txtB.SetText("")
	g.txtB.DefaultText = ""

	g.txtDis = ui.NewTextInput("Disabled input")
	g.txtDis.SetDefault("Disabled value")
	g.txtDis.SetEnabled(false)

	g.chkA = ui.NewCheckbox("Enable main button")
	g.chkA.SetChecked(true)

	g.chkDis = ui.NewCheckbox("Disabled checkbox")
	g.chkDis.SetChecked(true)
	g.chkDis.SetEnabled(false)

	g.btnA = ui.NewButton("Action (enabled)")
	g.btnA.OnClick = func() {
		g.footer.SetText("Button clicked!")
	}

	g.btnDis = ui.NewButton("Action (disabled)")
	g.btnDis.SetEnabled(false)

	g.footer = ui.NewLabel("")

	// Register in draw order
	g.ctx.Add(g.title)
	g.ctx.Add(g.focusInfo)
	g.ctx.Add(g.txtA)
	g.ctx.Add(g.txtB)
	g.ctx.Add(g.txtDis)
	g.ctx.Add(g.chkA)
	g.ctx.Add(g.chkDis)
	g.ctx.Add(g.btnA)
	g.ctx.Add(g.btnDis)
	g.ctx.Add(g.footer)
}

func (g *Game) Layout(outW, outH int) (int, int) {
	g.initOnce()

	// outsideWidth/outsideHeight are logical pixels.
	g.winW, g.winH = outW, outH

	g.scale = ui.Scale{Device: ebiten.DeviceScaleFactor(), UI: 1.0}

	// Render on a physical canvas for crisp UI.
	canvasW := int(math.Ceil(float64(outW) * g.scale.Device))
	canvasH := int(math.Ceil(float64(outH) * g.scale.Device))

	// etxt renderer should not apply extra scaling when we already render in physical pixels.
	g.renderer.SetScale(1)

	return canvasW, canvasH
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Layout in logical units, then convert to physical pixels.
	padding := 12
	x := g.scale.PxI(padding)
	y := g.scale.PxI(padding)
	w := g.scale.PxI(g.winW - padding*2)

	// Title / info
	g.title.SetRectByWidth(x, y, w)
	y += g.theme.ControlH + g.theme.SpaceS

	fw := g.ctx.Focused()
	if fw == nil {
		g.focusInfo.SetText("Focused: (none) — click a widget or TAB")
	} else {
		g.focusInfo.SetText(fmt.Sprintf("Focused: %T", fw))
	}
	g.focusInfo.SetRectByWidth(x, y, w)
	y += g.theme.ControlH + g.theme.SpaceM

	// Inputs
	g.txtA.SetRectByWidth(x, y, w)
	y += g.theme.ControlH + g.theme.SpaceS

	g.txtB.SetRectByWidth(x, y, w)
	y += g.theme.ControlH + g.theme.SpaceS

	g.txtDis.SetRectByWidth(x, y, w)
	y += g.theme.ControlH + g.theme.SpaceM

	// Checkboxes
	g.chkA.SetRectByWidth(x, y, w)
	y += g.theme.ControlH + g.theme.SpaceS

	g.chkDis.SetRectByWidth(x, y, w)
	y += g.theme.ControlH + g.theme.SpaceM

	// Buttons
	g.btnA.SetEnabled(g.chkA.Checked())
	g.btnA.SetRectByWidth(x, y, w)
	y += g.theme.ControlH + g.theme.SpaceS

	g.btnDis.SetRectByWidth(x, y, w)
	y += g.theme.ControlH + g.theme.SpaceM

	// Footer
	g.footer.SetRectByWidth(x, y, w)

	g.ctx.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 22, 26, 255})
	g.ctx.Draw(screen)
}
