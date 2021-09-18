package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

type CabinCfg struct {
	Dir Direction
}

var cabinVerts = []box2d.B2Vec2{
	{-0.5, -0.25},
	{-0.25, -0.25},
	{0, -0.125},
	{0, 0.125},
	{-0.5, 0.25},
	{-0.25, 0.25},
}

type Cabin struct {
	PartBase
	cfg   CabinCfg
	verts []box2d.B2Vec2
}

func NewCabin(cfg CabinCfg, img *ebiten.Image) *Cabin {
	return &Cabin{
		PartBase: PartBase{img: img, dir: cfg.Dir},
		verts:    Rotate(cfg.Dir.GetAng(), cabinVerts...),
	}
}

func (c *Cabin) GetPos() box2d.B2Vec2 {
	return c.pos
}

func (c *Cabin) Construct(ship *Ship, pos box2d.B2Vec2, size box2d.B2Vec2) {
	c.pos = pos
	pos.OperatorPlusInplace(box2d.B2Vec2MulScalar(0.5, size).OperatorNegate())
	pos.OperatorPlusInplace(box2d.MakeB2Vec2(0.5, 0.5))
	verts := Translate(pos, c.verts...)

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
