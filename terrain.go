package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

type Terrain struct {
	body *box2d.B2Body
}

func (t *Terrain) Draw(screen *ebiten.Image, cam Cam) {
	//for _, body := range t.bodies {
	//	DrawDebugBody(screen, body, cam, color.White)
	//}
}
