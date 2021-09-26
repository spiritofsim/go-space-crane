package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

const (
	// TODO: rename to ShipImpulseDestructionThreshold
	ShipImpulseThreshold = 20
)

type ShipDef struct {
	Parts   OneOfParts
	Energy  float64
	Fuel    float64
	MaxFuel float64
}

// TODO: think about construct ship from multiple bodies and joint them with Weld joint
// With this joints we can destroy ship parts in the future
type Ship struct {
	parts []Part
	size  box2d.B2Vec2
	// This is the angle for the first part
	// We need it to calculate ship orientation. At init it must be 0
	originalAng float64

	ps *ParticleSystem

	energy  float64
	fuel    float64
	maxFuel float64

	// if ship is landed, it is a pointer to platform to refuel ship
	currentPlatform *Platform
}

func NewShip(
	world *box2d.B2World,
	ps *ParticleSystem,
	pos box2d.B2Vec2,
	def ShipDef) *Ship {

	shipSize := box2d.MakeB2Vec2(float64(len(def.Parts[0])), float64(len(def.Parts)))

	ship := &Ship{
		size:    shipSize,
		ps:      ps,
		energy:  def.Energy,
		fuel:    def.Fuel,
		maxFuel: def.MaxFuel,
	}

	parts := make([]Part, 0)
	iparts := make([][]Part, len(def.Parts))
	for y, row := range def.Parts {
		iparts[y] = make([]Part, len(row))
		for x, partDef := range row {
			if partDef != nil {
				part := partDef.ToPartDef().Construct(
					world,
					ship,
					ps,
					pos,
					shipSize,
					box2d.MakeB2Vec2(float64(x), float64(y)))
				parts = append(parts, part)
				iparts[y][x] = part
			}
		}
	}

	// Create Weld joints to upper and left parts
	for y, row := range iparts {
		for x, part := range row {
			if part != nil {
				left := GetLeftPart(iparts, x, y)
				if left != nil {
					jd := box2d.MakeB2WeldJointDef()
					jd.BodyA = part.GetBody()
					jd.BodyB = left.GetBody()

					jd.ReferenceAngle = left.GetAng() - part.GetAng()

					rotA := box2d.NewB2RotFromAngle(math.Pi - part.GetAng())
					jd.LocalAnchorA = box2d.MakeB2Vec2(rotA.C/2, rotA.S/2)

					rotB := box2d.NewB2RotFromAngle(0 - left.GetAng())
					jd.LocalAnchorB = box2d.MakeB2Vec2(rotB.C/2, rotB.S/2)
					world.CreateJoint(&jd)
				}

				upper := GetUpperPart(iparts, x, y)
				if upper != nil {
					jd := box2d.MakeB2WeldJointDef()
					jd.BodyA = part.GetBody()
					jd.BodyB = upper.GetBody()

					jd.ReferenceAngle = upper.GetAng() - part.GetAng()

					rotA := box2d.NewB2RotFromAngle(-math.Pi/2 - part.GetAng())
					jd.LocalAnchorA = box2d.MakeB2Vec2(rotA.C/2, rotA.S/2)
					rotB := box2d.NewB2RotFromAngle(math.Pi/2 - upper.GetAng())
					jd.LocalAnchorB = box2d.MakeB2Vec2(rotB.C/2, rotB.S/2)

					world.CreateJoint(&jd)
				}
			}
		}
	}

	ship.parts = parts
	ship.originalAng = parts[0].GetAng() // TODO: Use cabin for ship center and angle!

	return ship
}

func (s *Ship) ApplyForce(force box2d.B2Vec2) {
	for _, part := range s.parts {
		body := part.GetBody()
		body.ApplyForce(force, body.GetPosition(), true)
	}
}

func (s *Ship) GetPos() box2d.B2Vec2 {
	// TODO:
	return s.parts[0].GetPos()
}

func (s *Ship) GetAng() float64 {
	return s.parts[0].GetAng() - s.originalAng
}

func (s *Ship) GetVelocity() float64 {
	for _, part := range s.parts {
		return part.GetVel()
	}
	panic("ship have no parts")
}

func (s *Ship) Update() {
	for _, part := range s.parts {
		part.Update()
	}

	vel := s.GetVelocity()

	// TODO: land only if ship is in horizontal state
	// TODO: align angle ! it can be negative or > 2*pi
	if s.currentPlatform != nil && s.currentPlatform.fuel > 0 && FloatEquals(vel, 0) {
		if s.fuel < s.maxFuel {
			s.currentPlatform.fuel--
			s.fuel++
		}
	}
}

func (s *Ship) Draw(screen *ebiten.Image, cam Cam) {
	for _, part := range s.parts {
		part.Draw(screen, cam)
	}
}

func GetLeftPart(pds [][]Part, x, y int) Part {
	if x == 0 {
		return nil
	}
	return pds[y][x-1]
}

func GetUpperPart(pds [][]Part, x, y int) Part {
	if y == 0 {
		return nil
	}
	return pds[y-1][x]
}
