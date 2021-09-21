package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

const (
	ShipImpulseThreshold = 20
)

// TODO: think about construct ship from multiple bodies and joint them with Weld joint
// With this joints we can destroy ship parts in the future
type Ship struct {
	box2d.B2RopeJoint

	body *box2d.B2Body
	// TODO: maybe store as flat slice?
	parts [][]Part
	size  box2d.B2Vec2

	ps *ParticleSystem

	energy  float64
	fuel    float64
	maxFuel float64

	// if ship is landed, it is a pointer to platform to refuel ship
	currentPlatform *Platform
}

func NewShip(
	world *box2d.B2World,
	pos box2d.B2Vec2,
	parts [][]Part,
	ps *ParticleSystem,
	energy float64,
	fuel float64,
	maxFuel float64) *Ship {

	size := box2d.MakeB2Vec2(float64(len(parts[0])), float64(len(parts)))

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(pos.X, pos.Y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.AllowSleep = false
	body := world.CreateBody(&bd)

	ship := &Ship{
		body:    body,
		parts:   parts,
		size:    size,
		ps:      ps,
		energy:  energy,
		fuel:    fuel,
		maxFuel: maxFuel,
	}
	for y, row := range parts {
		for x, part := range row {
			if part == nil {
				continue
			}
			part.Construct(ship, box2d.MakeB2Vec2(float64(x), float64(y)), size)
		}
	}

	body.SetUserData(ship)
	return ship
}

func (s *Ship) GetVelocity() float64 {
	lVel := s.body.GetLinearVelocity()
	return math.Sqrt(lVel.X*lVel.X + lVel.Y*lVel.Y)
}

func (s *Ship) Update() {
	for _, row := range s.parts {
		for _, part := range row {
			if part == nil {
				continue
			}
			part.Update()
		}
	}

	ang := s.body.GetAngle()
	vel := s.GetVelocity()

	// TODO: align angle ! it can be negative or > 2*pi
	if s.currentPlatform != nil && s.currentPlatform.fuel > 0 && FloatEquals(ang, 0) && FloatEquals(vel, 0) {
		if s.fuel < s.maxFuel {
			s.currentPlatform.fuel--
			s.fuel++
		}
	}
}

func (s *Ship) Draw(screen *ebiten.Image, cam Cam) {
	for _, row := range s.parts {
		for _, part := range row {
			if part == nil {
				continue
			}
			part.Draw(screen, cam)
		}
	}
}
