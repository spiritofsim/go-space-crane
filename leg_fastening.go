package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

var legFasteningVerts = []box2d.B2Vec2{
	{0.5, -legThickness / 2},
	{0.5, legThickness / 2},
	{legThickness / 2, 0.5},
	{-legThickness / 2, 0.5},
}

type LegFasteningCfg struct {
	Dir Direction
}

type LegFastening struct {
	PartBase
	cfg   LegFasteningCfg
	verts []box2d.B2Vec2
}

func NewLegFastening(cfg LegFasteningCfg, img *ebiten.Image) *LegFastening {
	return &LegFastening{
		PartBase: PartBase{img: img, dir: cfg.Dir},
		cfg:      cfg,
		verts:    Rotate(cfg.Dir.GetAng(), legFasteningVerts...),
	}
}

func (l *LegFastening) GetPos() box2d.B2Vec2 {
	return l.pos
}

func (l *LegFastening) Construct(ship *Ship, pos box2d.B2Vec2, size box2d.B2Vec2) {
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
