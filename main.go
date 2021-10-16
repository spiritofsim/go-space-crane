package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := NewGameDecorator()
	checkErr(ebiten.RunGame(g))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
