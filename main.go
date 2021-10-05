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

	// Ship parts
	tankSprite = LoadPart("tank")
	engineSprite = LoadPart("engine")
	legSprite = LoadPart("leg")
	legFasteningSprite = LoadPart("leg_fastening")
	cabinSprite = LoadPart("cabin")
	craneSprite = LoadPart("crane")
	craneUpperJawSprite = LoadPart("crane_upper_jaw")
	craneLowerJawSprite = LoadPart("crane_lower_jaw")
	chainElSprite = LoadPart("chain_el")

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
	ps := NewParticleSystem()

	level := LoadLevel(&world, ps, "level1")

	bg := NewBackground()
	game := NewGame(&world, cam, bg, ps, level)
	world.SetContactListener(game)
	err := ebiten.RunGame(game)
	checkErr(err)

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
