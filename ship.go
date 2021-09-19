package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
)

const (
	ShipImpulseThreshold = 20
)

type Ship struct {
	body *box2d.B2Body
	// TODO: maybe store as flat slice?
	parts [][]Part
	size  box2d.B2Vec2

	// TODO: get rid of tanks and engines
	// store fuel in ship
	engines []*Engine
	tanks   []*Tank
	ps      *ParticleSystem

	energy float64

	// if ship is landed, it is a pointer to platform to refuel ship
	currentPlatform *Platform
}

func NewShip(world *box2d.B2World, pos box2d.B2Vec2, parts [][]Part, ps *ParticleSystem, energy float64) *Ship {
	size := box2d.MakeB2Vec2(float64(len(parts[0])), float64(len(parts)))

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(pos.X, pos.Y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.AllowSleep = false
	body := world.CreateBody(&bd)

	engines := make([]*Engine, 0)
	tanks := make([]*Tank, 0)

	ship := &Ship{
		body:    body,
		parts:   parts,
		size:    size,
		engines: engines,
		ps:      ps,
		energy:  energy,
	}
	for y, row := range parts {
		for x, part := range row {
			if part == nil {
				continue
			}
			part.Construct(ship, box2d.MakeB2Vec2(float64(x), float64(y)), size)

			if engine, ok := part.(*Engine); ok {
				engines = append(engines, engine)
			}
			if tank, ok := part.(*Tank); ok {
				tanks = append(tanks, tank)
			}
		}
	}
	ship.engines = engines
	ship.tanks = tanks

	body.SetUserData(ship)
	return ship
}

func (s *Ship) GetFuel() float64 {
	fuel := 0.0
	for _, tank := range s.tanks {
		fuel += tank.Fuel
	}
	return fuel
}

func (s *Ship) GetMaxFuel() float64 {
	fuel := 0.0
	for _, tank := range s.tanks {
		fuel += tank.MaxFuel
	}
	return fuel
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
		for _, tank := range s.tanks {
			if tank.Fuel >= tank.MaxFuel {
				continue
			}

			s.currentPlatform.fuel--
			tank.Fuel++
		}
	}
}

func (s *Ship) Draw(screen *ebiten.Image, cam Cam) {
	DrawDebugBody(screen, s.body, cam, color.White)
	return // TODO: remove
	for _, row := range s.parts {
		for _, part := range row {
			if part == nil {
				continue
			}
			part.Draw(screen, cam)
		}
	}

	for _, engine := range s.engines {
		if !engine.isActive {
			continue
		}

		ev := box2d.B2Vec2MulScalar(0.5, engine.cfg.Dir.GetVec())
		pt := box2d.B2RotVec2Mul(
			*box2d.NewB2RotFromAngle(s.body.GetAngle()),
			box2d.MakeB2Vec2(engine.GetPos().X-s.size.X/2+0.5+ev.X, engine.GetPos().Y-s.size.Y/2+0.5+ev.Y))
		pt = box2d.B2Vec2Add(pt, s.body.GetPosition())
		s.ps.Emit(pt, engine.cfg.Dir.GetAng()+s.body.GetAngle(), math.Pi/2)

	}

}
