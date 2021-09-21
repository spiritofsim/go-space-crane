package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
	"time"
)

type CraneCfg struct{}

// TODO: get rid
const chainElLen = 0.5
const chainElTh = 0.1

// TODO: add direction!!!!
// TODO: try to add elements near to ship, not at the end
type Crane struct {
	world *box2d.B2World // We need world here to construct chain elements
	PartBase
	elements []*box2d.B2Body

	lastControlled time.Time

	currentCargo *Cargo
}

func NewCrane(cfg CraneCfg, world *box2d.B2World) *Crane {
	return &Crane{
		world:    world,
		PartBase: PartBase{sprite: crainSprite, dir: DirectionRight},
	}
}

func (c *Crane) GetPos() box2d.B2Vec2 {
	return c.pos
}

func (c *Crane) Draw(screen *ebiten.Image, cam Cam) {
	c.PartBase.Draw(screen, cam)
}

func (c *Crane) Construct(ship *Ship, pos box2d.B2Vec2, size box2d.B2Vec2) {
	c.pos = pos

	tPos := box2d.B2Vec2Add(pos, box2d.B2Vec2MulScalar(0.5, size).OperatorNegate())
	tPos = box2d.B2Vec2Add(tPos, box2d.MakeB2Vec2(0.5, 0.5))
	verts := Translate(tPos, c.sprite.verts...)

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
	// TODO: delay to const
	if c.lastControlled.Add(time.Second / 5).After(time.Now()) {
		return
	}
	c.lastControlled = time.Now()

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

	// crain deployed and no cargo
	if len(c.elements) > 0 && c.currentCargo == nil {
		for body := c.world.GetBodyList(); body != nil; body = body.GetNext() {
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
			djd.BodyA = chainEl
			djd.BodyB = cargo.body
			djd.CollideConnected = true
			c.world.CreateJoint(&djd)
			c.currentCargo = cargo
		}
	}

}

func (c *Crane) windup() {
	if len(c.elements) == 0 {
		return
	}

	x := box2d.MakeB2Vec2(c.pos.X-c.ship.size.X/2+0.5, c.pos.Y-c.ship.size.Y/2+0.5)
	shipAnchor := box2d.B2Vec2Add(x, box2d.MakeB2Vec2(0, chainElLen/2))
	c.world.DestroyBody(c.elements[0])
	c.elements = c.elements[1:]

	if len(c.elements) > 0 {
		// TODO: check if previous join destroyed
		c.createChainJoint(c.world, c.ship.body, shipAnchor, c.elements[0], box2d.MakeB2Vec2(0, -chainElLen/2))
	}
}

// TODO: unwind only on good ship orientation
// TODO: if last segment has a contact or body near, do not unwind
func (c *Crane) unwind() {
	// TODO: to svg
	verts := []box2d.B2Vec2{
		{-chainElTh / 2, -chainElLen / 2},
		{chainElTh / 2, -chainElLen / 2},
		{chainElTh / 2, chainElLen / 2},
		{-chainElTh / 2, chainElLen / 2},
	}

	world := c.ship.body.GetWorld()
	bPos := c.ship.body.GetPosition()

	bd := box2d.MakeB2BodyDef()
	// TODO: use prev body(element) position and angle, not c.pos!!!!!
	x := box2d.MakeB2Vec2(c.pos.X-c.ship.size.X/2+0.5, c.pos.Y-c.ship.size.Y/2+0.5)

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

	if len(c.elements) > 0 {
		// TODO: apply force to prev body if not ship
		prevBody := c.elements[0]
		c.destroyShipJoints(world, prevBody)
		c.createChainJoint(world, prevBody, box2d.MakeB2Vec2(0, -chainElLen/2), elBody, box2d.MakeB2Vec2(0, chainElLen/2))
		c.createChainJoint(world, c.ship.body, box2d.B2Vec2Add(x, box2d.MakeB2Vec2(0, chainElLen/2)), elBody, box2d.MakeB2Vec2(0, -chainElLen/2))
	} else {
		c.createChainJoint(world, c.ship.body, box2d.B2Vec2Add(x, box2d.MakeB2Vec2(0, chainElLen/2)), elBody, box2d.MakeB2Vec2(0, -chainElLen/2))
	}

	// TODO: use linked list?
	c.elements = append([]*box2d.B2Body{elBody}, c.elements...)
}

func (c *Crane) createChainJoint(
	world *box2d.B2World,
	bodyA *box2d.B2Body,
	lpA box2d.B2Vec2,
	bodyB *box2d.B2Body,
	lpB box2d.B2Vec2) {

	rjd := box2d.MakeB2RevoluteJointDef()
	rjd.BodyA = bodyA
	rjd.LocalAnchorA = lpA
	rjd.BodyB = bodyB
	rjd.LocalAnchorB = lpB
	rjd.CollideConnected = false
	world.CreateJoint(&rjd)

	//djd := box2d.MakeB2DistanceJointDef()
	//djd.BodyA = bodyA
	//djd.LocalAnchorA = lpA
	//djd.BodyB = bodyB
	//djd.LocalAnchorB = lpB
	//djd.CollideConnected = false
	//djd.Length = chainElLen
	//world.CreateJoint(&djd)
}

func (c *Crane) destroyShipJoints(world *box2d.B2World, body *box2d.B2Body) {
	type IHaveBodyA interface {
		GetBodyA() *box2d.B2Body
	}

	for joint := body.GetJointList(); joint != nil; joint = joint.Next {
		ja, ok := joint.Joint.(IHaveBodyA)
		if ok && ja.GetBodyA() == c.ship.body {
			world.DestroyJoint(joint.Joint)
			continue
		}
	}
}
