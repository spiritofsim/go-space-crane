package main

import (
	"github.com/ByteArena/box2d"
)

type CabinCfg struct {
	Dir Direction
}

type Cabin struct {
	PartBase
	cfg CabinCfg
}

func NewCabin(cfg CabinCfg) *Cabin {
	s := &Sprite{
		img:   cabinSprite.img,
		verts: Rotate(cfg.Dir.GetAng(), cabinSprite.verts...),
	}
	return &Cabin{
		PartBase: PartBase{sprite: s, dir: cfg.Dir},
	}
}

func (c *Cabin) GetPos() box2d.B2Vec2 {
	return c.pos
}

func (c *Cabin) Construct(ship *Ship, pos box2d.B2Vec2, size box2d.B2Vec2) {
	c.pos = pos
	pos.OperatorPlusInplace(box2d.B2Vec2MulScalar(0.5, size).OperatorNegate())
	pos.OperatorPlusInplace(box2d.MakeB2Vec2(0.5, 0.5))
	verts := Translate(pos, c.sprite.verts...)

	shape := box2d.MakeB2PolygonShape()
	shape.Set(verts, len(verts))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter = box2d.MakeB2Filter()
	fd.Shape = &shape
	fd.Density = FixtureDensity
	fd.Restitution = FixtureRestitution
	ship.body.CreateFixtureFromDef(&fd)
	c.ship = ship
}
