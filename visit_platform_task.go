package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
)

type VisitPlatformTask struct {
	Platform   *Platform
	isComplete bool
}

func NewVisitPlatformTask(platform *Platform) *VisitPlatformTask {
	return &VisitPlatformTask{Platform: platform}
}

func (t *VisitPlatformTask) Pos() box2d.B2Vec2 {
	return t.Platform.GetPos()
}

func (t *VisitPlatformTask) TargetName() string {
	return fmt.Sprintf(VisitPlatformText, t.Platform.id)
}

func (t *VisitPlatformTask) IsComplete() bool {
	return t.isComplete
}

func (t *VisitPlatformTask) ShipLanded(p *Platform) {
	if t.IsComplete() {
		return
	}
	if p.id == t.Platform.id {
		t.isComplete = true
	}
}

func (t *VisitPlatformTask) CargoOnPlatform(c *Cargo, p *Platform) {
}
