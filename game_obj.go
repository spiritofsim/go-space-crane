package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type GameObj struct {
	world *box2d.B2World
	body  *box2d.B2Body
	img   *ebiten.Image
}

func NewGameObj(
	world *box2d.B2World,
	sprite Sprite,
	pos box2d.B2Vec2,
	ang float64,
	aVel float64,
	lVel box2d.B2Vec2,
	friction float64,
	density float64,
	restitution float64) *GameObj {

	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(pos.X, pos.Y)
	bd.Angle = ang
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	body := world.CreateBody(&bd)

	for _, verts := range sprite.vertsSet {
		shape := box2d.MakeB2PolygonShape()
		shape.Set(verts, len(verts))

		fd := box2d.MakeB2FixtureDef()
		fd.Filter = box2d.MakeB2Filter()
		fd.Shape = &shape
		fd.Friction = friction
		fd.Density = density
		fd.Restitution = restitution
		body.CreateFixtureFromDef(&fd)
	}

	body.SetAngularVelocity(aVel)
	body.SetLinearVelocity(lVel)

	return &GameObj{
		world: world,
		body:  body,
		img:   sprite.img,
	}
}

func (g *GameObj) GetVel() float64 {
	lVel := g.body.GetLinearVelocity()
	return math.Sqrt(lVel.X*lVel.X + lVel.Y*lVel.Y)
}

func (g *GameObj) GetAng() float64 {
	return g.body.GetAngle()
}

func (g *GameObj) GetPos() box2d.B2Vec2 {
	return g.body.GetPosition()
}

func (g *GameObj) Draw(screen *ebiten.Image, cam Cam) {
	opts := &ebiten.DrawImageOptions{}

	bounds := g.img.Bounds()
	pos := g.body.GetPosition()
	opts.GeoM.Translate(-float64(bounds.Max.X/2), -float64(bounds.Max.Y/2))
	opts.GeoM.Scale(1/float64(bounds.Max.X), 1/float64(bounds.Max.Y))
	opts.GeoM.Rotate(g.body.GetAngle())
	opts.GeoM.Translate(pos.X, pos.Y)
	opts.GeoM.Translate(-cam.Pos.X, -cam.Pos.Y)
	opts.GeoM.Scale(cam.Zoom, cam.Zoom)
	opts.GeoM.Rotate(cam.Ang)
	opts.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)

	screen.DrawImage(g.img, opts)
}
