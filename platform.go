package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Platform struct {
	body *box2d.B2Body
	fuel float64

	ship *Ship
}

func NewPlatform(world *box2d.B2World, pos box2d.B2Vec2, size box2d.B2Vec2, fuel float64) *Platform {
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(pos.X, pos.Y)
	bd.Type = box2d.B2BodyType.B2_staticBody
	//bd.AllowSleep = false
	body := world.CreateBody(&bd)

	verts := []box2d.B2Vec2{
		{-size.X / 2, -size.Y / 2},
		{size.X / 2, -size.Y / 2},
		{size.X / 2, size.Y / 2},
		{-size.X / 2, size.Y / 2},
	}

	shape := box2d.MakeB2PolygonShape()
	shape.Set(verts, len(verts))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter = box2d.MakeB2Filter()
	fd.Shape = &shape
	fd.Density = FixtureDensity
	fd.Restitution = FixtureRestitution
	body.CreateFixtureFromDef(&fd)

	platform := &Platform{
		fuel: fuel,
		body: body,
	}
	body.SetUserData(platform)
	return platform
}

func (p *Platform) Draw(screen *ebiten.Image, cam Cam) {
	//DrawDebugBody(screen, p.body, cam, color.White)

	pos := p.body.GetPosition()
	x := box2d.MakeB2Vec2(-1, -0.5)
	x = cam.Project(x, pos, 0)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Fuel: %0.2f", p.fuel), int(x.X), int(x.Y))
}
