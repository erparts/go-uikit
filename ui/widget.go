package ui

import "github.com/hajimehoshi/ebiten/v2"

type Widget interface {
	Base() *Base
	Focusable() bool
	SetRectByWidth(x, y, w int)
	Update(ctx *Context)
	Draw(ctx *Context, dst *ebiten.Image)
}
