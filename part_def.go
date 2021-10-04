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

// This is for serialization
type OneOfPart struct {
	Cabin        *CabinDef        `yaml:",omitempty"`
	Engine       *EngineDef       `yaml:",omitempty"`
	Crane        *CraneDef        `yaml:",omitempty"`
	Tank         *TankDef         `yaml:",omitempty"`
	Leg          *LegDef          `yaml:",omitempty"`
	LegFastening *LegFasteningDef `yaml:",omitempty"`
}

func NewOneOfPart(d PartDef) *OneOfPart {
	switch pd := d.(type) {
	case CabinDef:
		return &OneOfPart{Cabin: &pd}
	case EngineDef:
		return &OneOfPart{Engine: &pd}
	case CraneDef:
		return &OneOfPart{Crane: &pd}
	case TankDef:
		return &OneOfPart{Tank: &pd}
	case LegDef:
		return &OneOfPart{Leg: &pd}
	case LegFasteningDef:
		return &OneOfPart{LegFastening: &pd}
	default:
		panic("unknown part defenition")
	}
}

func (p OneOfPart) ToPartDef() PartDef {
	if p.Cabin != nil {
		return p.Cabin
	}
	if p.Engine != nil {
		return p.Engine
	}
	if p.Crane != nil {
		return p.Crane
	}
	if p.Tank != nil {
		return p.Tank
	}
	if p.Leg != nil {
		return p.Leg
	}
	if p.LegFastening != nil {
		return p.LegFastening
	}
	return nil
}

type OneOfParts [][]*OneOfPart

func (ps OneOfParts) ToPartDefs() [][]PartDef {
	result := make([][]PartDef, len(ps))
	for y, row := range ps {
		result[y] = make([]PartDef, len(row))
		for x, op := range row {
			if op != nil {
				result[y][x] = op.ToPartDef()
			} else {
				result[y][x] = nil
			}
		}
	}
	return result
}
