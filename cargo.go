package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

type Cargo struct {
	*GameObj
	id   string
	size box2d.B2Vec2
}

func NewCargo(id string, world *box2d.B2World, pos box2d.B2Vec2) *Cargo {
	gobj := NewGameObj(
		world,
		cargoSprite,
		pos,
		0,
		0,
		box2d.B2Vec2_zero,
		DefaultFriction,
		DefaultFixtureDensity,
		DefaultFixtureRestitution)

	cargo := &Cargo{
		GameObj: gobj,
		id:      id,
		size:    getShapeSize(cargoSprite.vertsSet[0]),
	}
	cargo.body.SetUserData(cargo)
	return cargo
}

func (c *Cargo) Draw(screen *ebiten.Image, cam Cam) {
	c.GameObj.Draw(screen, cam)
	pos := c.body.GetPosition()
	px := box2d.MakeB2Vec2(-c.size.X/2, 0)
	px = cam.Project(px, pos, 0)

	// TODO: draw text on image, then apply cam and copy image to screen
	msg := fmt.Sprintf("%v", c.id)
	text.Draw(screen, msg, face, int(px.X), int(px.Y), color.White)
}
