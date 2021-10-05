package main

type Task interface {
	IsComplete() bool
	ShipLanded(p *Platform)
	CargoOnPlatform(p *Platform)
}

type OneOfTask struct {
	VisitPlatforms *VisitPlatformsTask `yaml:",omitempty"`
}

func (t OneOfTask) ToTask() Task {
	if t.VisitPlatforms != nil {
		return t.VisitPlatforms
	}
	return nil
}

type VisitPlatformsTask struct {
	Platforms []string
}

func NewVisitPlatformsTask(platformsToVisit ...string) *VisitPlatformsTask {
	return &VisitPlatformsTask{Platforms: platformsToVisit}
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
