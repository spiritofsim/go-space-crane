package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Particle struct {
	img *ebiten.Image

	pos  box2d.B2Vec2
	lvel box2d.B2Vec2
	ang  float64
	avel float64

	age int
	ttl int
}

func NewParticle(
	ttl int,
	pos box2d.B2Vec2,
	lvel box2d.B2Vec2,
	ang float64,
	avel float64) *Particle {

	return &Particle{
		img:  flameParticleSprite.img,
		pos:  pos,
		lvel: lvel,
		ang:  ang,
		avel: avel,
		age:  0,
		ttl:  ttl,
	}
}

func (p *Particle) IsDead() bool {
	return p.age > p.ttl
}

func (p *Particle) IncAge() {
	p.age++
}

func (p *Particle) Update() {
	p.age++
	p.ang += p.avel
	p.pos.OperatorPlusInplace(p.lvel)
}

func (p *Particle) Draw(screen *ebiten.Image, cam Cam) {
	opts := &ebiten.DrawImageOptions{}

	bounds := p.img.Bounds()
	opts.GeoM.Translate(-float64(bounds.Max.X/2), -float64(bounds.Max.Y/2))
	opts.GeoM.Scale(1/float64(bounds.Max.X), 1/float64(bounds.Max.Y))
	opts.GeoM.Rotate(p.ang)
	opts.GeoM.Translate(p.pos.X, p.pos.Y)
	opts.GeoM.Translate(-cam.Pos.X, -cam.Pos.Y)
	opts.GeoM.Scale(cam.Zoom, cam.Zoom)
	opts.GeoM.Rotate(cam.Ang)
	opts.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	screen.DrawImage(p.img, opts)
}

type ParticleSystem struct {
	particles map[*Particle]struct{}
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		particles: make(map[*Particle]struct{}),
	}
}

func (ps *ParticleSystem) Emit(pos box2d.B2Vec2, dir float64, angDisp float64) {
	count := RandInt(1, 50)

	for i := 0; i < count; i++ {
		ang := RandFloat(dir-angDisp/2, dir+angDisp/2)
		avel := RandFloat(-1, 1)
		speed := RandFloat(0.1, 0.2)
		ttl := RandInt(20, 50)

		c, s := math.Cos(ang), math.Sin(ang)
		lvel := box2d.MakeB2Vec2(c, s)
		lvel.OperatorScalarMulInplace(speed)

		rpos := box2d.B2Vec2Add(pos, box2d.B2Vec2{RandFloat(-0.2, 0.2), RandFloat(-0.5, 0.5)})

		p := NewParticle(ttl, rpos, lvel, ang, avel)
		ps.particles[p] = struct{}{}
	}
}

func (ps *ParticleSystem) Update() {
	for p := range ps.particles {
		p.IncAge()
		p.Update()
		if p.IsDead() {
			delete(ps.particles, p)
		}
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image, cam Cam) {
	for particle, _ := range ps.particles {
		particle.Draw(screen, cam)
	}
}
