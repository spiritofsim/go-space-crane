package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"time"
)

type CraneJaws struct {
	motor          *box2d.B2RevoluteJoint
	upper, lower   *GameObj
	lastControlled time.Time
}

func NewCraneJaws(c *Crane) *CraneJaws {
	upper := NewGameObj(
		c.world,
		craneUpperJawSprite,
		box2d.B2Vec2Add(c.GetPos(), box2d.MakeB2Vec2(0, 0.5)),
		DirectionDown.GetAng(), 0,
		box2d.B2Vec2_zero,
		1, 1, 0.0, true)
	lower := NewGameObj(
		c.world,
		craneLowerJawSprite,
		box2d.B2Vec2Add(c.GetPos(), box2d.MakeB2Vec2(0, 0.5)),
		DirectionDown.GetAng(), 0,
		box2d.B2Vec2_zero,
		1, 1, 0.0, true)

	lower.body.SetGravityScale(40)
	upper.body.SetGravityScale(40)

	// motor joint
	rjd := box2d.MakeB2RevoluteJointDef()
	rjd.BodyA = upper.body
	rjd.LocalAnchorA = box2d.MakeB2Vec2(-0.5, 0)
	rjd.BodyB = lower.body
	rjd.LocalAnchorB = box2d.MakeB2Vec2(-0.5, 0)
	rjd.CollideConnected = false
	rjd.EnableMotor = true
	rjd.EnableLimit = true
	rjd.UpperAngle = math.Pi / 2
	rjd.LowerAngle = -math.Pi / 4
	rjd.MaxMotorTorque = 100
	m := c.world.CreateJoint(&rjd)

	// joint to last chain element (upper)
	cjd := box2d.MakeB2RevoluteJointDef()
	cjd.BodyA = c.chain[len(c.chain)-1].body
	cjd.LocalAnchorA = box2d.MakeB2Vec2(0, c.chainElSize.Y/2)
	cjd.BodyB = upper.body
	cjd.LocalAnchorB = box2d.MakeB2Vec2(-0.5, 0)
	cjd.CollideConnected = false
	cjd.EnableMotor = true
	c.world.CreateJoint(&cjd)

	// joint to last chain element (lower)
	cjd = box2d.MakeB2RevoluteJointDef()
	cjd.BodyA = c.chain[len(c.chain)-1].body
	cjd.LocalAnchorA = box2d.MakeB2Vec2(0, c.chainElSize.Y/2)
	cjd.BodyB = lower.body
	cjd.LocalAnchorB = box2d.MakeB2Vec2(-0.5, 0)
	cjd.CollideConnected = false
	cjd.EnableMotor = true
	c.world.CreateJoint(&cjd)

	return &CraneJaws{
		motor: m.(*box2d.B2RevoluteJoint),
		upper: upper,
		lower: lower,
	}
}

func (j *CraneJaws) Draw(screen *ebiten.Image, cam Cam) {
	j.upper.Draw(screen, cam)
	j.lower.Draw(screen, cam)
}

func (j *CraneJaws) Update() {
	if j.lastControlled.Add(time.Second).After(time.Now()) {
		return
	}
	j.motor.SetMotorSpeed(0)
}

func (j *CraneJaws) OpenClose() {
	if j.lastControlled.Add(time.Second / 2).After(time.Now()) {
		return
	}
	j.lastControlled = time.Now()

	ms := j.motor.GetMotorSpeed()
	if ms < 0 {
		j.motor.SetMotorSpeed(10)
		return
	}
	j.motor.SetMotorSpeed(-10)
}
