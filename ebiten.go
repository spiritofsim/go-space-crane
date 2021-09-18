package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

func drawCross(screen *ebiten.Image, p box2d.B2Vec2, size float64, clr color.Color) {
	ebitenutil.DrawLine(screen, p.X, p.Y-size/2, p.X, p.Y+size/2, clr)
	ebitenutil.DrawLine(screen, p.X-size/2, p.Y, p.X+size/2, p.Y, clr)
}
