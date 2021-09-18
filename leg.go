package main

import (
	"github.com/ByteArena/box2d"
)

const legThickness = 0.1

var legVerts = []box2d.B2Vec2{
	{-0.5, -legThickness / 2},
	{0.5, -legThickness / 2},
	{0.5, legThickness / 2},
	{-0.5, legThickness / 2},
}

type LegCfg struct {
	Dir Direction
}

type Leg struct {
	PartBase
	cfg   LegCfg
	verts []box2d.B2Vec2
}

func NewLeg(cfg LegCfg) *Leg {
	return &Leg{
		PartBase: PartBase{img: legImg, dir: cfg.Dir},
		cfg:      cfg,
		verts:    Rotate(cfg.Dir.GetAng(), legVerts...),
	}
}

func (l *Leg) GetPos() box2d.B2Vec2 {
	return l.pos
}

func (l *Leg) Construct(ship *Ship, pos box2d.B2Vec2, size box2d.B2Vec2) {
	l.pos = pos
	pos.OperatorPlusInplace(box2d.B2Vec2MulScalar(0.5, size).OperatorNegate())
	pos.OperatorPlusInplace(box2d.MakeB2Vec2(0.5, 0.5))
	verts := Translate(pos, l.verts...)

	shape := box2d.MakeB2PolygonShape()
	shape.Set(verts, len(verts))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter = box2d.MakeB2Filter()
	fd.Shape = &shape
	fd.Density = FixtureDensity
	fd.Restitution = FixtureRestitution
	ship.body.CreateFixtureFromDef(&fd)
	l.ship = ship
}
