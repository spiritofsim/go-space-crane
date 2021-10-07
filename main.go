package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
var face font.Face
var hoodFace font.Face
var radarArrowImg *ebiten.Image
var hoodImg *ebiten.Image

func init() {
	rand.Seed(time.Now().UnixNano())

	f, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	checkErr(err)

	const dpi = 72
	face, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	checkErr(err)

	hoodFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    30,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	checkErr(err)

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

	radarArrowImg = loadImage("hood/radar_arrow.png")
	hoodImg = loadImage("hood/hood.png")
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Space Crane")
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
