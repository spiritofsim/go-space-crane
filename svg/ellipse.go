package svg

import (
	"encoding/xml"
	"github.com/ByteArena/box2d"
)

type Ellipse struct {
	// TODO: id, title, desc to struct
	ID          string
	Title       string
	Description string

	Pos    box2d.B2Vec2
	Radius box2d.B2Vec2
}

func (e *Ellipse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var prj struct {
		ID          string  `xml:"id,attr"`
		Title       string  `xml:"title"`
		Description string  `xml:"desc"`
		Cx          float64 `xml:"cx,attr"`
		Cy          float64 `xml:"cy,attr"`
		Rx          float64 `xml:"rx,attr"`
		Ry          float64 `xml:"ry,attr"`
	}

	if err := d.DecodeElement(&prj, &start); err != nil {
		return err
	}

	*e = Ellipse{
		ID:          prj.ID,
		Title:       prj.Title,
		Description: prj.Description,
		Pos:         box2d.MakeB2Vec2(prj.Cx, prj.Cy),
		Radius:      box2d.MakeB2Vec2(prj.Rx, prj.Ry),
	}
	return nil
}
