package main

type Task interface {
	IsComplete() bool
	ShipLanded(p *Platform)
	CargoOnPlatform(p *Platform)
}

type ShipPlatformsTask struct {
	PlatformsToVisit []string
}

func NewShipPlatformsTask(platformsToVisit ...string) *ShipPlatformsTask {
	return &ShipPlatformsTask{PlatformsToVisit: platformsToVisit}
}

func (t *ShipPlatformsTask) IsComplete() bool {
	return len(t.PlatformsToVisit) == 0
}

func (t *ShipPlatformsTask) ShipLanded(p *Platform) {
	if t.IsComplete() {
		return
	}

	if p.id == t.PlatformsToVisit[0] {
		t.PlatformsToVisit = t.PlatformsToVisit[1:]
	}
}

func (t *ShipPlatformsTask) CargoOnPlatform(p *Platform) {
}
