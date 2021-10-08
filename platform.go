package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

type Platform struct {
	*GameObj
	id   string
	fuel float64
	ship *Ship
	size box2d.B2Vec2
}

func NewPlatform(id string, world *box2d.B2World, pos box2d.B2Vec2, size box2d.B2Vec2, fuel float64) *Platform {
	verts := []box2d.B2Vec2{
		{-size.X / 2, -size.Y / 2},
		{size.X / 2, -size.Y / 2},
		{size.X / 2, size.Y / 2},
		{-size.X / 2, size.Y / 2},
	}
	gobj := NewGameObj(
		world,
		NewSprite(nil, [][]box2d.B2Vec2{verts}),
		pos,
		0,
		0,
		box2d.B2Vec2_zero,
		DefaultFriction,
		DefaultFixtureDensity,
		DefaultFixtureRestitution)

	platform := &Platform{
		GameObj: gobj,
		id:      id,
		fuel:    fuel,
		size:    size,
	}
	platform.body.SetUserData(platform)
	return platform
}

func (p *Platform) Draw(screen *ebiten.Image, cam Cam) {
	return
	// TODO: same as hood
	pos := p.body.GetPosition()
	px := box2d.MakeB2Vec2(-p.size.X/2, 0)
	px = cam.Project(px, pos, 0)

	pSize := box2d.B2Vec2MulScalar(300, p.size)
	txt := fmt.Sprintf("%v : %v", p.id, int(p.fuel))
	textBounds := text.BoundString(platformFace, txt)
	img := ebiten.NewImage(int(pSize.X), int(pSize.Y))
	text.Draw(img, txt, platformFace, int(pSize.X/2)-textBounds.Max.X/2, int(pSize.Y/2)-textBounds.Min.Y/2, color.White)

	// TODO: dupe from gameObj.Draw
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-pSize.X/2, -pSize.Y/2)
	opts.GeoM.Scale(1/float64(pSize.Y), 1/float64(pSize.Y))
	opts.GeoM.Translate(pos.X, pos.Y)
	opts.GeoM.Translate(-cam.Pos.X, -cam.Pos.Y)
	opts.GeoM.Scale(cam.Zoom, cam.Zoom)
	opts.GeoM.Rotate(cam.Ang)
	opts.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	screen.DrawImage(img, opts)
}
