package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"time"
)

type CraneDef struct {
	Dir Direction
}

// TODO: move to svg
const chainElLen = 0.2
const chainElTh = 0.1

type Crane struct {
	*GameObj
	elements       []*box2d.B2Body
	lastControlled time.Time
}

func (d CraneDef) Construct(
	world *box2d.B2World,
	ship *Ship,
	ps *ParticleSystem,
	shipPos box2d.B2Vec2,
	shipSize box2d.B2Vec2,
	pos box2d.B2Vec2) Part {

	// TODO: duplicate in basic_part
	shipHalfSize := box2d.B2Vec2MulScalar(0.5, shipSize)
	pos.OperatorPlusInplace(shipPos)
	pos.OperatorPlusInplace(shipHalfSize.OperatorNegate())
	pos.OperatorPlusInplace(box2d.MakeB2Vec2(0.5, 0.5))

	crane := &Crane{
		GameObj: NewGameObj(
			world,
			crainSprite,
			box2d.B2Vec2Add(shipPos, pos),
			d.Dir.GetAng(), 0,
			box2d.B2Vec2_zero),
	}
	crane.GetBody().SetUserData(crane)

	return crane
}

// TODO: try to remove. Draw is already in GameObj
func (c *Crane) Draw(screen *ebiten.Image, cam Cam) {
	c.GameObj.Draw(screen, cam)
}

func (c *Crane) GetBody() *box2d.B2Body {
	return c.body
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
}

func (c *Crane) windup() {
	if len(c.elements) == 0 {
		return
	}

	c.world.DestroyBody(c.elements[0])
	c.elements = c.elements[1:]

	if len(c.elements) > 0 {
		// TODO: check if previous join destroyed
		// TODO: use part rotation. now it is hardcoded
		c.createChainJoint(c.body, box2d.B2Vec2{0.5, 0}, c.elements[0], box2d.MakeB2Vec2(0, -chainElLen/2))
	}
}

func (c *Crane) unwind() {
	// TODO: to svg
	verts := []box2d.B2Vec2{
		{-chainElTh / 2, -chainElLen / 2},
		{chainElTh / 2, -chainElLen / 2},
		{chainElTh / 2, chainElLen / 2},
		{-chainElTh / 2, chainElLen / 2},
	}

	// TODO: use angle (see engine)
	pos := box2d.B2Vec2Add(c.body.GetPosition(), box2d.B2Vec2{0, 0.5 + chainElLen/2})

	bd := box2d.MakeB2BodyDef()
	elPos := pos
	bd.Position.Set(elPos.X, elPos.Y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.AllowSleep = false
	elBody := c.world.CreateBody(&bd)
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
		c.destroyCrainJoints(prevBody)
		c.createChainJoint(prevBody, box2d.MakeB2Vec2(0, -chainElLen/2), elBody, box2d.MakeB2Vec2(0, chainElLen/2))
	}
	// TODO: use rotation. now its hardcoded
	c.createChainJoint(c.body, box2d.B2Vec2{0.5, 0}, elBody, box2d.MakeB2Vec2(0, -chainElLen/2))

	// TODO: use linked list?
	c.elements = append([]*box2d.B2Body{elBody}, c.elements...)
}

func (c *Crane) createChainJoint(
	bodyA *box2d.B2Body,
	lpA box2d.B2Vec2,
	bodyB *box2d.B2Body,
	lpB box2d.B2Vec2) {

	// TODO: try chainJoint
	rjd := box2d.MakeB2RevoluteJointDef()
	rjd.BodyA = bodyA
	rjd.LocalAnchorA = lpA
	rjd.BodyB = bodyB
	rjd.LocalAnchorB = lpB
	rjd.CollideConnected = false
	c.world.CreateJoint(&rjd)

	//djd := box2d.MakeB2DistanceJointDef()
	//djd.BodyA = bodyA
	//djd.LocalAnchorA = lpA
	//djd.BodyB = bodyB
	//djd.LocalAnchorB = lpB
	//djd.CollideConnected = false
	//djd.Length = chainElLen
	//c.world.CreateJoint(&djd)
}

func (c *Crane) destroyCrainJoints(body *box2d.B2Body) {
	type IHaveBodyA interface {
		GetBodyA() *box2d.B2Body
	}

	for joint := body.GetJointList(); joint != nil; joint = joint.Next {
		ja, ok := joint.Joint.(IHaveBodyA)
		if ok && ja.GetBodyA() == c.GetBody() {
			c.world.DestroyJoint(joint.Joint)
			continue
		}
	}
}
