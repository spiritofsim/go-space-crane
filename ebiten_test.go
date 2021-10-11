package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"testing"
)

func BenchmarkEbitenNewImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ebiten.NewImage(500, 30)
	}
}

func BenchmarkClearNewImage(b *testing.B) {
	img := ebiten.NewImage(500, 30)
	for i := 0; i < b.N; i++ {
		img.Clear()
	}
}

func BenchmarkFillNewImage(b *testing.B) {
	img := ebiten.NewImage(500, 30)
	for i := 0; i < b.N; i++ {
		img.Fill(color.White)
	}
}
