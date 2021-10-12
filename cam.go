package main

import "github.com/ByteArena/box2d"

type Cam struct {
	Pos  box2d.B2Vec2
	Zoom float64
}

func NewCam() *Cam {
	return &Cam{
		Pos:  box2d.MakeB2Vec2(50, 20),
		Zoom: 20,
	}
}

func (c *Cam) Project(v box2d.B2Vec2, pos box2d.B2Vec2, ang float64) box2d.B2Vec2 {
	v = box2d.B2RotVec2Mul(*box2d.NewB2RotFromAngle(ang), v)                // rotate
	v = box2d.B2Vec2Add(v, pos)                                             // set position
	v = box2d.B2Vec2Add(v, box2d.MakeB2Vec2(-c.Pos.X, -c.Pos.Y))            // camera pos
	v = box2d.B2Vec2MulScalar(c.Zoom, v)                                    // camera zoom
	v = box2d.B2Vec2Add(v, box2d.MakeB2Vec2(ScreenWidth/2, ScreenHeight/2)) // cam screen center
	return v
}
