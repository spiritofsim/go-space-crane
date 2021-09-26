package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

type Terrain struct {
	body *box2d.B2Body
	img  *ebiten.Image
}

func NewTerrain(
	world *box2d.B2World,
	sprite Sprite) *Terrain {

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(0, 0)
	bd.Type = box2d.B2BodyType.B2_staticBody
	body := world.CreateBody(&bd)

	for _, verts := range sprite.vertsSet {
		shape := box2d.MakeB2ChainShape()
		shape.CreateLoop(verts, len(verts))

		fd := box2d.MakeB2FixtureDef()
		fd.Filter = box2d.MakeB2Filter()
		fd.Shape = &shape
		fd.Friction = DefaultFriction
		fd.Density = DefaultFixtureDensity
		fd.Restitution = DefaultFixtureRestitution
		body.CreateFixtureFromDef(&fd)
	}

	return &Terrain{
		body: body,
		img:  sprite.img,
	}
}

func (g *Terrain) Draw(screen *ebiten.Image, cam Cam) {
	opts := &ebiten.DrawImageOptions{}

	//bounds := g.img.Bounds()
	// PS resolution: 100px/cm
	opts.GeoM.Scale(0.01, 0.01)
	opts.GeoM.Translate(-cam.Pos.X, -cam.Pos.Y)
	opts.GeoM.Scale(cam.Zoom, cam.Zoom)
	opts.GeoM.Rotate(cam.Ang)
	opts.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	screen.DrawImage(g.img, opts)
}
