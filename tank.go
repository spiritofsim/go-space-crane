package main

import (
	"github.com/ByteArena/box2d"
)

type TankCfg struct {
	Fuel    float64
	MaxFuel float64
}

type Tank struct {
	PartBase
	Fuel    float64
	MaxFuel float64
}

func NewTank(cfg TankCfg) *Tank {
	return &Tank{
		PartBase: PartBase{img: tankImg, dir: DirectionRight},
		Fuel:     cfg.Fuel,
		MaxFuel:  cfg.MaxFuel,
	}
}

func (t *Tank) GetPos() box2d.B2Vec2 {
	return t.pos
}

func (t *Tank) Construct(ship *Ship, pos box2d.B2Vec2, size box2d.B2Vec2) {
	t.pos = pos
	verts := []box2d.B2Vec2{
		{-0.5, -0.25},
		{-0.25, -0.5},
		{0.25, -0.5},
		{0.5, -0.25},
		{0.5, 0.25},
		{0.25, 0.5},
		{-0.25, 0.5},
		{-0.5, 0.25},
	}

	pos.OperatorPlusInplace(box2d.B2Vec2MulScalar(0.5, size).OperatorNegate())
	pos.OperatorPlusInplace(box2d.MakeB2Vec2(0.5, 0.5))
	verts = Translate(pos, verts...)

	shape := box2d.MakeB2PolygonShape()
	shape.Set(verts, len(verts))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter = box2d.MakeB2Filter()
	fd.Shape = &shape
	fd.Density = FixtureDensity
	fd.Restitution = FixtureRestitution
	ship.body.CreateFixtureFromDef(&fd)
	t.ship = ship
}
