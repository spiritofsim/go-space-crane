package main

import (
	"github.com/ByteArena/box2d"
	svg2 "go-space-crane/svg"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"strconv"
)

// All path elements are terrain parts
// All rects with "platform" title are platforms
//   platform fuel stored in description
// All rects with "cargo" title are carcos
func LoadLevel(world *box2d.B2World, name string) (
	shipDef ShipDef,
	shipPos box2d.B2Vec2,
	terrain *Terrain,
	platforms []*Platform,
	cargos []*Cargo) {

	svg, err := svg2.Load(path.Join(AssetsDir, name+".svg"))
	checkErr(err)

	vertsSet := make([][]box2d.B2Vec2, len(svg.Layers[0].Pathes))
	for i, path := range svg.Layers[0].Pathes {
		vertsSet[i] = path.Verts
	}
	sprite := NewSprite(loadImage(name+".png"), vertsSet)
	terrain = NewTerrain(world, sprite)

	platforms = make([]*Platform, 0)
	cargos = make([]*Cargo, 0)
	for _, rect := range svg.Layers[0].Rects {
		switch rect.Title {
		case "platform":
			fuel, err := strconv.Atoi(rect.Description)
			checkErr(err)
			platforms = append(platforms, NewPlatform(
				rect.ID,
				world,
				box2d.B2Vec2Add(rect.Pos, box2d.B2Vec2MulScalar(0.5, rect.Size)),
				rect.Size,
				float64(fuel)))
		}
	}

	for _, ellipse := range svg.Layers[0].Ellipses {
		switch ellipse.Title {
		case "cargo":
			cargo := NewCargo(ellipse.ID, world, ellipse.Pos, box2d.MakeB2Vec2(0.5, 0.5))
			//cargo := NewGameObj(world, cargoSprite, ellipse.Pos, 0, 0, box2d.B2Vec2_zero, DefaultFriction, DefaultFixtureDensity, DefaultFixtureRestitution)
			cargos = append(cargos, cargo)
		case "ship":
			shipPos = ellipse.Pos
			shipName := ellipse.Description
			shipData, err := ioutil.ReadFile(path.Join(AssetsDir, ShipsDir, shipName+".yaml"))
			checkErr(err)
			checkErr(yaml.Unmarshal(shipData, &shipDef))
		}
	}

	return
}
