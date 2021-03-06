package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2/text"
	svg2 "go-space-crane/svg"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
	"strconv"
)

type LevelDef struct {
	Name        string
	Description string
	Terrain     string
	Ship        ShipDef
	TaskDefs    []string `yaml:"tasks"`
}

// TODO: get rid of level yaml. Move it to SVG!
type Level struct {
	Ship      *Ship
	Terrain   *Terrain
	Platforms map[string]*Platform
	Cargos    map[string]*Cargo
	Tasks     []Task
	bounds    box2d.B2AABB
}

func LoadLevel(world *box2d.B2World, ps *ParticleSystem, name string) Level {
	levelDefData, err := ioutil.ReadFile(path.Join(LevelsDir, name))
	checkErr(err)
	var levelDef LevelDef
	checkErr(yaml.Unmarshal(levelDefData, &levelDef))

	// this svg holds terrain data, ship position, cargos position, platforms data
	svg, err := svg2.Load(path.Join(AssetsDir, TerrainsDir, levelDef.Terrain+".svg"))
	checkErr(err)

	// TODO: also load rects
	// For now all terrain data stored in pathes
	terrainVertsSet := make([][]box2d.B2Vec2, len(svg.Layers[0].Pathes))
	for i, path := range svg.Layers[0].Pathes {
		terrainVertsSet[i] = path.Verts
	}

	// Platforms are rects with "platform" title
	platforms := make(map[string]*Platform)
	levelBounds := box2d.MakeB2AABB()
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
			levelBounds.LowerBound = box2d.MakeB2Vec2(rect.Pos.X, rect.Pos.Y)
			levelBounds.UpperBound = box2d.MakeB2Vec2(rect.Pos.X+rect.Size.X, rect.Pos.Y+rect.Size.Y)
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
		Terrain:   NewTerrain(world, terrainVertsSet),
		Platforms: platforms,
		Cargos:    cargos,
		Tasks:     tasks,
		bounds:    levelBounds,
	}
}
