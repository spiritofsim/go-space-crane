package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

type Part interface {
	GetBody() *box2d.B2Body
	Update(keys []ebiten.Key)
	Draw(screen *ebiten.Image, cam Cam)
	GetVel() float64
	GetVelVec() box2d.B2Vec2
	GetAng() float64
	GetPos() box2d.B2Vec2
}
