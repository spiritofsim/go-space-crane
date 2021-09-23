package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	svg2 "go-space-crane/svg"
	"path"
)

// Sprite contains sprite image and shape vertices
type Sprite struct {
	img   *ebiten.Image
	verts []box2d.B2Vec2
}

func LoadSpriteObj(name string) Sprite {
	img := loadImage(name + ".png")
	svg, err := svg2.Load(path.Join(AssetsDir, name+".svg"))
	checkErr(err)

	return Sprite{
		img:   img,
		verts: svg.Layers[0].Pathes[0].Verts,
	}
}
