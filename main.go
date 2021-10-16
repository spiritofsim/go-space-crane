package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	world := box2d.MakeB2World(box2d.MakeB2Vec2(0, Gravity))

	cam := NewCam()
	ps := NewParticleSystem()
	level := LoadLevel(&world, ps, "level1")
	bg := NewBackground()
	game := NewGame(&world, cam, bg, ps, level)
	world.SetContactListener(game)

	checkErr(ebiten.RunGame(game))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
