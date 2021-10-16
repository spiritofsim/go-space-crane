package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: description
// TODO: rename
// We need it because ebiten allows to run only one game
type GameDecorator struct {
	g ebiten.Game

	cam *Cam
	ps  *ParticleSystem
	bg  Background
}

func NewGameDecorator() *GameDecorator {
	cam := NewCam()
	ps := NewParticleSystem()
	bg := NewBackground()

	return &GameDecorator{
		g: NewLevelSelector(),

		cam: cam,
		ps:  ps,
		bg:  bg,
	}
}

func (gd *GameDecorator) Update() error {
	err := gd.g.Update()
	if selectedLevel, ok := err.(*SelectedLevel); ok {
		world := box2d.MakeB2World(box2d.MakeB2Vec2(0, Gravity))
		level := LoadLevel(&world, gd.ps, selectedLevel.name)
		game := NewGame(&world, gd.cam, gd.bg, gd.ps, level)
		world.SetContactListener(game)
		gd.g = game
	}

	switch err {
	case levelComplete, levelFailed:
		gd.g = NewLevelSelector()
	}
	return nil
}

func (gd *GameDecorator) Draw(screen *ebiten.Image) {
	gd.g.Draw(screen)
}

func (gd *GameDecorator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
