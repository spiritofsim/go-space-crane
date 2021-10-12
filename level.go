package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2/text"
	svg2 "go-space-crane/svg"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"strconv"
)

type PlatformDef struct {
	Fuel float64
}

type LevelDef struct {
	Name        string
	Description string
	Terrain     string
	Ship        ShipDef

	// Platforms overrides existing platforms props
	PlatformDefs map[string]PlatformDef `yaml:"platforms"`
	//Cargos map[string] // TODO

	TaskDefs []string `yaml:"tasks"`
}

type Level struct {
	Ship      *Ship
	Terrain   *Terrain
	Platforms map[string]*Platform
	Cargos    map[string]*Cargo
	Tasks     []Task
	boundsMin box2d.B2Vec2
	boundsMax box2d.B2Vec2
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
	platforms := make(map[string]*Platform)
	boundsMin := box2d.B2Vec2_zero
	boundsMax := box2d.B2Vec2_zero
	for _, rect := range svg.Layers[0].Rects {
		switch rect.Title {
		case "platform":
			fuel, err := strconv.Atoi(rect.Description)
			checkErr(err)
			platform := NewPlatform(
				rect.ID,
				world,
				box2d.B2Vec2Add(rect.Pos, box2d.B2Vec2MulScalar(0.5, rect.Size)),
				rect.Size,
				float64(fuel))
			platforms[platform.id] = platform
		case "bounds":
			// TODO: bounds
			boundsMin = box2d.MakeB2Vec2(rect.Pos.X, rect.Pos.Y)
			boundsMax = box2d.MakeB2Vec2(rect.Pos.X+rect.Size.X, rect.Pos.Y+rect.Size.Y)
		}
	}

	// Cargos are ellipses with "cargo" title
	// Ship default position stored in ellipse with "ship" title
	cargos := make(map[string]*Cargo)
	var shipPos box2d.B2Vec2
	for _, ellipse := range svg.Layers[0].Ellipses {
		switch ellipse.Title {
		case "cargo":
			cargo := NewCargo(ellipse.ID, world, ellipse.Pos)
			cargos[cargo.id] = cargo
		case "ship":
			shipPos = ellipse.Pos
		}
	}

	if levelDef.Ship.Pos != nil {
		shipPos = *levelDef.Ship.Pos
	}
	ship := LoadShip(world, shipPos, ps, levelDef.Ship)

	// TODO: platforms, cargos to map
	for id, pd := range levelDef.PlatformDefs {
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

	tasks := make([]Task, len(levelDef.TaskDefs))
	for i, def := range levelDef.TaskDefs {
		tasks[i] = ParseTaskDef(def, platforms, cargos)
		text.CacheGlyphs(hoodFace, tasks[i].TargetName())
	}

	for _, cargo := range cargos {
		cargo.tasks = tasks
	}

	return Level{
		Ship:      ship,
		Terrain:   terrain,
		Platforms: platforms,
		Cargos:    cargos,
		Tasks:     tasks,
		boundsMin: boundsMin,
		boundsMax: boundsMax,
	}
}

func LoadShip(world *box2d.B2World, pos box2d.B2Vec2, ps *ParticleSystem, def ShipDef) *Ship {
	data, err := ioutil.ReadFile(path.Join(ShipsDir, def.Name+".yaml"))
	checkErr(err)

	var parts OneOfParts
	checkErr(yaml.Unmarshal(data, &parts))

	return NewShip(world, ps, pos, def, parts)
}
