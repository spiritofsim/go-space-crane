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
var cargoSprite Sprite
var flameParticleSprite Sprite

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
	cargoSprite = LoadSpriteObj("cargo")
	flameParticleSprite = LoadSpriteObj("flame_particle")
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Lander")
	ebiten.SetWindowResizable(true)

	gravity := box2d.MakeB2Vec2(0, Gravity)
	world := box2d.MakeB2World(gravity)

	cam := NewCam()
	particles := NewParticleSystem(&world, gravity)

	shipDef, shipPos, terrain, platforms, cargos := LoadLevel(&world, "test_level")
	ship := NewShip(&world, particles, shipPos, shipDef)

	bg := NewBackground()
	game := NewGame(&world, cam, ship, terrain, bg, particles, platforms, cargos)
	world.SetContactListener(game)
	err := ebiten.RunGame(game)
	checkErr(err)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
