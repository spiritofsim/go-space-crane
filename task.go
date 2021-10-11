package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"strings"
)

type Task interface {
	// Pos is the position of task object
	TargetName() string
	Pos() box2d.B2Vec2
	IsComplete() bool
	ShipLanded(p *Platform)
	CargoOnPlatform(c *Cargo, p *Platform)
}

func ParseTaskDef(def string, platforms map[string]*Platform, cargos map[string]*Cargo) Task {
	taskType := def[0:1]
	switch taskType {
	case "v":
		pid := def[2:]
		p, found := platforms[pid]
		if !found {
			checkErr(fmt.Errorf("%v platform not found", pid))
		}
		return NewVisitPlatformTask(p)
	case "d":
		parts := strings.Split(def[2:], "->")
		if len(parts) != 2 {
			checkErr(fmt.Errorf("bad deliver cargo task defenition: %v", def[2:]))
		}
		cid := parts[0]
		pid := parts[1]
		p, found := platforms[pid]
		if !found {
			checkErr(fmt.Errorf("%v platform not found", pid))
		}
		c, found := cargos[cid]
		if !found {
			checkErr(fmt.Errorf("%v cargo not found", cid))
		}
		return NewDeliverCargoTask(c, p)
	default:
		checkErr(fmt.Errorf("unknown task type %v", taskType))
		return nil
	}
}
