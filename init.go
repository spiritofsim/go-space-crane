package main

import (
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
var flameParticleImg *ebiten.Image
var platformFace font.Face
var cargoFace font.Face
var hoodFace font.Face
var radarArrowImg *ebiten.Image
var hoodImg *ebiten.Image

func init() {
	rand.Seed(time.Now().UnixNano())

	f, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	checkErr(err)

	platformFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    50,
		DPI:     FontDpi,
		Hinting: font.HintingFull,
	})
	checkErr(err)

	cargoFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    100,
		DPI:     FontDpi,
		Hinting: font.HintingFull,
	})
	checkErr(err)

	hoodFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    30,
		DPI:     FontDpi,
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
	flameParticleImg = loadImage("flame_particle.png")

	radarArrowImg = loadImage("hood/radar_arrow.png")
	hoodImg = loadImage("hood/hood.png")
}
