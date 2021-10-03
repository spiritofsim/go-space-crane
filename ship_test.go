package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestName(t *testing.T) {
	mp := map[string]int{}
	mp["x"] = 10
	mp["y"] = 20

	data, _ := yaml.Marshal(mp)
	t.Log(string(data))
}

func TestMarshalShipDef(t *testing.T) {
	prts := OneOfParts{
		{
			nil,
			nil,
			NewOneOfPart(CabinDef{Dir: DirectionUp}),
			nil,
			nil,
		},
		{
			NewOneOfPart(LegFasteningDef{DirectionRight}),
			NewOneOfPart(TankDef{}),
			NewOneOfPart(CraneDef{Dir: DirectionDown}),
			NewOneOfPart(TankDef{}),
			NewOneOfPart(LegFasteningDef{DirectionDown}),
		},
		{
			NewOneOfPart(LegDef{Dir: DirectionDown}),
			NewOneOfPart(EngineDef{
				Dir:   DirectionDown,
				Power: 100,
				Keys:  []ebiten.Key{ebiten.KeyRight, ebiten.KeyUp},
			}),
			nil,
			NewOneOfPart(EngineDef{
				Dir:   DirectionDown,
				Power: 100,
				Keys:  []ebiten.Key{ebiten.KeyLeft, ebiten.KeyUp},
			}),
			NewOneOfPart(LegDef{Dir: DirectionDown}),
		},
	}
	sd := ShipDef{
		Parts:   prts,
		Energy:  100,
		Fuel:    5000,
		MaxFuel: 5000,
	}

	d, err := yaml.Marshal(sd)
	require.NoError(t, err)
	t.Log(string(d))
}
