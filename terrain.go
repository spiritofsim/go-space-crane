package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"go-space-crane/svg"
	"image/color"
	"path"
)

type Terrain struct {
	bodies []*box2d.B2Body
}

func NewTerrain(world *box2d.B2World, asset string) *Terrain {
	parts, err := svg.Load(path.Join(AssetsDir, asset))
	checkErr(err)

	terrain := &Terrain{}
	bodies := make([]*box2d.B2Body, len(parts))
	for i, verts := range parts {
		shape := box2d.MakeB2ChainShape()
		shape.CreateLoop(verts, len(verts))
		fd := box2d.MakeB2FixtureDef()
		fd.Filter = box2d.MakeB2Filter()
		fd.Shape = &shape
		fd.Density = FixtureDensity
		fd.Restitution = FixtureRestitution

		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(0, 0)
		bd.Type = box2d.B2BodyType.B2_staticBody
		body := world.CreateBody(&bd)
		body.CreateFixtureFromDef(&fd)

		bodies[i] = body
		body.SetUserData(terrain)
	}

	terrain.bodies = bodies
	return terrain
}

func (t *Terrain) Draw(screen *ebiten.Image, cam Cam) {
	for _, body := range t.bodies {
		DrawDebugBody(screen, body, cam, color.White)
	}
}
