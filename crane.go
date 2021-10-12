package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"time"
)

type CraneDef struct {
	Dir Direction
}

type Crane struct {
	*GameObj
	chain          []*GameObj
	jaws           *CraneJaws
	lastControlled time.Time

	chainElSize box2d.B2Vec2
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
	worldPos := box2d.B2Vec2Add(shipPos, pos)
	worldPos = box2d.B2Vec2Add(worldPos, shipHalfSize.OperatorNegate())
	worldPos = box2d.B2Vec2Add(worldPos, box2d.MakeB2Vec2(0.5, 0.5))

	crane := &Crane{
		GameObj: NewGameObj(
			world,
			craneSprite,
			worldPos,
			d.Dir.GetAng(), 0,
			box2d.B2Vec2_zero,
			DefaultFriction, DefaultFixtureDensity, DefaultFixtureRestitution),
		chainElSize: getShapeSize(chainElSprite.vertsSet[0]),
	}
	crane.GetBody().SetUserData(crane)

	crane.unwind()
	crane.jaws = NewCraneJaws(crane)

	return crane
}

func (c *Crane) Draw(screen *ebiten.Image, cam Cam) {
	c.jaws.Draw(screen, cam)
	c.GameObj.Draw(screen, cam)
	for _, chainEl := range c.chain {
		chainEl.Draw(screen, cam)
	}
}

func (c *Crane) GetBody() *box2d.B2Body {
	return c.body
}

func (c *Crane) Update() {
	c.jaws.Update()

	// TODO: delay to const
	if c.lastControlled.Add(time.Second / 5).After(time.Now()) {
		return
	}
	c.lastControlled = time.Now()

	// TODO: pass keys from game
	keys := inpututil.AppendPressedKeys(nil)
	for _, key := range keys {
		if key == ebiten.KeyQ {
			c.windup()
		}
		if key == ebiten.KeyA {
			c.unwind()
		}
		if key == ebiten.KeyTab {
			fmt.Println("X")
			c.jaws.OpenClose()
		}
	}
}

func (c *Crane) windup() {
	if len(c.chain) <= 1 {
		return
	}

	c.world.DestroyBody(c.chain[0].body)
	c.chain = c.chain[1:]

	f := box2d.MakeB2Vec2(0, 100)
	c.jaws.upper.body.ApplyForce(f, c.jaws.upper.body.GetPosition(), true)
	c.jaws.lower.body.ApplyForce(f, c.jaws.upper.body.GetPosition(), true)

	// TODO: apply additional force jaws
	if len(c.chain) > 0 {
		// TODO: check if previous join destroyed by destroying its body
		// TODO: use part rotation. now it is hardcoded
		c.createChainJoint(c.body, box2d.B2Vec2_zero, c.chain[0].body, box2d.MakeB2Vec2(0, -c.chainElSize.Y/2))
	}
}

func (c *Crane) unwind() {
	// TODO: use angle (see engine)
	pos := box2d.B2Vec2Add(c.body.GetPosition(), box2d.MakeB2Vec2(0, 0.5+c.chainElSize.Y/2))

	// Chain must be massive (see density) to joint work well
	chainEl := NewGameObj(
		c.world,
		chainElSprite,
		pos, 0, 0,
		box2d.B2Vec2_zero,
		DefaultFriction, 100, DefaultFixtureRestitution)

	if len(c.chain) > 0 {
		// TODO: apply additional force jaws
		prevBody := c.chain[0].body
		c.destroyCrainJoints(prevBody)
		c.createChainJoint(prevBody, box2d.MakeB2Vec2(0, -c.chainElSize.Y/2), chainEl.body, box2d.MakeB2Vec2(0, c.chainElSize.Y/2))
	}
	// TODO: use rotation. now its hardcoded
	c.createChainJoint(c.body, box2d.B2Vec2_zero, chainEl.body, box2d.MakeB2Vec2(0, -c.chainElSize.Y/2))

	// TODO: use linked list?
	c.chain = append([]*GameObj{chainEl}, c.chain...)
}

func (c *Crane) createChainJoint(
	bodyA *box2d.B2Body,
	lpA box2d.B2Vec2,
	bodyB *box2d.B2Body,
	lpB box2d.B2Vec2) {

	// TODO: try ropeJoint
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
