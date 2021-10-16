package main

import (
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"strconv"
	"strings"
)

type PartType string

const (
	PartTypeBox          PartType = "box"
	PartTypeCabin        PartType = "cab"
	PartTypeEngine       PartType = "eng"
	PartTypeCrane        PartType = "crn"
	PartTypeLeg          PartType = "leg"
	PartTypeLegFastening PartType = "lft"
)

type PartParam string

const (
	PartParamDir   PartParam = "dir"
	PartParamPower PartParam = "pow"
	PartParamKeys  PartParam = "keys"
)

type PartDef interface {
	Construct(world *box2d.B2World,
		tanker Tank,
		ps *ParticleSystem,
		shipPos box2d.B2Vec2,
		shipSize box2d.B2Vec2,
		pos box2d.B2Vec2) Part
}

func ParsePartDef(strP *string) PartDef {
	if strP == nil {
		return nil
	}
	str := *strP
	tp := strings.ToLower(str[:3])
	params := parsePartParams(str[3:])
	switch PartType(tp) {
	case PartTypeBox:
		return &BoxDef{}
	case PartTypeCabin:
		return &CabinDef{
			Dir: params[PartParamDir].AsDirection(),
		}
	case PartTypeEngine:
		return &EngineDef{
			Dir:   params[PartParamDir].AsDirection(),
			Power: params[PartParamPower].AsFloat(),
			Keys:  params[PartParamKeys].AsKeys(),
		}
	case PartTypeCrane:
		return &CraneDef{
			Dir: params[PartParamDir].AsDirection(),
		}
	case PartTypeLeg:
		return &LegDef{
			Dir: params[PartParamDir].AsDirection(),
		}
	case PartTypeLegFastening:
		return &LegFasteningDef{
			Dir: params[PartParamDir].AsDirection(),
		}
	default:
		checkErr(fmt.Errorf("unknown part type %v", tp))
		return nil
	}
}

type paramVal string

func (v paramVal) AsDirection() Direction {
	return Direction(v)
}
func (v paramVal) AsFloat() float64 {
	result, err := strconv.ParseFloat(string(v), 64)
	checkErr(err)
	return result
}
func (v paramVal) AsKeys() Keys {
	result := make(Keys)
	for _, elem := range strings.Split(string(v), ",") {
		k, err := strconv.Atoi(elem)
		checkErr(err)
		result[ebiten.Key(k)] = struct{}{}
	}
	return result
}

func parsePartParams(str string) map[PartParam]paramVal {
	result := make(map[PartParam]paramVal)
	params := strings.Split(str, ";")
	if len(params) == 1 && params[0] == "" {
		return result
	}
	for _, param := range params {
		kv := strings.Split(param, "=")
		if len(kv) != 2 {
			checkErr(fmt.Errorf("bad part param %v", str))
		}
		result[PartParam(strings.Trim(kv[0], " "))] = paramVal(strings.Trim(kv[1], " "))
	}
	return result
}
