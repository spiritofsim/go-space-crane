package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Cargo struct {
	body *box2d.B2Body
}

func NewCargo(world *box2d.B2World, pos box2d.B2Vec2, size float64) *Cargo {
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(pos.X, pos.Y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	body := world.CreateBody(&bd)

	verts := []box2d.B2Vec2{
		{-size / 2, -size / 2},
		{size / 2, -size / 2},
		{size / 2, size / 2},
		{-size / 2, size / 2},
	}

	shape := box2d.MakeB2PolygonShape()
	shape.Set(verts, len(verts))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter = box2d.MakeB2Filter()
	fd.Shape = &shape
	fd.Density = FixtureDensity
	fd.Restitution = FixtureRestitution
	body.CreateFixtureFromDef(&fd)

	cargo := &Cargo{
		body: body,
	}
	body.SetUserData(cargo)
	return cargo
}

func (c *Cargo) Draw(screen *ebiten.Image, cam Cam) {
	DrawDebugBody(screen, c.body, cam, color.White)
}
