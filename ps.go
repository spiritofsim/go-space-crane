package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

type ParticleSystem interface {
	Emit(pos box2d.B2Vec2, dir float64, angDisp float64)
	Update()
	Draw(screen *ebiten.Image, cam Cam)
}
