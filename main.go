package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/erparts/go-uikit/demo"
)

func main() {
	ebiten.SetWindowSize(520, 280)
	ebiten.SetWindowTitle("uikitdemo (consistent proportions)")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(demo.New()); err != nil {
		log.Fatal(err)
	}
}
