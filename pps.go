package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type PhysicalParticle struct {
	*GameObj
	age int
	ttl int
}

func NewPhysicalParticle(
	world *box2d.B2World,
	ttl int,
	pos box2d.B2Vec2,
	lvel box2d.B2Vec2,
	ang float64,
	avel float64) *PhysicalParticle {

	obj := NewGameObj(world, flameParticleSprite, pos, ang, avel, lvel, DefaultFriction)
	return &PhysicalParticle{
		GameObj: obj,
		age:     0,
		ttl:     ttl,
	}
}

func (p *PhysicalParticle) IsDead() bool {
	return p.age > p.ttl
}

func (p *PhysicalParticle) IncAge() {
	p.age++
}

func (p *PhysicalParticle) Destroy() {
	p.world.DestroyBody(p.body)
}

func (p *PhysicalParticle) Update() {
	p.age++
}

type PhysicalParticleSystem struct {
	world     *box2d.B2World
	particles map[*PhysicalParticle]struct{}
}

func NewPhysicalParticleSystem(world *box2d.B2World) *PhysicalParticleSystem {
	return &PhysicalParticleSystem{
		world:     world,
		particles: make(map[*PhysicalParticle]struct{}),
	}
}

func (ps *PhysicalParticleSystem) Emit(pos box2d.B2Vec2, dir float64, angDisp float64) {
	count := RandInt(5, 20)

	for i := 0; i < count; i++ {
		ang := RandFloat(dir-angDisp/2, dir+angDisp/2)
		avel := RandFloat(-1, 1)
		speed := RandFloat(5, 10)
		ttl := RandInt(10, 50)

		c, s := math.Cos(ang), math.Sin(ang)
		lvel := box2d.MakeB2Vec2(c, s)
		lvel.OperatorScalarMulInplace(speed)

		p := NewPhysicalParticle(ps.world, ttl, pos, lvel, ang, avel)
		ps.particles[p] = struct{}{}
	}
}

func (ps *PhysicalParticleSystem) Update() {
	for p := range ps.particles {
		p.IncAge()
		p.Update()
		if p.IsDead() {
			p.Destroy()
			delete(ps.particles, p)
		}
	}
}

func (ps *PhysicalParticleSystem) Draw(screen *ebiten.Image, cam Cam) {
	for p := range ps.particles {
		p.Draw(screen, cam)
	}
}
