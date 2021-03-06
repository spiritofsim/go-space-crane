package main

import (
	"github.com/ByteArena/box2d"
)

type BasicPart struct {
	*GameObj
}

func (p *BasicPart) GetBody() *box2d.B2Body {
	return p.body
}

func (p *BasicPart) Update(Keys) {}

func ConstructBasicPart(
	world *box2d.B2World,
	_ Tank,
	_ *ParticleSystem,
	shipPos box2d.B2Vec2, // TODO: create interface: IHavePosAndSize
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2,
	sprite Sprite,
	ang float64) Part {

	shipHalfSize := box2d.B2Vec2MulScalar(0.5, shipSize)
	worldPos := box2d.B2Vec2Add(shipPos, pos)
	worldPos = box2d.B2Vec2Add(worldPos, shipHalfSize.OperatorNegate())
	worldPos = box2d.B2Vec2Add(worldPos, box2d.MakeB2Vec2(0.5, 0.5))

	// TODO: fix ship pos. now we use ship pos twice
	//pos.OperatorPlusInplace(shipPos)
	//pos.OperatorPlusInplace(shipHalfSize.OperatorNegate())
	//pos.OperatorPlusInplace(box2d.MakeB2Vec2(0.5, 0.5))
	////x := box2d.B2Vec2Add(shipPos, pos)
	//x := pos
	part := &BasicPart{GameObj: NewGameObj(
		world,
		sprite,
		worldPos,
		ang,
		0,
		box2d.B2Vec2_zero,
		DefaultFriction, DefaultFixtureDensity, DefaultFixtureRestitution, true)}
	part.GetBody().SetUserData(part)
	return part
}

type CabinDef struct {
	Dir Direction
}

func (d CabinDef) Construct(
	world *box2d.B2World,
	tanker Tank,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {

	return ConstructBasicPart(
		world,
		tanker,
		ps,
		shipPos,
		shipSize,
		pos,
		cabinSprite,
		d.Dir.GetAng())
}

type BoxDef struct {
}

func (d BoxDef) Construct(
	world *box2d.B2World,
	tanker Tank,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {

	return ConstructBasicPart(
		world,
		tanker,
		ps,
		shipPos,
		shipSize,
		pos,
		boxSprite,
		0)
}

type LegDef struct {
	Dir Direction
}

func (d LegDef) Construct(
	world *box2d.B2World,
	tanker Tank,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2, // TODO: remove. use ship
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {

	return ConstructBasicPart(
		world,
		tanker,
		ps,
		shipPos,
		shipSize,
		pos,
		legSprite,
		d.Dir.GetAng())
}

type LegFasteningDef struct {
	Dir Direction
}

func (d LegFasteningDef) Construct(
	world *box2d.B2World,
	tanker Tank,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {

	return ConstructBasicPart(
		world,
		tanker,
		ps, shipPos,
		shipSize,
		pos,
		legFasteningSprite,
		d.Dir.GetAng())
}
