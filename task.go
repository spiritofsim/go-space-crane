package main

import "github.com/ByteArena/box2d"

type Task interface {
	// Pos is the position of task object
	TargetName() string
	Pos() box2d.B2Vec2
	IsComplete() bool
	ShipLanded(p *Platform)
	CargoOnPlatform(p *Platform)
}

type OneOfTask struct {
	VisitPlatforms *VisitPlatformsTask `yaml:",omitempty"`
}

func (t OneOfTask) ToTask(platforms []*Platform, cargos []*Cargo) Task {
	pm := make(map[string]*Platform)
	for _, platform := range platforms {
		pm[platform.id] = platform
	}

	if t.VisitPlatforms != nil {
		t.VisitPlatforms.allPlatforms = pm
		return t.VisitPlatforms
	}
	return nil
}

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
	return t.allPlatforms[t.Platforms[0]].id
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
