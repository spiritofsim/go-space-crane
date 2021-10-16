package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

type Game struct {
	world      *box2d.B2World
	cam        *Cam
	ship       *Ship
	terrain    *Terrain
	background Background
	ps         *ParticleSystem
	platforms  map[string]*Platform
	cargos     map[string]*Cargo
	tasks      []Task
	bounds     box2d.B2AABB

	// Optimizations
	prevTargetDistance    int
	prevTargetDistanceImg *ebiten.Image
	prevTargetName        string
	prevTargetNameImg     *ebiten.Image
}

func NewGame(
	world *box2d.B2World,
	cam *Cam,
	background Background,
	ps *ParticleSystem, level Level) *Game {

	return &Game{
		world:      world,
		cam:        cam,
		ship:       level.Ship,
		terrain:    level.Terrain,
		background: background,
		ps:         ps,
		platforms:  level.Platforms,
		cargos:     level.Cargos,
		tasks:      level.Tasks,
		bounds:     level.bounds,
	}
}

func (g *Game) Update() error {
	keys := KeysFromSlice(inpututil.AppendPressedKeys(nil))

	// Tasks
	if len(g.tasks) == 0 {
		// TODO: level complete
		//return nil
	} else {
		if g.tasks[0].IsComplete() {
			g.tasks = g.tasks[1:]
		}

		for _, task := range g.tasks {
			// is landed
			if platform := g.ship.GetLandedPlatform(); platform != nil {
				task.ShipLanded(g.ship.contactPlatform)
			}
		}
	}

	g.cam.Pos = g.ship.GetPos()

	targetZoom := MaxCamZoom - g.ship.GetVel()*20
	if targetZoom <= MinCamZoom {
		targetZoom = MinCamZoom
	}
	if targetZoom > MaxCamZoom {
		targetZoom = MaxCamZoom
	}

	x := targetZoom - g.cam.Zoom
	g.cam.Zoom += x / 100

	g.checkWorldBounds()

	g.ps.Update()
	g.ship.Update(keys)
	for _, cargo := range g.cargos {
		cargo.Update()
	}

	g.world.Step(1.0/60.0, 8, 3)
	return nil
}

// TODO: apply force depends on ship impulse just to stop it
func (g *Game) checkWorldBounds() {
	//force := 50.0
	shipPos := g.ship.GetPos()
	shipVel := g.ship.GetVelVec()
	mult := 10.0

	force := box2d.B2Vec2_zero
	if shipPos.X < g.bounds.LowerBound.X {
		force = box2d.B2Vec2Add(force, box2d.MakeB2Vec2(-shipVel.X*(g.bounds.LowerBound.X-shipPos.X)*mult, 0))
	}
	if shipPos.Y < g.bounds.LowerBound.Y {
		force = box2d.B2Vec2Add(force, box2d.MakeB2Vec2(0, -shipVel.Y*(g.bounds.LowerBound.Y-shipPos.Y)*mult))
	}
	if shipPos.X > g.bounds.UpperBound.X {
		force = box2d.B2Vec2Add(force, box2d.MakeB2Vec2(-shipVel.X*(shipPos.X-g.bounds.UpperBound.X)*mult, 0))
	}
	if shipPos.Y > g.bounds.UpperBound.Y {
		force = box2d.B2Vec2Add(force, box2d.MakeB2Vec2(0, -shipVel.Y*(shipPos.Y-g.bounds.UpperBound.Y)*mult))
	}
	g.ship.ApplyForce(force)

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

	if DrawDebugBodies {
		g.drawDebugBodies(screen)
	}
	if PrintDebugInfo {
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
	screen.DrawImage(hoodImg, nil)

	// Fuel
	func() {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(200, 30)
		opts.GeoM.Translate(200, 974)
		opts.ColorM.Translate(1, 0, 0, 1)
		screen.DrawImage(emptyTransparentImage, opts)

		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(Remap(g.ship.fuel, 0, g.ship.maxFuel, 0, 200), 30)
		opts.GeoM.Translate(200, 974)
		opts.ColorM.Translate(0, 1, 0, 1)
		screen.DrawImage(emptyTransparentImage, opts)
	}()

	// Energy
	func() {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(200, 30)
		opts.GeoM.Translate(200, 1018)
		opts.ColorM.Translate(1, 0, 0, 1)
		screen.DrawImage(emptyTransparentImage, opts)

		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(Remap(g.ship.energy, 0, g.ship.maxEnergy, 0, 200), 30)
		opts.GeoM.Translate(200, 1018)
		opts.ColorM.Translate(0, 1, 0, 1)
		screen.DrawImage(emptyTransparentImage, opts)
	}()

	g.drawRadar(screen)
}

// drawRadar draws radar, pointing to current task object
func (g *Game) drawRadar(screen *ebiten.Image) {
	if len(g.tasks) == 0 {
		return
	}

	ang, dist := GetVecsAng(g.ship.GetPos(), g.tasks[0].Pos())
	iDist := int(dist)
	targetName := g.tasks[0].TargetName()

	if targetName != g.prevTargetName {
		g.prevTargetNameImg = ebiten.NewImage(500, 30)
		txt := targetName
		bounds := text.BoundString(hoodFace, txt)
		text.Draw(g.prevTargetNameImg, txt, hoodFace, -bounds.Min.X, -bounds.Min.Y, color.White)
		g.prevTargetName = targetName
	}
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(750, 974)
	screen.DrawImage(g.prevTargetNameImg, opts)

	if iDist != g.prevTargetDistance {
		g.prevTargetDistanceImg = ebiten.NewImage(500, 30)
		txt := fmt.Sprintf(DistanceText, iDist)
		bounds := text.BoundString(hoodFace, txt)
		text.Draw(g.prevTargetDistanceImg, txt, hoodFace, -bounds.Min.X, -bounds.Min.Y, color.White)
		g.prevTargetDistance = iDist
	}
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(750, 1018)
	screen.DrawImage(g.prevTargetDistanceImg, opts)

	opts = &ebiten.DrawImageOptions{}
	bounds := radarArrowImg.Bounds()
	opts.GeoM.Translate(-float64(bounds.Max.X)/2, -float64(bounds.Max.Y)/2)
	opts.GeoM.Rotate(ang)
	opts.GeoM.Translate(470, 1000)
	screen.DrawImage(radarArrowImg, opts)

}

func (g *Game) resolveContact(ct ContactType, contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
	if ct == ContactTypeBegin && !contact.IsTouching() {
		return
	}

	a := contact.GetFixtureA().GetBody().GetUserData()
	b := contact.GetFixtureB().GetBody().GetUserData()

	if part, ok := a.(Part); ok {
		g.ShipPartContact(ct, contact, impulse, part, b)
		return
	}
	if part, ok := b.(Part); ok {
		g.ShipPartContact(ct, contact, impulse, part, a)
		return
	}

	if cargo, ok := a.(*Cargo); ok {
		g.CargoContact(ct, contact, impulse, cargo, b)
		return
	}
	if cargo, ok := b.(*Cargo); ok {
		g.CargoContact(ct, contact, impulse, cargo, a)
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

func (g *Game) ShipPartContact(
	ct ContactType,
	contact box2d.B2ContactInterface,
	impulse *box2d.B2ContactImpulse,
	part Part,
	other interface{}) {

	if ct == ContactTypePostSolve {
		imp := impulse.NormalImpulses[0]
		if imp > ShipImpulseDestructionThreshold {
			g.ship.energy -= imp
			if g.ship.energy < 0 {
				g.ship.energy = 0
			}
		}
	}

	switch obj := other.(type) {
	case *Platform:
		g.ShipPlatformContact(ct, obj)
	default:
		//fmt.Printf("unknown body %v\n", obj)
	}
}

func (g *Game) ShipPlatformContact(
	ct ContactType,
	platform *Platform) {

	switch ct {
	case ContactTypeBegin:
		g.ship.contactPlatform = platform
		platform.ship = g.ship

	case ContactTypeEnd:
		g.ship.contactPlatform = nil
		platform.ship = nil
	}
}

// ----------------------

func (g *Game) CargoContact(
	ct ContactType,
	contact box2d.B2ContactInterface,
	impulse *box2d.B2ContactImpulse,
	cargo *Cargo,
	other interface{}) {

	switch obj := other.(type) {
	case *Platform:
		g.CargoPlatformContact(ct, cargo, obj)
	default:
		//fmt.Printf("unknown body %v\n", obj)
	}
}

func (g *Game) CargoPlatformContact(
	ct ContactType,
	cargo *Cargo,
	platform *Platform) {

	switch ct {
	case ContactTypeBegin:
		cargo.platform = platform
	case ContactTypeEnd:
		cargo.platform = nil
	}
}

// ----------------------

type ContactType string

const (
	ContactTypeBegin     = "BeginContact"
	ContactTypeEnd       = "EndContact"
	ContactTypePreSolve  = "PreSolve"
	ContactTypePostSolve = "PostSolve"
)
