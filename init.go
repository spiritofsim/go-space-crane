package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"math/rand"
	"time"
)

var boxSprite Sprite
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
var modalImg *ebiten.Image
var emptyTransparentImage = ebiten.NewImage(1, 1)
var emptyImage = ebiten.NewImage(3, 3)
var emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
var modalTitleFace font.Face
var modalTextFace font.Face
var levelNameFace font.Face
var levelDescFace font.Face

func init() {
	initEbiten()

	emptyImage.Fill(color.White)
	rand.Seed(time.Now().UnixNano())

	f, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	checkErr(err)

	platformFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    30,
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

	modalTitleFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    50,
		DPI:     FontDpi,
		Hinting: font.HintingFull,
	})
	checkErr(err)

	modalTextFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    20,
		DPI:     FontDpi,
		Hinting: font.HintingFull,
	})
	checkErr(err)

	levelNameFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    30,
		DPI:     FontDpi,
		Hinting: font.HintingFull,
	})
	checkErr(err)

	levelDescFace, err = opentype.NewFace(f, &opentype.FaceOptions{
		Size:    30,
		DPI:     FontDpi,
		Hinting: font.HintingFull,
	})
	checkErr(err)

	// Ship parts
	boxSprite = LoadPart("box")
	engineSprite = LoadPart("engine")
	legSprite = LoadPart("leg")
	legFasteningSprite = LoadPart("leg_fastening")
	cabinSprite = LoadPart("cabin")
	craneSprite = LoadPart("crane")
	craneUpperJawSprite = LoadPart("crane_upper_jaw")
	craneLowerJawSprite = LoadPart("crane_lower_jaw")
	chainElSprite = LoadPart("chain_el")

	cargoSprite = LoadPart("cargo")
	flameParticleImg = loadImage("flame_particle.png")

	radarArrowImg = loadImage("hood/radar_arrow.png")

	hoodImg = loadImage("hood/hood.png")
	// Print text on image for future localization
	text.Draw(hoodImg, fmt.Sprintf("%v\n%v", FuelLabelText, EnergyLabelText), hoodFace, 50, 1000, color.White)
	text.Draw(hoodImg, fmt.Sprintf("%v\n%v", TargetLabelText, DistanceLabelText), hoodFace, 550, 1000, color.White)

	// Cache distance glyphs
	for i := 0; i < 50; i++ {
		text.CacheGlyphs(hoodFace, fmt.Sprintf(DistanceText, i))
	}

	modalImg = loadImage("modal.png")
}

func initEbiten() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Space Crane")
	ebiten.SetWindowResizable(true)
}
