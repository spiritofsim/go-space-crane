package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
)

type DeliverCargoTask struct {
	Cargo      *Cargo
	Platform   *Platform
	isComplete bool
}

func NewDeliverCargoTask(cargo *Cargo, platform *Platform) *DeliverCargoTask {
	return &DeliverCargoTask{
		Cargo:    cargo,
		Platform: platform,
	}
}

func (t *DeliverCargoTask) Pos() box2d.B2Vec2 {
	// TODO: is cargo not captured, show its position
	return t.Platform.GetPos()
}

func (t *DeliverCargoTask) TargetName() string {
	return fmt.Sprintf("Deliver %v to %v", t.Cargo.id, t.Platform.id)
}

func (t *DeliverCargoTask) IsComplete() bool {
	return t.isComplete
}

func (t *DeliverCargoTask) ShipLanded(p *Platform) {
}

func (t *DeliverCargoTask) CargoOnPlatform(c *Cargo, p *Platform) {
	if c.id == t.Cargo.id && p.id == t.Platform.id {
		t.isComplete = true
	}
}
