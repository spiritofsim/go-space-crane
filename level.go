package main

import (
	"github.com/ByteArena/box2d"
	svg2 "go-space-crane/svg"
	"gopkg.in/yaml.v3"
	"path"
	"strconv"
)

// All path elements are terrain parts
// All rects with "platform" title are platforms
//   platform fuel stored in description
// All rects with "cargo" title are carcos
func LoadLevel(world *box2d.B2World, name string) (shipDef ShipDef, shipPos box2d.B2Vec2, terrain *Terrain, platforms []*Platform, cargos []*GameObj) {
	svg, err := svg2.Load(path.Join(AssetsDir, name+".svg"))
	checkErr(err)

	vertsSet := make([][]box2d.B2Vec2, len(svg.Layers[0].Pathes))
	for i, path := range svg.Layers[0].Pathes {
		vertsSet[i] = path.Verts
	}
	sprite := NewSprite(loadImage(name+".png"), vertsSet)
	terrain = NewTerrain(world, sprite)

	platforms = make([]*Platform, 0)
	cargos = make([]*GameObj, 0)
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
			cargo := NewGameObj(world, cargoSprite, rect.Pos, 0, 0, box2d.B2Vec2_zero, DefaultFriction)
			cargos = append(cargos, cargo)
		case "ship":
			shipPos = rect.Pos
			checkErr(yaml.Unmarshal([]byte(rect.Description), &shipDef))
		}
	}

	return
}
