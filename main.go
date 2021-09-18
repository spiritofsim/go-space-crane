package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

var tankImg *ebiten.Image
var engineImg *ebiten.Image
var legImg *ebiten.Image
var legFasteningImg *ebiten.Image
var cabinImg *ebiten.Image
var crainImg *ebiten.Image

func init() {
	rand.Seed(time.Now().UnixNano())
	tankImg = loadImage("tank.png")
	engineImg = loadImage("engine.png")
	legImg = loadImage("leg.png")
	legFasteningImg = loadImage("leg_fastening.png")
	cabinImg = loadImage("cabin.png")
	crainImg = loadImage("crane.png")
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
			NewTank(TankCfg{
				Fuel:    5000,
				MaxFuel: 5000,
			}),
			NewTank(TankCfg{
				Fuel:    5000,
				MaxFuel: 5000,
			}),
			NewTank(TankCfg{
				Fuel:    5000,
				MaxFuel: 5000,
			}),
			NewLegFastening(LegFasteningCfg{DirectionDown}),
		},
		{
			NewLeg(LegCfg{DirectionDown}),
			NewEngine(EngineCfg{
				Dir:   DirectionDown,
				Power: 100,
				Keys:  []ebiten.Key{ebiten.KeyRight, ebiten.KeyUp},
			}),
			NewCrain(CraneCfg{}),
			NewEngine(EngineCfg{
				Dir:   DirectionDown,
				Power: 100,
				Keys:  []ebiten.Key{ebiten.KeyLeft, ebiten.KeyUp},
			}),
			NewLeg(LegCfg{DirectionDown}),
		},
	}, particles, 100)

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
