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
			NewCrane(CraneCfg{}),
			NewEngine(EngineCfg{
				Dir:   DirectionDown,
				Power: 100,
				Keys:  []ebiten.Key{ebiten.KeyLeft, ebiten.KeyUp},
			}),
			NewLeg(LegCfg{DirectionDown}),
		},
	}, particles, 100, 3000, 3000)

	terrain, platforms, cargos := LoadLevel(&world, "level1")

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
