package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

type EngineDef struct {
	Dir   Direction
	Power float64
	Keys  []ebiten.Key
}

type Engine struct {
	*GameObj
	ship     *Ship
	ps       *ParticleSystem
	power    float64
	km       map[ebiten.Key]struct{}
	isActive bool
}

func (d EngineDef) Construct(
	world *box2d.B2World,
	ship *Ship,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {

	// TODO: duplicate in basic_part
	shipHalfSize := box2d.B2Vec2MulScalar(0.5, shipSize)
	worldPos := box2d.B2Vec2Add(shipPos, pos)
	worldPos = box2d.B2Vec2Add(worldPos, shipHalfSize.OperatorNegate())
	worldPos = box2d.B2Vec2Add(worldPos, box2d.MakeB2Vec2(0.5, 0.5))

	km := make(map[ebiten.Key]struct{})
	for _, key := range d.Keys {
		km[key] = struct{}{}
	}

	engine := &Engine{
		GameObj: NewGameObj(
			world,
			engineSprite,
			worldPos,
			d.Dir.GetAng(), 0,
			box2d.B2Vec2_zero,
			DefaultFriction, DefaultFixtureDensity, DefaultFixtureRestitution, true),
		ship:  ship,
		power: d.Power,
		ps:    ps,
		km:    km,
	}
	engine.GetBody().SetUserData(engine)

	return engine
}

func (e *Engine) Draw(screen *ebiten.Image, cam Cam) {
	e.GameObj.Draw(screen, cam)

	if !e.isActive {
		return
	}

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

func (e *Engine) Update() {
	// TODO: pass keys from game
	keys := inpututil.AppendPressedKeys(nil)
	e.isActive = false
	if e.ship.fuel <= 0 {
		return
	}

	// TODO: to func
	keyFound := false
	for _, key := range keys {
		if _, ok := e.km[key]; ok {
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

	// Reduce fuel
	e.ship.fuel -= e.power * EngineFuelConsumption
	if e.ship.fuel < 0 {
		e.ship.fuel = 0
	}
}
