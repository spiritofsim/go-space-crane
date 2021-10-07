package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
)

type VisitPlatformsTask struct {
	allPlatforms map[string]*Platform
	Platforms    []string
}

func (t *VisitPlatformsTask) Pos() box2d.B2Vec2 {
	if len(t.Platforms) == 0 {
		return box2d.B2Vec2_zero
	}
	return t.allPlatforms[t.Platforms[0]].GetPos()
}

func (t *VisitPlatformsTask) TargetName() string {
	if len(t.Platforms) == 0 {
		return ""
	}
	return fmt.Sprintf("Visit platform %v", t.allPlatforms[t.Platforms[0]].id)
}

func (t *VisitPlatformsTask) IsComplete() bool {
	return len(t.Platforms) == 0
}

func (t *VisitPlatformsTask) ShipLanded(p *Platform) {
	if t.IsComplete() {
		return
	}

	if p.id == t.Platforms[0] {
		t.Platforms = t.Platforms[1:]
	}
}

func (t *VisitPlatformsTask) CargoOnPlatform(p *Platform) {
}
