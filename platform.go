package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

type Platform struct {
	*GameObj
	id   string
	fuel float64
	ship *Ship
	size box2d.B2Vec2
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
		size:    size,
	}
	platform.body.SetUserData(platform)
	return platform
}

func (p *Platform) Draw(screen *ebiten.Image, cam Cam) {
	pos := p.body.GetPosition()
	px := box2d.MakeB2Vec2(-p.size.X/2, 0)
	px = cam.Project(px, pos, 0)

	// TODO: draw text on image, then apply cam and copy image to screen
	// TODO: use text.BoundString()
	msg := fmt.Sprintf("%v : %v", p.id, int(p.fuel))
	text.Draw(screen, msg, face, int(px.X), int(px.Y), color.White)
}
