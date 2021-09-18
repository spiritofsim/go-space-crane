package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

type Part interface {
	Construct(ship *Ship, pos box2d.B2Vec2, size box2d.B2Vec2)
	// GetPos returns part position relative to ship upper left corner
	GetPos() box2d.B2Vec2
	Update()
	Draw(screen *ebiten.Image, cam Cam)
}

type PartBase struct {
	img  *ebiten.Image
	ship *Ship
	pos  box2d.B2Vec2
	dir  Direction
}

func (p PartBase) Update() {
}

func (p PartBase) Draw(screen *ebiten.Image, cam Cam) {
	if p.img == nil {
		return
	}

	bounds := p.img.Bounds()
	shipPos := p.ship.body.GetPosition()

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-float64(bounds.Max.X/2), -float64(bounds.Max.Y/2))
	opts.GeoM.Scale(1/float64(bounds.Max.X), 1/float64(bounds.Max.Y))
	opts.GeoM.Rotate(p.dir.GetAng())
	opts.GeoM.Translate(p.pos.X-p.ship.size.X/2+0.5, p.pos.Y-p.ship.size.Y/2+0.5)
	opts.GeoM.Rotate(p.ship.body.GetAngle())
	opts.GeoM.Translate(shipPos.X, shipPos.Y)
	opts.GeoM.Translate(-cam.Pos.X, -cam.Pos.Y)
	opts.GeoM.Scale(cam.Zoom, cam.Zoom)
	opts.GeoM.Rotate(cam.Ang)
	opts.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	screen.DrawImage(p.img, opts)
}
