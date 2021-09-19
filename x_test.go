package main

import (
	"encoding/json"
	"io"
	"strings"
	"testing"
)

const j = `[
    [
        null,
        {
            "Cabin": {"Dir": "U"}
        },
        {
            "Engine": {"Dir": "U"}
        },
        null
    ]
]`

func ParseShipParts(reader io.Reader) [][]Part {
	type CP struct {
		Cabin        *CabinCfg
		Crain        *CraneCfg
		Engine       *EngineCfg
		Leg          *LegCfg
		LegFastening *LegFasteningCfg
		Tank         *TankCfg
	}
	var cps [][]CP

	checkErr(json.NewDecoder(reader).Decode(&cps))
	result := make([][]Part, len(cps))
	for y, cpRow := range cps {
		result[y] = make([]Part, len(cpRow))
		for x, cp := range cpRow {
			if cp.Cabin != nil {
				result[y][x] = NewCabin(*cp.Cabin)
				continue
			}
			if cp.Crain != nil {
				result[y][x] = NewCrane(*cp.Crain)
				continue
			}
			if cp.Engine != nil {
				result[y][x] = NewEngine(*cp.Engine)
				continue
			}
			if cp.Leg != nil {
				result[y][x] = NewLeg(*cp.Leg)
				continue
			}
			if cp.LegFastening != nil {
				result[y][x] = NewLegFastening(*cp.LegFastening)
				continue
			}
			if cp.Tank != nil {
				result[y][x] = NewTank(*cp.Tank)
				continue
			}

		}
	}
	return result
}

func TestName(t *testing.T) {
	parts := ParseShipParts(strings.NewReader(j))
	t.Log(parts)
}
