package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	svg2 "go-space-crane/svg"
	"path"
)

// Sprite contains sprite image and shape vertices
type Sprite struct {
	img      *ebiten.Image
	vertsSet [][]box2d.B2Vec2
}

func NewSprite(img *ebiten.Image, vertsSet [][]box2d.B2Vec2) Sprite {
	return Sprite{img: img, vertsSet: vertsSet}
}

func LoadPart(name string) Sprite {
	name = path.Join(PartsDir, name)
	img := loadImage(name + ".png")
	svg, err := svg2.Load(path.Join(AssetsDir, name+".svg"))
	checkErr(err)

	vSet := make([][]box2d.B2Vec2, len(svg.Layers[0].Pathes))
	for i, path := range svg.Layers[0].Pathes {
		vSet[i] = path.Verts
	}

	return Sprite{
		img:      img,
		vertsSet: vSet,
	}
}
