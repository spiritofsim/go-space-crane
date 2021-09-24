package main

import (
	"github.com/ByteArena/box2d"
	svg2 "go-space-crane/svg"
	"path"
	"strconv"
)

// All path elements are terrain parts
// All rects with "platform" title are platforms
//   platform fuel stored in description
// All rects with "cargo" title are carcos
func LoadLevel(world *box2d.B2World, name string) (*Terrain, []*Platform, []*Cargo) {
	svg, err := svg2.Load(path.Join(AssetsDir, name+".svg"))
	checkErr(err)

	terrain := &Terrain{}
	bodies := make([]*box2d.B2Body, len(svg.Layers[0].Pathes))
	for i, path := range svg.Layers[0].Pathes {
		verts := path.Verts
		shape := box2d.MakeB2ChainShape()
		shape.CreateLoop(verts, len(verts))
		fd := box2d.MakeB2FixtureDef()
		fd.Filter = box2d.MakeB2Filter()
		fd.Shape = &shape
		fd.Density = DefaultFixtureDensity
		fd.Restitution = DefaultFixtureRestitution

		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(0, 0)
		bd.Type = box2d.B2BodyType.B2_staticBody
		body := world.CreateBody(&bd)
		body.CreateFixtureFromDef(&fd)

		bodies[i] = body
		body.SetUserData(terrain)
	}
	terrain.bodies = bodies

	platforms := make([]*Platform, 0)
	cargos := make([]*Cargo, 0)
	for _, rect := range svg.Layers[0].Rects {
		switch rect.Title {
		case "platform":
			fuel, err := strconv.Atoi(rect.Description)
			checkErr(err)
			platforms = append(platforms, NewPlatform(
				world,
				box2d.B2Vec2Add(rect.Pos, box2d.B2Vec2MulScalar(0.5, rect.Size)),
				rect.Size,
				float64(fuel)))
		case "cargo":
			cargos = append(cargos, NewCargo(
				world,
				box2d.B2Vec2Add(rect.Pos, box2d.B2Vec2MulScalar(0.5, rect.Size)),
				rect.Size))
		}
	}

	return terrain, platforms, cargos
}
