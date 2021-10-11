package main

import (
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

	bounds := text.BoundString(platformFace, id)
	img := ebiten.NewImage(-bounds.Min.X+bounds.Max.X, -bounds.Min.Y+bounds.Max.Y+50)
	text.Draw(img, id, platformFace, -bounds.Min.X, -bounds.Min.Y+50, color.White)

	gobj := NewGameObj(
		world,
		NewSprite(img, [][]box2d.B2Vec2{verts}),
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
