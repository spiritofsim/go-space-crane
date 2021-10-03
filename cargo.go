package main

import (
	"github.com/ByteArena/box2d"
)

type Cargo struct {
	*GameObj
	id string
}

func NewCargo(id string, world *box2d.B2World, pos box2d.B2Vec2, size box2d.B2Vec2) *Cargo {
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
	}
	cargo.body.SetUserData(cargo)
	return cargo
}
