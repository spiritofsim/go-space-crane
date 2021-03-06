package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type EngineDef struct {
	Dir   Direction
	Power float64
	Keys  Keys
	Size  float64
}

type Engine struct {
	*GameObj
	tanker   Tank
	ps       *ParticleSystem
	power    float64
	keys     Keys
	isActive bool
}

func (d EngineDef) Construct(
	world *box2d.B2World,
	tanker Tank,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {

	// TODO: duplicate in basic_part
	shipHalfSize := box2d.B2Vec2MulScalar(0.5, shipSize)
	worldPos := box2d.B2Vec2Add(shipPos, pos)
	worldPos = box2d.B2Vec2Add(worldPos, shipHalfSize.OperatorNegate())
	worldPos = box2d.B2Vec2Add(worldPos, box2d.MakeB2Vec2(0.5, 0.5))

	sprite := engineSprite.Scale(d.Size, d.Dir)
	gameObj := NewGameObj(
		world,
		sprite,
		worldPos,
		d.Dir.GetAng(),
		0,
		box2d.B2Vec2_zero,
		DefaultFriction, DefaultFixtureDensity, DefaultFixtureRestitution, true)

	engine := &Engine{
		GameObj: gameObj,
		tanker:  tanker,
		power:   d.Power,
		ps:      ps,
		keys:    d.Keys,
	}
	engine.GetBody().SetUserData(engine)

	return engine
}

func (e *Engine) Draw(screen *ebiten.Image, cam Cam) {
	e.GameObj.Draw(screen, cam)

	if !e.isActive {
		return
	}

	// TODO: fix emit for small engines
	// Flame particles
	pos := box2d.B2Vec2Add(
		e.GetPos(),
		box2d.B2RotVec2Mul(*box2d.NewB2RotFromAngle(e.GetAng()), box2d.MakeB2Vec2(0.5, 0)))
	e.ps.
		Emit(pos, e.GetAng(), math.Pi/4)
}

func (e *Engine) GetBody() *box2d.B2Body {
	return e.body
}

func (e *Engine) Update(keys Keys) {
	e.isActive = false
	if e.tanker.GetFuel() <= 0 {
		return
	}

	// TODO: to func
	keyFound := false
	for key := range keys {
		if e.keys.IsPressed(key) {
			keyFound = true
			break
		}
	}
	if !keyFound {
		return
	}
	e.isActive = true

	rot := box2d.NewB2RotFromAngle(e.GetAng())
	force := box2d.B2RotVec2Mul(*rot, box2d.MakeB2Vec2(e.power, 0))
	force = force.OperatorNegate()
	e.body.ApplyForce(force, e.body.GetPosition(), true)

	e.tanker.ReduceFuel(e.power * EngineFuelConsumption)
}
