package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

var tankSprite *Sprite
var engineSprite *Sprite
var legSprite *Sprite
var legFasteningSprite *Sprite
var cabinSprite *Sprite
var crainSprite *Sprite

func init() {
	rand.Seed(time.Now().UnixNano())
	tankSprite = LoadSpriteObj("tank")
	engineSprite = LoadSpriteObj("engine")
	legSprite = LoadSpriteObj("leg")
	legFasteningSprite = LoadSpriteObj("leg_fastening")
	cabinSprite = LoadSpriteObj("cabin")
	crainSprite = LoadSpriteObj("crane")
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Lander")
	ebiten.SetWindowResizable(true)

	gravity := box2d.MakeB2Vec2(0, Gravity)
	world := box2d.MakeB2World(gravity)

	cam := NewCam()
	particles := NewParticleSystem(&world, gravity)

	ship := NewShip(&world, box2d.MakeB2Vec2(5, 20), [][]Part{
		{
			nil,
			nil,
			NewCabin(CabinCfg{Dir: DirectionUp}),
			nil,
			nil,
		},
		{
			NewLegFastening(LegFasteningCfg{DirectionRight}),
			NewTank(TankCfg{}),
			NewTank(TankCfg{}),
			NewTank(TankCfg{}),
			NewLegFastening(LegFasteningCfg{DirectionDown}),
		},
		{
			NewLeg(LegCfg{DirectionDown}),
			NewEngine(EngineCfg{
				Dir:   DirectionDown,
				Power: 100,
				Keys:  []ebiten.Key{ebiten.KeyRight, ebiten.KeyUp},
			}),
			NewCrane(CraneCfg{}, &world),
			NewEngine(EngineCfg{
				Dir:   DirectionDown,
				Power: 100,
				Keys:  []ebiten.Key{ebiten.KeyLeft, ebiten.KeyUp},
			}),
			NewLeg(LegCfg{DirectionDown}),
		},
	}, particles, 100, 30000, 30000)

	terrain, platforms, cargos := LoadLevel(&world, "level1")

	// TODO: sandbox
	func() {
		return
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(5, 18)
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.AllowSleep = false
		elBody := world.CreateBody(&bd)
		shape := box2d.MakeB2PolygonShape()
		shape.SetAsBox(0.5, 0.5)
		fd := box2d.MakeB2FixtureDef()
		fd.Filter = box2d.MakeB2Filter()
		fd.Shape = &shape
		fd.Density = FixtureDensity
		fd.Restitution = FixtureRestitution
		elBody.CreateFixtureFromDef(&fd)
		elBody.ApplyLinearImpulse(box2d.B2Vec2{1000, 1000}, box2d.B2Vec2{0, 0}, true)

		bd2 := box2d.MakeB2BodyDef()
		bd2.Position.Set(5, 17)
		bd2.Type = box2d.B2BodyType.B2_dynamicBody
		bd2.AllowSleep = false
		elBody2 := world.CreateBody(&bd2)
		shape2 := box2d.MakeB2PolygonShape()
		shape2.SetAsBox(0.5, 0.5)
		fd2 := box2d.MakeB2FixtureDef()
		fd2.Filter = box2d.MakeB2Filter()
		fd2.Shape = &shape
		fd2.Density = FixtureDensity
		fd2.Restitution = FixtureRestitution
		elBody2.CreateFixtureFromDef(&fd2)

		jd := box2d.MakeB2WeldJointDef()
		jd.BodyA = elBody
		jd.LocalAnchorA = box2d.B2Vec2{0, -0.5}
		jd.BodyB = elBody2
		jd.LocalAnchorB = box2d.B2Vec2{0, 0.5}
		jd.CollideConnected = false
		world.CreateJoint(&jd)
	}()

	game := NewGame(&world, cam, ship, terrain, particles, platforms, cargos)
	world.SetContactListener(game)
	err := ebiten.RunGame(game)
	checkErr(err)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
