package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"math/rand"
	"path"
)

func RandFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min) + min
}

const FloatTolerance = 0.01

func FloatEquals(a, b float64) bool {
	delta := math.Abs(a - b)
	if delta < FloatTolerance {
		return true
	}
	return false
}

// TODO: use b2vec
func Remap(val, from1, to1, from2, to2 float64) float64 {
	return (val-from1)/(to1-from1)*(to2-from2) + from2
}

func loadImage(name string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path.Join(AssetsDir, name))
	checkErr(err)
	return img
}

func getShapeSize(verts []box2d.B2Vec2) box2d.B2Vec2 {
	min, max := verts[0], verts[0]
	for i := 1; i < len(verts); i++ {
		if verts[i].X < min.X {
			min.X = verts[i].X
		}
		if verts[i].Y < min.Y {
			min.Y = verts[i].Y
		}
		if verts[i].X > max.X {
			max.X = verts[i].X
		}
		if verts[i].Y > max.Y {
			max.Y = verts[i].Y
		}
	}

	return box2d.B2Vec2{max.X - min.X, max.Y - min.Y}
}

func GetVecsAng(a, b box2d.B2Vec2) (ang float64, dist float64) {
	x := box2d.B2Vec2Add(b, a.OperatorNegate())
	dist = x.Length()
	x.Normalize()
	rot := box2d.B2Rot{
		S: x.Y,
		C: x.X,
	}
	ang = rot.GetAngle()
	return
}
