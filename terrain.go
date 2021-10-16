package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Terrain struct {
	body *box2d.B2Body
	vSet [][]box2d.B2Vec2
}

func NewTerrain(world *box2d.B2World, vSet [][]box2d.B2Vec2) *Terrain {
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(0, 0)
	bd.Type = box2d.B2BodyType.B2_staticBody
	body := world.CreateBody(&bd)

	for _, verts := range vSet {
		shape := box2d.MakeB2ChainShape()
		shape.CreateLoop(verts, len(verts))

		fd := box2d.MakeB2FixtureDef()
		fd.Filter = box2d.MakeB2Filter()
		fd.Shape = &shape
		fd.Friction = DefaultFriction
		fd.Density = DefaultFixtureDensity
		fd.Restitution = DefaultFixtureRestitution
		body.CreateFixtureFromDef(&fd)
	}

	return &Terrain{
		body: body,
		vSet: vSet,
	}
}

func (g *Terrain) Draw(screen *ebiten.Image, cam Cam) {
	for _, vecs := range g.vSet {
		var path vector.Path

		v := cam.Project(vecs[0], box2d.B2Vec2_zero, 0)
		path.MoveTo(float32(v.X), float32(v.Y))
		for i := 1; i < len(vecs); i++ {
			v := cam.Project(vecs[i], box2d.B2Vec2_zero, 0)
			path.LineTo(float32(v.X), float32(v.Y))
		}

		opts := &ebiten.DrawTrianglesOptions{
			FillRule: ebiten.EvenOdd,
		}

		vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vs {
			vs[i].SrcX = 1
			vs[i].SrcY = 1
			vs[i].ColorR = 0.2
			vs[i].ColorG = 0.2
			vs[i].ColorB = 0.2
		}
		screen.DrawTriangles(vs, is, emptySubImage, opts)
	}
}
