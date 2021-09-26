package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Particle struct {
	*GameObj
	age int
	ttl int
}

func NewParticle(
	world *box2d.B2World,
	ttl int,
	pos box2d.B2Vec2,
	lvel box2d.B2Vec2,
	ang float64,
	avel float64) *Particle {

	obj := NewGameObj(world, flameParticleSprite, pos, ang, avel, lvel, DefaultFriction)
	return &Particle{
		GameObj: obj,
		age:     0,
		ttl:     ttl,
	}
}

func (p *Particle) IsDead() bool {
	return p.age > p.ttl
}

func (p *Particle) IncAge() {
	p.age++
}

func (p *Particle) Destroy() {
	p.world.DestroyBody(p.body)
}

func (p *Particle) Update() {
	p.age++
}

type ParticleSystem struct {
	world     *box2d.B2World
	gravity   box2d.B2Vec2
	particles map[*Particle]struct{}
}

func NewParticleSystem(world *box2d.B2World, gravity box2d.B2Vec2) *ParticleSystem {
	return &ParticleSystem{
		world:     world,
		gravity:   box2d.B2Vec2MulScalar(0.001, gravity),
		particles: make(map[*Particle]struct{}),
	}
}

func (ps *ParticleSystem) Emit(pos box2d.B2Vec2, dir float64, angDisp float64) {
	count := RandInt(5, 20)

	for i := 0; i < count; i++ {
		ang := RandFloat(dir-angDisp/2, dir+angDisp/2)
		avel := RandFloat(-1, 1)
		speed := RandFloat(5, 10)
		ttl := RandInt(10, 50)

		c, s := math.Cos(ang), math.Sin(ang)
		lvel := box2d.MakeB2Vec2(c, s)
		lvel.OperatorScalarMulInplace(speed)

		p := NewParticle(ps.world, ttl, pos, lvel, ang, avel)
		ps.particles[p] = struct{}{}
	}
}

func (ps *ParticleSystem) Update() {
	for p := range ps.particles {
		p.IncAge()
		p.Update()
		if p.IsDead() {
			p.Destroy()
			delete(ps.particles, p)
		}
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image, cam Cam) {
	for p := range ps.particles {
		p.Draw(screen, cam)
	}
}
