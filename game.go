package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

type Game struct {
	world     *box2d.B2World
	cam       *Cam
	ship      *Ship
	terrain   *Terrain
	ps        *ParticleSystem
	platforms []*Platform
	cargos    []*Cargo
}

func NewGame(world *box2d.B2World, cam *Cam, ship *Ship, terrain *Terrain, ps *ParticleSystem, platforms []*Platform, cargos []*Cargo) *Game {
	return &Game{
		world:     world,
		cam:       cam,
		ship:      ship,
		terrain:   terrain,
		ps:        ps,
		platforms: platforms,
		cargos:    cargos,
	}
}

func (g *Game) Update() error {
	g.cam.Pos = g.ship.body.GetPosition()
	g.cam.Zoom = MaxCamZoom - g.ship.GetVelocity()*10
	//g.cam.Ang = -g.ship.body.GetAngle()

	if g.cam.Zoom <= MinCamZoom {
		g.cam.Zoom = MinCamZoom
	}
	if g.cam.Zoom > MaxCamZoom {
		g.cam.Zoom = MaxCamZoom
	}

	g.ps.Update()
	g.ship.Update()

	g.world.Step(1.0/60.0, 8, 3)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ship.Draw(screen, *g.cam)
	g.terrain.Draw(screen, *g.cam)
	for _, platform := range g.platforms {
		platform.Draw(screen, *g.cam)
	}
	for _, cargo := range g.cargos {
		cargo.Draw(screen, *g.cam)
	}
	g.ps.Draw(screen, *g.cam)

	g.drawHood(screen)

	if Debug {
		g.drawDebugBodies(screen)
	}
}

func (g *Game) drawDebugBodies(screen *ebiten.Image) {
	for body := g.world.GetBodyList(); body != nil; body = body.GetNext() {
		DrawDebugBody(screen, body, *g.cam, color.White)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) drawHood(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			"Fuel: %0.2f/%02.f\nEnergy: %0.2f",
			g.ship.GetFuel(), g.ship.GetMaxFuel(),
			g.ship.energy),
		10, 10)
}

func (g *Game) resolveContact(ct ContactType, contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	if ct == ContactTypeBegin && !contact.IsTouching() {
		return
	}

	a := contact.GetFixtureA().GetBody().GetUserData()
	b := contact.GetFixtureB().GetBody().GetUserData()

	if ship, ok := a.(*Ship); ok {
		g.ShipContact(ct, contact, impulse, ship, b)
		return
	}
	if ship, ok := b.(*Ship); ok {
		g.ShipContact(ct, contact, impulse, ship, a)
		return
	}
}

func (g *Game) BeginContact(contact box2d.B2ContactInterface) {
	g.resolveContact(ContactTypeBegin, contact, nil)
}

func (g *Game) EndContact(contact box2d.B2ContactInterface) {
	g.resolveContact(ContactTypeEnd, contact, nil)
}

func (g *Game) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
	g.resolveContact(ContactTypePreSolve, contact, nil)
}

func (g *Game) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	g.resolveContact(ContactTypePostSolve, contact, impulse)
}

func (g *Game) ShipContact(
	ct ContactType,
	contact box2d.B2ContactInterface,
	impulse *box2d.B2ContactImpulse,
	ship *Ship,
	other interface{}) {

	if ct == ContactTypePostSolve {
		imp := impulse.NormalImpulses[0]
		if imp > ShipImpulseThreshold {
			ship.energy -= imp
		}
	}

	switch obj := other.(type) {
	case *Platform:
		g.ShipPlatformContact(ct, ship, obj)
	default:
		//fmt.Printf("unknown body %v\n", obj)
	}
}

func (g *Game) ShipPlatformContact(
	ct ContactType,
	ship *Ship,
	platform *Platform) {

	switch ct {
	case ContactTypeBegin:
		// ship must be aligned to refuel
		ship.currentPlatform = platform
		platform.ship = ship
	case ContactTypeEnd:
		ship.currentPlatform = nil
		platform.ship = nil
	}
}

type ContactType string

const (
	ContactTypeBegin     = "BeginContact"
	ContactTypeEnd       = "EndContact"
	ContactTypePreSolve  = "PreSolve"
	ContactTypePostSolve = "PostSolve"
)
