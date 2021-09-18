package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"math"
)

type CraneCfg struct{}

// TODO: add direction!!!!
// TODO: try to add elements near to ship, not at the end
type Crane struct {
	PartBase
	elements []*box2d.B2Body

	currentCargo *Cargo
}

func NewCrain(cfg CraneCfg) *Crane {
	return &Crane{
		PartBase: PartBase{img: crainImg, dir: DirectionRight},
	}
}

func (c *Crane) GetPos() box2d.B2Vec2 {
	return c.pos
}

func (c *Crane) Draw(screen *ebiten.Image, cam Cam) {
	c.PartBase.Draw(screen, cam)
	for _, element := range c.elements {
		DrawDebugBody(screen, element, cam, color.White)
	}
}

func (c *Crane) Construct(ship *Ship, pos box2d.B2Vec2, size box2d.B2Vec2) {
	c.pos = pos
	verts := []box2d.B2Vec2{
		{-0.5, -0.5},
		{0.5, -0.5},
		{0.5, 0},
		{-0.5, 0},
	}

	tPos := box2d.B2Vec2Add(pos, box2d.B2Vec2MulScalar(0.5, size).OperatorNegate())
	tPos = box2d.B2Vec2Add(tPos, box2d.MakeB2Vec2(0.5, 0.5))
	verts = Translate(tPos, verts...)

	shape := box2d.MakeB2PolygonShape()
	shape.Set(verts, len(verts))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter = box2d.MakeB2Filter()
	fd.Shape = &shape
	fd.Density = FixtureDensity
	fd.Restitution = FixtureRestitution
	ship.body.CreateFixtureFromDef(&fd)
	c.ship = ship
}

func (c *Crane) Update() {
	// TODO: pass keys from game
	keys := inpututil.PressedKeys()
	for _, key := range keys {
		if key == ebiten.KeyQ {
			c.windup()
			break
		}
		if key == ebiten.KeyA {
			c.unwind()
			break
		}
	}

	// TODO: get cargos from game
	world := c.ship.body.GetWorld()

	// crain deployed and no cargo
	if len(c.elements) > 0 && c.currentCargo == nil {
		for body := world.GetBodyList(); body != nil; body = body.GetNext() {
			chainEl := c.elements[len(c.elements)-1] //TODO:  c.elements[ len(c.elements)-1] -> getLast

			cargo, ok := body.GetUserData().(*Cargo)
			if !ok {
				continue
			}

			// check distance
			x := chainEl.GetPosition().X - cargo.body.GetPosition().X
			y := chainEl.GetPosition().Y - cargo.body.GetPosition().Y
			dist := math.Sqrt(x*x + y*y)
			if dist > 1 { // TODO: to const
				continue
			}

			djd := box2d.MakeB2DistanceJointDef()
			djd.BodyA = chainEl //TODO:  c.elements[ len(c.elements)-1] -> getLast
			//djd.LocalAnchorA = anchorA // TODO: h:-elsize/2
			djd.BodyB = cargo.body
			//djd.LocalAnchorB = box2d.MakeB2Vec2(0, 0)
			djd.CollideConnected = true
			//djd.Length = elLen
			world.CreateJoint(&djd)
			c.currentCargo = cargo
		}
	}

}

// TODO: slowdown
func (c *Crane) windup() {
	// todo: fix it. well be great to change chain len
	if c.currentCargo != nil {
		return
	}

	if len(c.elements) == 0 {
		return
	}
	world := c.ship.body.GetWorld()

	world.DestroyBody(c.elements[len(c.elements)-1])
	c.elements = c.elements[:len(c.elements)-1]
}

// TODO: slowdown
// TODO: unwind only on good ship orientation
// TODO: if last segment has a contact or body near, do not unwind
func (c *Crane) unwind() {
	// todo: fix it. well be great to change chain len
	if c.currentCargo != nil {
		return
	}

	const elLen = 0.1
	const elTh = 0.1

	verts := []box2d.B2Vec2{
		{-elTh, -elLen / 2},
		{elTh, -elLen / 2},
		{elTh, elLen / 2},
		{-elTh, elLen / 2},
	}

	world := c.ship.body.GetWorld()
	bPos := c.ship.body.GetPosition()

	bd := box2d.MakeB2BodyDef()
	// TODO: use prev body(element) position and angle, not c.pos!!!!!
	x := box2d.MakeB2Vec2(c.pos.X-c.ship.size.X/2+0.5, (c.pos.Y-c.ship.size.Y/2+0.5)+float64(len(c.elements))*elLen+elLen/2)

	elPos := box2d.B2Vec2Add(bPos, x)
	bd.Position.Set(elPos.X, elPos.Y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.AllowSleep = false
	elBody := world.CreateBody(&bd)

	shape := box2d.MakeB2PolygonShape()
	shape.Set(verts, len(verts))
	fd := box2d.MakeB2FixtureDef()
	fd.Filter = box2d.MakeB2Filter()
	fd.Shape = &shape
	fd.Density = FixtureDensity
	fd.Restitution = FixtureRestitution
	elBody.CreateFixtureFromDef(&fd)

	prevBody := c.ship.body
	if len(c.elements) > 0 {
		prevBody = c.elements[len(c.elements)-1]
	}
	anchorA := box2d.B2Vec2Add(x, box2d.MakeB2Vec2(0, -elLen/2))
	if len(c.elements) > 0 {
		anchorA = box2d.MakeB2Vec2(0, elLen/2)
	}

	rjd := box2d.MakeB2RevoluteJointDef()
	rjd.BodyA = prevBody
	rjd.LocalAnchorA = anchorA
	rjd.BodyB = elBody
	rjd.LocalAnchorB = box2d.MakeB2Vec2(0, -elLen/2)
	rjd.CollideConnected = true
	world.CreateJoint(&rjd)

	djd := box2d.MakeB2DistanceJointDef()
	djd.BodyA = prevBody
	//djd.LocalAnchorA = anchorA // TODO: h:-elsize/2
	djd.LocalAnchorA = box2d.B2Vec2Add(anchorA, box2d.MakeB2Vec2(0, -elLen/2))
	djd.BodyB = elBody
	djd.LocalAnchorB = box2d.MakeB2Vec2(0, 0)
	djd.CollideConnected = true
	djd.Length = elLen
	world.CreateJoint(&djd)

	c.elements = append(c.elements, elBody)
}
