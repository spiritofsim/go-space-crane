package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

type Platform struct {
	*GameObj
	id      string
	fuel    float64
	maxFuel float64
	ship    *Ship
	size    box2d.B2Vec2
}

func NewPlatform(id string, world *box2d.B2World, pos box2d.B2Vec2, size box2d.B2Vec2, fuel float64) *Platform {
	verts := []box2d.B2Vec2{
		box2d.MakeB2Vec2(-size.X/2, -size.Y/2),
		box2d.MakeB2Vec2(size.X/2, -size.Y/2),
		box2d.MakeB2Vec2(size.X/2, size.Y/2),
		box2d.MakeB2Vec2(-size.X/2, size.Y/2),
	}

	gObj := NewGameObj(
		world,
		NewSprite(emptyTransparentImage, [][]box2d.B2Vec2{verts}),
		pos,
		0,
		0,
		box2d.B2Vec2_zero,
		DefaultFriction,
		DefaultFixtureDensity,
		DefaultFixtureRestitution, false)

	platform := &Platform{
		GameObj: gObj,
		id:      id,
		fuel:    fuel,
		maxFuel: fuel,
		size:    size,
	}
	platform.body.SetUserData(platform)
	return platform
}

func (p *Platform) Draw(screen *ebiten.Image, cam Cam) {
	p.GameObj.Draw(screen, cam)

	// Draw platform
	func() {
		size := box2d.MakeB2Vec2(p.size.X*cam.Zoom, p.size.Y*cam.Zoom)

		opts := &ebiten.DrawImageOptions{}
		pos := cam.Project(box2d.B2Vec2_zero, p.GetPos(), 0)
		opts.ColorM.Translate(0, 0, 0, 1)
		opts.GeoM.Scale(size.X, size.Y)
		opts.GeoM.Translate(pos.X-(size.X/2), pos.Y-(size.Y/2))
		screen.DrawImage(emptyTransparentImage, opts)

	}()

	// Draw fuel
	func() {
		size := box2d.MakeB2Vec2(p.size.X*cam.Zoom/2, p.size.Y*cam.Zoom/2)
		opts := &ebiten.DrawImageOptions{}
		pos := cam.Project(box2d.B2Vec2_zero, p.GetPos(), 0)
		opts.ColorM.Translate(1, 0, 0, 1)
		opts.GeoM.Scale(size.X, size.Y)
		opts.GeoM.Translate(pos.X-(size.X/2), pos.Y-(size.Y/2))
		screen.DrawImage(emptyTransparentImage, opts)

		opts = &ebiten.DrawImageOptions{}
		opts.ColorM.Translate(0, 1, 0, 1)
		opts.GeoM.Scale(Remap(p.fuel, 0, p.maxFuel, 0, size.X), size.Y)
		opts.GeoM.Translate(pos.X-(size.X/2), pos.Y-(size.Y/2))
		screen.DrawImage(emptyTransparentImage, opts)
	}()

	// Draw platformID
	func() {
		bounds := text.BoundString(platformFace, p.id)
		img := ebiten.NewImage(bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y)

		text.Draw(img, p.id, platformFace, -bounds.Min.X, -bounds.Min.Y, color.White)
		pos := p.GetPos()
		pos = box2d.B2Vec2Add(pos, box2d.MakeB2Vec2(0.5, 0))
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(-float64(bounds.Max.X/2), -float64(bounds.Max.Y/2))
		opts.GeoM.Scale(1/float64(100), 1/float64(100))
		opts.GeoM.Translate(pos.X-p.size.X/2, pos.Y)
		opts.GeoM.Translate(-cam.Pos.X, -cam.Pos.Y)
		opts.GeoM.Scale(cam.Zoom, cam.Zoom)
		opts.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

		//opts := &ebiten.DrawImageOptions{}
		//opts.GeoM.Translate(100,100)

		screen.DrawImage(img, opts)
	}()

}
