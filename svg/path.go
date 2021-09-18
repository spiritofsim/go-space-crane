package svg

import (
	"encoding/xml"
	"fmt"
	"github.com/ByteArena/box2d"
	"strconv"
	"strings"
)

type Path struct {
	ID          string
	Title       string
	Description string
	Verts       []box2d.B2Vec2
}

func (p *Path) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var prj struct {
		ID          string `xml:"id,attr"`
		Data        string `xml:"d,attr"`
		Title       string `xml:"title"`
		Description string `xml:"desc"`
	}

	if err := d.DecodeElement(&prj, &start); err != nil {
		return err
	}

	verts, err := parsePath(prj.Data)
	if err != nil {
		return err
	}

	*p = Path{
		ID:          prj.ID,
		Title:       prj.Title,
		Description: prj.Description,
		Verts:       verts,
	}
	return nil
}

func parsePath(str string) ([]box2d.B2Vec2, error) {
	parts := strings.Split(str, " ")

	verts := make([]box2d.B2Vec2, 0)
	cmd := ""
	var prev *box2d.B2Vec2 = nil
	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case "M", "m", "L", "l", "V", "v", "H", "h":
			cmd = parts[i]
		case "Z", "z": // we finish here because don't need 8like shapes
			return verts, nil
		default:
			var v box2d.B2Vec2
			switch cmd {
			case "M", "L":
				_v, err := parseCoords(parts[i])
				if err != nil {
					return nil, err
				}
				v = _v
			case "m", "l":
				_v, err := parseCoords(parts[i])
				if err != nil {
					return nil, err
				}
				if prev == nil {
					v = _v
				} else {
					v = box2d.MakeB2Vec2(prev.X+_v.X, prev.Y+_v.Y)
				}
			case "V":
				y, err := strconv.ParseFloat(parts[i], 64)
				if err != nil {
					return nil, err
				}
				v = box2d.MakeB2Vec2(prev.X, y)
			case "v":
				y, err := strconv.ParseFloat(parts[i], 64)
				if err != nil {
					return nil, err
				}
				v = box2d.MakeB2Vec2(prev.X, prev.Y+y)
			case "H":
				x, err := strconv.ParseFloat(parts[i], 64)
				if err != nil {
					return nil, err
				}
				v = box2d.MakeB2Vec2(x, prev.Y)
			case "h":
				x, err := strconv.ParseFloat(parts[i], 64)
				if err != nil {
					return nil, err
				}
				v = box2d.MakeB2Vec2(prev.X+x, prev.Y)
			}

			prev = &v
			verts = append(verts, v)
		}
	}

	return verts, nil
}

func parseCoords(str string) (box2d.B2Vec2, error) {
	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		return box2d.B2Vec2{}, fmt.Errorf("2 coords expected in %v", str)
	}

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return box2d.B2Vec2{}, err
	}

	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return box2d.B2Vec2{}, err
	}

	return box2d.MakeB2Vec2(x, y), nil
}
