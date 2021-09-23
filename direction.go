package main

import (
	"github.com/ByteArena/box2d"
	"math"
)

type Direction string

const (
	DirectionDown  Direction = "D"
	DirectionLeft  Direction = "L"
	DirectionUp    Direction = "U"
	DirectionRight Direction = "R"
)

func (d Direction) GetAng() float64 {
	switch d {
	case DirectionRight:
		return 0
	case DirectionDown:
		return math.Pi / 2
	case DirectionLeft:
		return math.Pi
	case DirectionUp:
		return math.Pi + math.Pi/2
	}
	panic("bad direction")
}

func (d Direction) GetVec() box2d.B2Vec2 {
	switch d {
	case DirectionRight:
		return box2d.MakeB2Vec2(1, 0)
	case DirectionDown:
		return box2d.MakeB2Vec2(0, 1)
	case DirectionLeft:
		return box2d.MakeB2Vec2(-1, 0)
	case DirectionUp:
		return box2d.MakeB2Vec2(0, -1)
	}
	panic("bad direction")
}

func (d Direction) Negative() Direction {
	switch d {
	case DirectionRight:
		return DirectionLeft
	case DirectionDown:
		return DirectionUp
	case DirectionLeft:
		return DirectionRight
	case DirectionUp:
		return DirectionDown
	}
	panic("bad direction")
}
