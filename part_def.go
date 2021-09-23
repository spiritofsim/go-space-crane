package main

import "github.com/ByteArena/box2d"

type PartDef interface {
	Construct(world *box2d.B2World,
		ship *Ship, //TODO: pass interface
		ps *ParticleSystem,
		shipPos box2d.B2Vec2,
		shipSize box2d.B2Vec2,
		pos box2d.B2Vec2) Part
}

type PartDefs [][]PartDef
