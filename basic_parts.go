package main

import "github.com/ByteArena/box2d"

type BasicPart struct {
	*GameObj
}

func (p *BasicPart) GetBody() *box2d.B2Body {
	return p.body
}

func (p *BasicPart) Update() {}

func ConstructBasicPart(
	world *box2d.B2World,
	ship *Ship,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2,
	sprite Sprite,
	ang float64) Part {

	shipHalfSize := box2d.B2Vec2MulScalar(0.5, shipSize)
	pos.OperatorPlusInplace(shipPos)
	pos.OperatorPlusInplace(shipHalfSize.OperatorNegate())
	pos.OperatorPlusInplace(box2d.MakeB2Vec2(0.5, 0.5))

	return &BasicPart{GameObj: NewGameObj(
		world,
		sprite,
		box2d.B2Vec2Add(shipPos, pos),
		ang,
		0,
		box2d.B2Vec2_zero)}
}

type CabinDef struct {
	Dir Direction
}

func (d CabinDef) Construct(
	world *box2d.B2World,
	ship *Ship,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {
	return ConstructBasicPart(world, ship, ps, shipPos, shipSize, pos, cabinSprite, d.Dir.GetAng())
}

type TankDef struct {
}

func (d TankDef) Construct(
	world *box2d.B2World,
	ship *Ship,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {
	return ConstructBasicPart(world, ship, ps, shipPos, shipSize, pos, tankSprite, 0)
}

type LegDef struct {
	Dir Direction
}

func (d LegDef) Construct(
	world *box2d.B2World,
	ship *Ship,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2, // TODO: remove. use ship
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {

	return ConstructBasicPart(world, ship, ps, shipPos, shipSize, pos, legSprite, d.Dir.GetAng())
}

type LegFasteningDef struct {
	Dir Direction
}

func (d LegFasteningDef) Construct(
	world *box2d.B2World,
	ship *Ship,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {
	return ConstructBasicPart(world, ship, ps, shipPos, shipSize, pos, legFasteningSprite, d.Dir.GetAng())
}
