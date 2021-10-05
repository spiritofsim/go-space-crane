package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	svg2 "go-space-crane/svg"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"strconv"
)

type OneOfCargoDef struct {
	DeleteCargo *DeleteCargoDef
	AddCargo    *AddCargoDef
}

func (c OneOfCargoDef) Apply(world *box2d.B2World, cargos []*Cargo) []*Cargo {
	var cd CargoDef
	switch {
	case c.DeleteCargo != nil:
		cd = c.DeleteCargo
	case c.AddCargo != nil:
		cd = c.AddCargo
	}
	return cd.Apply(world, cargos)
}

type CargoDef interface {
	Apply(world *box2d.B2World, cargos []*Cargo) []*Cargo
}

type DeleteCargoDef struct {
	ID string
}

func (d *DeleteCargoDef) Apply(world *box2d.B2World, cargos []*Cargo) []*Cargo {
	result := make([]*Cargo, 0)
	for _, cargo := range cargos {
		if d.ID != cargo.id {
			result = append(result, cargo)
		}
	}

	if len(result) == len(cargos) {
		checkErr(fmt.Errorf("%v cargo not found", d.ID))
	}
	return result
}

type AddCargoDef struct {
	ID  string
	Pos box2d.B2Vec2
}

func (d *AddCargoDef) Apply(world *box2d.B2World, cargos []*Cargo) []*Cargo {
	return append(cargos, NewCargo(d.ID, world, d.Pos))
}

// -------------
type PlatformDef struct {
	Fuel float64
}

type LevelDef struct {
	Name        string
	Description string
	Terrain     string
	Ship        ShipDef

	// Platforms overrides existing platforms props
	Platforms map[string]PlatformDef
	//Cargos map[string] // TODO

	Tasks  []OneOfTask
	Cargos []OneOfCargoDef
}

func (ld LevelDef) GetTasks() []Task {
	result := make([]Task, len(ld.Tasks))
	for i, task := range ld.Tasks {
		result[i] = task.ToTask()
	}
	return result
}

type Level struct {
	Ship      *Ship
	Terrain   *Terrain
	Platforms []*Platform
	Cargos    []*Cargo
	Tasks     []Task
}

func LoadLevel(world *box2d.B2World, ps *ParticleSystem, name string) Level {
	levelDefData, err := ioutil.ReadFile(path.Join(AssetsDir, LevelsDir, name+".yaml"))
	checkErr(err)
	var levelDef LevelDef
	checkErr(yaml.Unmarshal(levelDefData, &levelDef))

	// this svg holds terrain data, ship position, cargos position, platforms data
	svg, err := svg2.Load(path.Join(AssetsDir, TerrainsDir, levelDef.Terrain+".svg"))
	checkErr(err)

	// TODO: also load rects
	// For now all terrain data stored in pathes
	vertsSet := make([][]box2d.B2Vec2, len(svg.Layers[0].Pathes))
	for i, path := range svg.Layers[0].Pathes {
		vertsSet[i] = path.Verts
	}
	terrainSprite := NewSprite(loadImage(path.Join(TerrainsDir, levelDef.Terrain+".png")), vertsSet)
	terrain := NewTerrain(world, terrainSprite)

	// Platforms are rects with "platform" title
	platforms := make([]*Platform, 0)
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

	// Cargos are ellipses with "cargo" title
	// Ship default position stored in ellipse with "ship" title
	cargos := make([]*Cargo, 0)
	var shipPos box2d.B2Vec2
	for _, ellipse := range svg.Layers[0].Ellipses {
		switch ellipse.Title {
		case "cargo":
			cargo := NewCargo(ellipse.ID, world, ellipse.Pos)
			cargos = append(cargos, cargo)
		case "ship":
			shipPos = ellipse.Pos
		}
	}

	if levelDef.Ship.Pos != nil {
		shipPos = *levelDef.Ship.Pos
	}
	ship := LoadShip(world, shipPos, ps, levelDef.Ship)

	// TODO: platforms, cargos to map
	for id, pd := range levelDef.Platforms {
		found := false
		for _, platform := range platforms {
			if platform.id == id {
				found = true
				platform.fuel = pd.Fuel
			}
		}
		if !found {
			checkErr(fmt.Errorf("%v platform not found in level", id))
		}
	}

	for _, cd := range levelDef.Cargos {
		cargos = cd.Apply(world, cargos)
	}

	return Level{
		Ship:      ship,
		Terrain:   terrain,
		Platforms: platforms,
		Cargos:    cargos,
		Tasks:     levelDef.GetTasks(),
	}
}

func LoadShip(world *box2d.B2World, pos box2d.B2Vec2, ps *ParticleSystem, def ShipDef) *Ship {
	data, err := ioutil.ReadFile(path.Join(ShipsDir, def.Name+".yaml"))
	checkErr(err)

	var parts OneOfParts
	checkErr(yaml.Unmarshal(data, &parts))

	return NewShip(world, ps, pos, def, parts)
}
