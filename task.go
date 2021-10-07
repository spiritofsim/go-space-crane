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
