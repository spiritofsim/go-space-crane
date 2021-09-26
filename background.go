package main

import "github.com/hajimehoshi/ebiten/v2"

type Background struct {
	img *ebiten.Image
}

func NewBackground() Background {
	return Background{
		img: loadImage("blueprint.png"),
	}
}

func (b *Background) Draw(screen *ebiten.Image, cam Cam) {
	bounds := b.img.Bounds()

	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(b.img, opts)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(bounds.Max.X), 0)
	screen.DrawImage(b.img, opts)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, float64(bounds.Max.Y))
	screen.DrawImage(b.img, opts)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(bounds.Max.X), float64(bounds.Max.Y))
	screen.DrawImage(b.img, opts)

}
