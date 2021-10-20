package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

const (
	// TODO: 50 is for dummies. Set to 20
	ShipImpulseDestructionThreshold = 30
)

type ShipDef struct {
	Name    string
	Energy  float64
	Fuel    float64
	MaxFuel float64
	Pos     *box2d.B2Vec2
}

type Ship struct {
	world *box2d.B2World
	parts []Part
	size  box2d.B2Vec2
	// This is the angle for the first part
	// We need it to calculate ship orientation. At init it must be 0
	originalAng float64

	energy      float64
	maxEnergy   float64
	fuel        float64
	maxFuel     float64
	isDestroyed bool

	contactPlatform *Platform
}

func NewShip(
	world *box2d.B2World,
	ps *ParticleSystem,
	pos box2d.B2Vec2,
	def ShipDef,
	partDefs [][]PartDef) *Ship {

	shipSize := box2d.MakeB2Vec2(float64(len(partDefs[0])), float64(len(partDefs)))

	ship := &Ship{
		world:     world,
		size:      shipSize,
		energy:    def.Energy,
		maxEnergy: def.Energy,
		fuel:      def.Fuel,
		maxFuel:   def.MaxFuel,
	}

	parts := make([]Part, 0)
	iparts := make([][]Part, len(partDefs))
	for y, row := range partDefs {
		iparts[y] = make([]Part, len(row))
		for x, partDef := range row {
			if partDef != nil {
				part := partDef.Construct(
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
	// TODO: dont need GetLeftPart, GetUpperPart. Can cache in loop
	for y, row := range iparts {
		for x, part := range row {
			if part != nil {
				// TODO a lot of dupes
				// Simplify and move to func
				// TODO: for now joins are not hard enough
				// Try to add joints to upped and lower parts and maybe add some more joints
				// TODO: try to simplify with distance+revolute joint

				if other := GetLeftPart(iparts, x, y); other != nil {
					jd := box2d.MakeB2WeldJointDef()
					jd.CollideConnected = false
					jd.BodyA = part.GetBody()
					jd.BodyB = other.GetBody()
					jd.ReferenceAngle = other.GetAng() - part.GetAng()

					rotA := box2d.NewB2RotFromAngle(math.Pi - part.GetAng())
					jd.LocalAnchorA = box2d.MakeB2Vec2(rotA.C/2, rotA.S/2)
					rotB := box2d.NewB2RotFromAngle(0 - other.GetAng())
					jd.LocalAnchorB = box2d.MakeB2Vec2(rotB.C/2, rotB.S/2)

					world.CreateJoint(&jd)
				}

				if other := GetUpperPart(iparts, x, y); other != nil {
					jd := box2d.MakeB2WeldJointDef()
					jd.CollideConnected = false
					jd.BodyA = part.GetBody()
					jd.BodyB = other.GetBody()
					jd.ReferenceAngle = other.GetAng() - part.GetAng()

					rotA := box2d.NewB2RotFromAngle(-math.Pi/2 - part.GetAng())
					jd.LocalAnchorA = box2d.MakeB2Vec2(rotA.C/2, rotA.S/2)
					rotB := box2d.NewB2RotFromAngle(math.Pi/2 - other.GetAng())
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

func (s *Ship) GetFuel() float64 {
	return s.fuel
}

func (s *Ship) ReduceFuel(val float64) {
	s.fuel -= val
	if s.fuel < 0 {
		s.fuel = 0
	}
}

func (s *Ship) destroy() {
	for _, part := range s.parts {
		for je := part.GetBody().GetJointList(); je != nil; je = je.Next {
			s.world.DestroyJoint(je.Joint)
		}
	}
	s.isDestroyed = true
}

// GetLandedPlatform returns current landed platform or nil if ship is not landed
// Ship is landed when it has contact with platform and zero velocity
// TODO: also check orientation!
func (s *Ship) GetLandedPlatform() *Platform {
	if s.contactPlatform != nil && FloatEquals(s.GetVel(), 0) {
		return s.contactPlatform
	}
	return nil
}

func (s *Ship) ApplyForce(force box2d.B2Vec2) {
	for _, part := range s.parts {
		body := part.GetBody()
		body.ApplyForce(force, body.GetPosition(), true)
	}
}

func (s *Ship) GetPos() box2d.B2Vec2 {
	// TODO: calc cabin pos
	return s.parts[0].GetPos()
}

func (s *Ship) GetAng() float64 {
	return s.parts[0].GetAng() - s.originalAng
}

func (s *Ship) GetVel() float64 {
	for _, part := range s.parts {
		return part.GetVel()
	}
	panic("ship have no parts")
}

func (s *Ship) GetVelVec() box2d.B2Vec2 {
	for _, part := range s.parts {
		return part.GetVelVec()
	}
	panic("ship have no parts")
}

func (s *Ship) Update(keys Keys) {
	for _, part := range s.parts {
		part.Update(keys)
	}

	if platform := s.GetLandedPlatform(); platform != nil {
		// Refueling
		if s.fuel < s.maxFuel && platform.fuel > 0 {
			s.contactPlatform.fuel -= 10
			s.fuel += 10
		}
	}

	// Destroy ship if no energy
	if s.energy <= 0 {
		s.destroy()
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
