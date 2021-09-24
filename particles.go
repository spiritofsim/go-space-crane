package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
)

var partSize = 0.005
var particleVerts = []box2d.B2Vec2{
	box2d.MakeB2Vec2(-partSize, -partSize),
	box2d.MakeB2Vec2(partSize, -partSize),
	box2d.MakeB2Vec2(partSize, partSize),
	box2d.MakeB2Vec2(-partSize, partSize),
}

type Particle interface {
	// TODO: get rid of this methods. leave update only
	IsDead() bool
	IncAge()
	Update(gravity box2d.B2Vec2)
	Pos() box2d.B2Vec2
	Ang() float64
	Destroy()
}

type ParticleBase struct {
	age int
	ttl int
}

func (p *ParticleBase) IsDead() bool {
	return p.age > p.ttl
}

func (p *ParticleBase) IncAge() {
	p.age++
}

// NPParticle is NOT physical particle
type NPParticle struct {
	*ParticleBase
	pos  box2d.B2Vec2
	lvel box2d.B2Vec2
	ang  float64
	avel float64
}

func NewNPParticle(
	ttl int,
	pos box2d.B2Vec2,
	lvel box2d.B2Vec2,
	ang float64,
	avel float64) *NPParticle {
	return &NPParticle{
		ParticleBase: &ParticleBase{
			age: 0,
			ttl: ttl,
		},
		pos:  pos,
		lvel: lvel,
		ang:  ang,
		avel: avel,
	}
}

func (p *NPParticle) Update(gravity box2d.B2Vec2) {
	p.lvel.OperatorPlusInplace(gravity)
	p.pos.OperatorPlusInplace(p.lvel)
	p.ang += p.avel
}

func (p *NPParticle) Pos() box2d.B2Vec2 {
	return p.pos
}

func (p *NPParticle) Ang() float64 {
	return p.ang
}

func (p *NPParticle) Destroy() {
}

// PParticle is physical particle
type PParticle struct {
	*ParticleBase
	world *box2d.B2World
	body  *box2d.B2Body
}

func NewPParticle(
	world *box2d.B2World,
	ttl int,
	pos box2d.B2Vec2,
	lvel box2d.B2Vec2,
	ang float64,
	avel float64) *PParticle {

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(pos.X, pos.Y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.AllowSleep = false
	bd.Angle = ang
	body := world.CreateBody(&bd)
	body.SetAngularVelocity(avel)
	body.SetLinearVelocity(lvel)

	shape := box2d.MakeB2PolygonShape()
	shape.Set(particleVerts, len(particleVerts))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter = box2d.MakeB2Filter()
	fd.Shape = &shape
	fd.Density = DefaultFixtureDensity
	fd.Restitution = 1
	body.CreateFixtureFromDef(&fd)

	return &PParticle{
		ParticleBase: &ParticleBase{
			age: 0,
			ttl: ttl,
		},
		world: world,
		body:  body,
	}
}

func (p *PParticle) Update(gravity box2d.B2Vec2) {
}

func (p *PParticle) Pos() box2d.B2Vec2 {
	return p.body.GetPosition()
}

func (p *PParticle) Ang() float64 {
	return p.body.GetAngle()
}

func (p *PParticle) Destroy() {
	p.world.DestroyBody(p.body)
}

type ParticleSystem struct {
	// if world not set, use not physical particles
	world     *box2d.B2World
	gravity   box2d.B2Vec2
	particles map[Particle]struct{}
}

func NewParticleSystem(world *box2d.B2World, gravity box2d.B2Vec2) *ParticleSystem {
	return &ParticleSystem{
		world:     world,
		gravity:   box2d.B2Vec2MulScalar(0.001, gravity),
		particles: make(map[Particle]struct{}),
	}
}

func (ps *ParticleSystem) Emit(pos box2d.B2Vec2, dir float64, angDisp float64) {
	count := RandInt(5, 10)

	for i := 0; i < count; i++ {
		ang := RandFloat(dir-angDisp/2, dir+angDisp/2)
		avel := RandFloat(-1, 1)
		speed := RandFloat(5, 10)
		ttl := RandInt(10, 50)

		c, s := math.Cos(ang), math.Sin(ang)
		lvel := box2d.MakeB2Vec2(c, s)
		lvel.OperatorScalarMulInplace(speed)

		var p Particle = NewNPParticle(ttl, pos, lvel, ang, avel)
		if ps.world != nil {
			p = NewPParticle(ps.world, ttl, pos, lvel, ang, avel)
		}

		ps.particles[p] = struct{}{}
	}
}

func (ps *ParticleSystem) Update() {
	for p := range ps.particles {
		p.IncAge()
		p.Update(ps.gravity)
		if p.IsDead() {
			p.Destroy()
			delete(ps.particles, p)
		}
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image, cam Cam) {
	for p := range ps.particles {
		clr := color.RGBA{
			R: uint8(RandInt(0xaa, 0xff)),
			G: 0x66,
			B: 0,
			A: uint8(RandInt(0x66, 0xff)),
		}
		drawDebugPolyFromVerts(screen, p.Pos(), p.Ang(), particleVerts, cam, clr)
	}
}
