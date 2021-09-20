package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EngineCfg struct {
	Dir   Direction
	Power float64
	Keys  []ebiten.Key
}

type Engine struct {
	PartBase
	cfg      EngineCfg
	km       map[ebiten.Key]struct{}
	isActive bool
}

func NewEngine(cfg EngineCfg) *Engine {
	km := make(map[ebiten.Key]struct{})
	for _, key := range cfg.Keys {
		km[key] = struct{}{}
	}

	s := &Sprite{
		img:   engineSprite.img,
		verts: Rotate(cfg.Dir.GetAng(), engineSprite.verts...),
	}

	return &Engine{
		PartBase: PartBase{sprite: s, dir: cfg.Dir},
		cfg:      cfg,
		km:       km,
	}
}

func (e *Engine) GetPos() box2d.B2Vec2 {
	return e.pos
}

func (engine *Engine) Update() {
	// TODO: pass keys from game
	keys := inpututil.PressedKeys()

	engine.isActive = false
	if engine.ship.fuel <= 0 {
		return
	}

	keyFound := false
	for _, key := range keys {
		if _, ok := engine.km[key]; ok {
			keyFound = true
			break
		}
	}
	if !keyFound {
		return
	}
	engine.isActive = true

	fAng := engine.ship.body.GetAngle() + engine.cfg.Dir.Negative().GetAng()
	rot := box2d.NewB2RotFromAngle(fAng)
	force := box2d.B2RotVec2Mul(*rot, box2d.MakeB2Vec2(engine.cfg.Power, 0))

	p := engine.ship.body.GetPosition()
	pt := box2d.B2RotVec2Mul(
		*box2d.NewB2RotFromAngle(engine.ship.body.GetAngle()),
		box2d.MakeB2Vec2(float64(engine.GetPos().X)-engine.ship.size.X/2+0.5, float64(engine.GetPos().Y)-engine.ship.size.Y/2+0.5))
	p.OperatorPlusInplace(pt)
	engine.ship.body.ApplyForce(force, p, true)

	engine.ship.fuel -= engine.cfg.Power * EngineFuelConsumption
	if engine.ship.fuel < 0 {
		engine.ship.fuel = 0
	}
}

func (e *Engine) Construct(ship *Ship, pos box2d.B2Vec2, size box2d.B2Vec2) {
	e.pos = pos
	pos.OperatorPlusInplace(box2d.B2Vec2MulScalar(0.5, size).OperatorNegate())
	pos.OperatorPlusInplace(box2d.MakeB2Vec2(0.5, 0.5))
	verts := Translate(pos, e.sprite.verts...)

	shape := box2d.MakeB2PolygonShape()
	shape.Set(verts, len(verts))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter = box2d.MakeB2Filter()
	fd.Shape = &shape
	fd.Density = FixtureDensity
	fd.Restitution = FixtureRestitution
	ship.body.CreateFixtureFromDef(&fd)
	e.ship = ship
}
