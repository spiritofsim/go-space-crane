package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type CraneJaws struct {
	motor        *box2d.B2RevoluteJoint
	upper, lower *GameObj
}

func NewCraneJaws(c *Crane) *CraneJaws {
	upper := NewGameObj(
		c.world,
		craneUpperJawSprite,
		box2d.B2Vec2Add(c.GetPos(), box2d.MakeB2Vec2(0, 0.5)), //TODO: 4debug
		DirectionDown.GetAng(), 0,
		box2d.B2Vec2_zero,
		1)
	lower := NewGameObj(
		c.world,
		craneLowerJawSprite,
		box2d.B2Vec2Add(c.GetPos(), box2d.MakeB2Vec2(0, 0.5)), //TODO: 4debug
		DirectionDown.GetAng(), 0,
		box2d.B2Vec2_zero,
		1)

	// motor joint
	rjd := box2d.MakeB2RevoluteJointDef()
	rjd.BodyA = upper.body
	rjd.LocalAnchorA = box2d.B2Vec2{-0.5, 0}
	rjd.BodyB = lower.body
	rjd.LocalAnchorB = box2d.B2Vec2{-0.5, 0}
	rjd.CollideConnected = false
	rjd.EnableMotor = true
	rjd.EnableLimit = true
	rjd.UpperAngle = math.Pi / 2
	rjd.LowerAngle = -math.Pi / 4
	rjd.MaxMotorTorque = 100
	m := c.world.CreateJoint(&rjd)

	// joint to last chain element
	cjd := box2d.MakeB2RevoluteJointDef()
	cjd.BodyA = c.chain[len(c.chain)-1].body
	cjd.LocalAnchorA = box2d.B2Vec2{0, c.chainElSize.Y / 2}
	cjd.BodyB = upper.body
	cjd.LocalAnchorB = box2d.B2Vec2{-0.5, 0}
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

func (j *CraneJaws) Open() {
	j.motor.SetMotorSpeed(10)
}

func (j *CraneJaws) Close() {
	j.motor.SetMotorSpeed(-10)
}
