package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Platform struct {
	*GameObj
	id   string
	fuel float64
	ship *Ship
}

func NewPlatform(id string, world *box2d.B2World, pos box2d.B2Vec2, size box2d.B2Vec2, fuel float64) *Platform {
	verts := []box2d.B2Vec2{
		{-size.X / 2, -size.Y / 2},
		{size.X / 2, -size.Y / 2},
		{size.X / 2, size.Y / 2},
		{-size.X / 2, size.Y / 2},
	}
	gobj := NewGameObj(
		world,
		NewSprite(nil, [][]box2d.B2Vec2{verts}),
		pos,
		0,
		0,
		box2d.B2Vec2_zero,
		DefaultFriction,
		DefaultFixtureDensity,
		DefaultFixtureRestitution)

	platform := &Platform{
		GameObj: gobj,
		id:      id,
		fuel:    fuel,
	}
	platform.body.SetUserData(platform)
	return platform
}

func (p *Platform) Draw(screen *ebiten.Image, cam Cam) {
	pos := p.body.GetPosition()
	x := box2d.MakeB2Vec2(-1, -0.5)
	x = cam.Project(x, pos, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Fuel: %0.2f", p.fuel), int(x.X), int(x.Y))
}
