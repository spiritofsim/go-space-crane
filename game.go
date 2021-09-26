package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

type Game struct {
	world      *box2d.B2World
	cam        *Cam
	ship       *Ship
	terrain    *Terrain
	background Background
	ps         ParticleSystem
	platforms  []*Platform
	cargos     []*GameObj
}

func NewGame(
	world *box2d.B2World,
	cam *Cam,
	ship *Ship,
	terrain *Terrain,
	background Background,
	ps ParticleSystem,
	platforms []*Platform,
	cargos []*GameObj) *Game {

	return &Game{
		world:      world,
		cam:        cam,
		ship:       ship,
		terrain:    terrain,
		background: background,
		ps:         ps,
		platforms:  platforms,
		cargos:     cargos,
	}
}

func (g *Game) Update() error {
	g.cam.Pos = g.ship.GetPos()
	g.cam.Zoom = MaxCamZoom - g.ship.GetVelocity()*10
	//g.cam.Ang = -g.ship.GetAng()

	if g.cam.Zoom <= MinCamZoom {
		g.cam.Zoom = MinCamZoom
	}
	if g.cam.Zoom > MaxCamZoom {
		g.cam.Zoom = MaxCamZoom
	}

	//g.collideWorldBox()

	g.ps.Update()
	g.ship.Update()

	g.world.Step(1.0/60.0, 8, 3)
	return nil
}

// TODO: apply force depends on ship impulse just to stop it
func (g *Game) collideWorldBox() {
	force := 20.0
	shipPos := g.ship.GetPos()
	if shipPos.X < 0 {
		g.ship.ApplyForce(box2d.MakeB2Vec2(force, 0))
	}
	if shipPos.Y < 0 {
		g.ship.ApplyForce(box2d.MakeB2Vec2(0, force))
	}
	if shipPos.X > g.terrain.size.X {
		g.ship.ApplyForce(box2d.MakeB2Vec2(-force, 0))
	}
	if shipPos.Y > g.terrain.size.Y {
		g.ship.ApplyForce(box2d.MakeB2Vec2(0, -force))
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.background.Draw(screen, *g.cam)
	g.ps.Draw(screen, *g.cam)
	g.ship.Draw(screen, *g.cam)

	g.terrain.Draw(screen, *g.cam)

	for _, platform := range g.platforms {
		platform.Draw(screen, *g.cam)
	}
	for _, cargo := range g.cargos {
		cargo.Draw(screen, *g.cam)
	}

	g.drawHood(screen)

	if Debug {
		g.drawDebugBodies(screen)
		g.printDebugInfo(screen)
	}
}

func (g *Game) printDebugInfo(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf(
			"TPS: %0.2f\nFPS: %0.2f\n",
			ebiten.CurrentTPS(),
			ebiten.CurrentFPS(),
		),
		20, 100)

}

func (g *Game) drawDebugBodies(screen *ebiten.Image) {
	clr := color.RGBA{
		R: 0xff,
		G: 0,
		B: 0xff,
		A: 0xff,
	}

	for body := g.world.GetBodyList(); body != nil; body = body.GetNext() {
		DrawDebugBody(screen, body, *g.cam, clr)
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
			g.ship.fuel, g.ship.maxFuel,
			g.ship.energy),
		10, 10)
}

func (g *Game) resolveContact(ct ContactType, contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	if ct == ContactTypeBegin && !contact.IsTouching() {
		return
	}

	a := contact.GetFixtureA().GetBody().GetUserData()
	b := contact.GetFixtureB().GetBody().GetUserData()

	if part, ok := a.(Part); ok {
		g.PartContact(ct, contact, impulse, part, b)
		return
	}
	if part, ok := b.(Part); ok {
		g.PartContact(ct, contact, impulse, part, a)
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

func (g *Game) PartContact(
	ct ContactType,
	contact box2d.B2ContactInterface,
	impulse *box2d.B2ContactImpulse,
	part Part,
	other interface{}) {

	if ct == ContactTypePostSolve {
		imp := impulse.NormalImpulses[0]
		if imp > ShipImpulseDestructionThreshold {
			g.ship.energy -= imp
		}
	}

	switch obj := other.(type) {
	case *Platform:
		g.PlatformContact(ct, obj)
	default:
		//fmt.Printf("unknown body %v\n", obj)
	}
}

func (g *Game) PlatformContact(
	ct ContactType,
	platform *Platform) {

	switch ct {
	case ContactTypeBegin:
		// ship must be aligned to refuel
		g.ship.currentPlatform = platform
		platform.ship = g.ship
	case ContactTypeEnd:
		g.ship.currentPlatform = nil
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
