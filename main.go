package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

var tankSprite Sprite
var engineSprite Sprite
var legSprite Sprite
var legFasteningSprite Sprite
var cabinSprite Sprite
var craneSprite Sprite
var craneUpperJawSprite Sprite
var craneLowerJawSprite Sprite
var chainElSprite Sprite

func init() {
	rand.Seed(time.Now().UnixNano())
	tankSprite = LoadSpriteObj("tank")
	engineSprite = LoadSpriteObj("engine")
	legSprite = LoadSpriteObj("leg")
	legFasteningSprite = LoadSpriteObj("leg_fastening")
	cabinSprite = LoadSpriteObj("cabin")
	craneSprite = LoadSpriteObj("crane")
	craneUpperJawSprite = LoadSpriteObj("crane_upper_jaw")
	craneLowerJawSprite = LoadSpriteObj("crane_lower_jaw")
	chainElSprite = LoadSpriteObj("chain_el")
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Lander")
	ebiten.SetWindowResizable(true)

	gravity := box2d.MakeB2Vec2(0, Gravity)
	world := box2d.MakeB2World(gravity)

	cam := NewCam()
	particles := NewParticleSystem(&world, gravity)

	ship := NewShip(&world, box2d.MakeB2Vec2(6, 80), PartDefs{
		{
			nil,
			nil,
			CabinDef{Dir: DirectionUp},
			nil,
			nil,
		},
		{
			LegFasteningDef{DirectionRight},
			TankDef{},
			CraneDef{Dir: DirectionDown},
			TankDef{},
			LegFasteningDef{DirectionDown},
		},
		{
			LegDef{Dir: DirectionDown},
			EngineDef{
				Dir:   DirectionDown,
				Power: 100,
				Keys:  []ebiten.Key{ebiten.KeyRight, ebiten.KeyUp},
			},
			nil,
			EngineDef{
				Dir:   DirectionDown,
				Power: 100,
				Keys:  []ebiten.Key{ebiten.KeyLeft, ebiten.KeyUp},
			},
			LegDef{Dir: DirectionDown},
		},
	}, particles, 100, 30000, 30000)

	terrain, platforms, cargos := LoadLevel(&world, "level2")

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
