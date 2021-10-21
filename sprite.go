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

// TODO: move from sprite
// TODO: rename to ScaleEngine
func (s Sprite) Scale(f float64, dir Direction) Sprite {
	if f == 1 {
		return s
	}

	vertsSet := make([][]box2d.B2Vec2, len(s.vertsSet))
	for y, row := range s.vertsSet {
		vertsSet[y] = make([]box2d.B2Vec2, len(row))
		for x, vec := range row {
			vertsSet[y][x] = box2d.B2Vec2MulScalar(f, vec)
			vertsSet[y][x].X -= 0.5 - f/2
		}
	}

	bounds := s.img.Bounds()
	img := ebiten.NewImage(bounds.Max.X, bounds.Max.Y)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(f, f)
	opts.GeoM.Translate(
		0,
		float64(bounds.Max.Y)/2-float64(bounds.Max.Y)*f/2)

	img.DrawImage(s.img, opts)

	return NewSprite(img, vertsSet)
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
