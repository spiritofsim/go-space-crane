package svg

import (
	"encoding/xml"
	"github.com/ByteArena/box2d"
)

type Rect struct {
	ID          string
	Title       string
	Description string
	Pos         box2d.B2Vec2
	Size        box2d.B2Vec2
}

func (r *Rect) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var prj struct {
		ID          string  `xml:"id,attr"`
		Title       string  `xml:"title"`
		Description string  `xml:"desc"`
		X           float64 `xml:"x,attr"`
		Y           float64 `xml:"y,attr"`
		Width       float64 `xml:"width,attr"`
		Height      float64 `xml:"height,attr"`
	}

	if err := d.DecodeElement(&prj, &start); err != nil {
		return err
	}

	*r = Rect{
		ID:          prj.ID,
		Title:       prj.Title,
		Description: prj.Description,
		Pos:         box2d.MakeB2Vec2(prj.X, prj.Y),
		Size:        box2d.MakeB2Vec2(prj.Width, prj.Height),
	}
	return nil
}
