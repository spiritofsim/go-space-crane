package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

type Cargo struct {
	*GameObj
	id       string
	size     box2d.B2Vec2
	platform *Platform
	tasks    []Task
}

func NewCargo(id string, world *box2d.B2World, pos box2d.B2Vec2) *Cargo {
	// Replace cargo image with new one with ID
	img := ebiten.NewImage(cargoSprite.img.Size())
	img.DrawImage(cargoSprite.img, nil)
	textBounds := text.BoundString(cargoFace, id)
	cargoSize := box2d.MakeB2Vec2(300, 300)
	text.Draw(img, id, cargoFace, int(cargoSize.X/2)-textBounds.Max.X/2, int(cargoSize.Y/2)-textBounds.Min.Y/2, color.White)
	sprite := NewSprite(img, cargoSprite.vertsSet)

	gobj := NewGameObj(
		world,
		sprite,
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

func (c *Cargo) Update() {
	if c.platform != nil && FloatEquals(c.GetVel(), 0) {
		for _, task := range c.tasks {
			task.CargoOnPlatform(c, c.platform)
		}
	}
}
