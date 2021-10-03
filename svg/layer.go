package svg

type Layer struct {
	ID       string    `xml:"id,attr"`
	Pathes   []Path    `xml:"path"`
	Rects    []Rect    `xml:"rect"`
	Ellipses []Ellipse `xml:"ellipse"`
}
