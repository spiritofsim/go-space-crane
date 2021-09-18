package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
)

func DrawDebugBody(screen *ebiten.Image, body *box2d.B2Body, cam Cam, clr color.Color) {
	for fix := body.GetFixtureList(); fix != nil; fix = fix.GetNext() {
		switch shape := fix.GetShape().(type) {
		case *box2d.B2CircleShape:
			DrawDebugCircleShape(screen, body, shape, cam, clr)
		case *box2d.B2PolygonShape:
			DrawDebugPolygonShape(screen, body, shape, cam, clr)
		case *box2d.B2ChainShape:
			DrawDebugChainShape(screen, body, shape, cam, clr)
		case *box2d.B2EdgeShape:
			DrawDebugEdgeShape(screen, body, shape, cam, clr)
		}
	}

	// Center
	crossSize := 0.3
	pos := body.GetPosition()
	ang := body.GetAngle()
	p1 := box2d.MakeB2Vec2(-crossSize/2, 0)
	p2 := box2d.MakeB2Vec2(crossSize/2, 0)
	p3 := box2d.MakeB2Vec2(0, -crossSize/2)
	p4 := box2d.MakeB2Vec2(0, crossSize/2)

	p1 = cam.Project(p1, pos, ang)
	p2 = cam.Project(p2, pos, ang)
	p3 = cam.Project(p3, pos, ang)
	p4 = cam.Project(p4, pos, ang)

	crossClr := color.RGBA{
		R: 0,
		G: 0xff,
		B: 0,
		A: 0xff,
	}
	ebitenutil.DrawLine(screen, p1.X, p1.Y, p2.X, p2.Y, crossClr)
	ebitenutil.DrawLine(screen, p3.X, p3.Y, p4.X, p4.Y, crossClr)
}

func DrawDebugChainShape(screen *ebiten.Image, body *box2d.B2Body, shape *box2d.B2ChainShape, cam Cam, clr color.Color) {
	drawDebugPolyFromVerts(screen, body.GetPosition(), body.GetAngle(), shape.M_vertices[:shape.M_count], cam, clr)
}

func DrawDebugEdgeShape(screen *ebiten.Image, body *box2d.B2Body, shape *box2d.B2EdgeShape, cam Cam, clr color.Color) {
	v1 := cam.Project(shape.M_vertex1, body.GetPosition(), body.GetAngle())
	v2 := cam.Project(shape.M_vertex2, body.GetPosition(), body.GetAngle())
	ebitenutil.DrawLine(
		screen,
		v1.X, v1.Y,
		v2.X, v2.Y,
		clr)
}

func DrawDebugPolygonShape(screen *ebiten.Image, body *box2d.B2Body, shape *box2d.B2PolygonShape, cam Cam, clr color.Color) {
	drawDebugPolyFromVerts(screen, body.GetPosition(), body.GetAngle(), shape.M_vertices[:shape.M_count], cam, clr)
}

func DrawDebugCircleShape(screen *ebiten.Image, body *box2d.B2Body, shape *box2d.B2CircleShape, cam Cam, clr color.Color) {
	count := box2d.B2_maxPolygonVertices
	verts := make([]box2d.B2Vec2, count)
	for i := 0; i < count; i++ {
		ang := 2 * math.Pi * float64(i) / float64(count)
		r := shape.GetRadius()
		verts[i] = box2d.MakeB2Vec2(math.Cos(ang)*r, math.Sin(ang)*r)
	}

	polyShape := box2d.MakeB2PolygonShape()
	polyShape.Set(verts, len(verts))
	DrawDebugPolygonShape(screen, body, &polyShape, cam, clr)
}

// TODO: use box2d.B2Transform instead of pos+ang
func drawDebugPolyFromVerts(screen *ebiten.Image, pos box2d.B2Vec2, ang float64, verts []box2d.B2Vec2, cam Cam, clr color.Color) {
	for i := 1; i < len(verts); i++ {
		v1 := cam.Project(verts[i], pos, ang)
		v2 := cam.Project(verts[i-1], pos, ang)
		ebitenutil.DrawLine(
			screen,
			v1.X, v1.Y,
			v2.X, v2.Y,
			clr)
	}

	first := cam.Project(verts[0], pos, ang)
	last := cam.Project(verts[len(verts)-1], pos, ang)
	ebitenutil.DrawLine(
		screen,
		first.X, first.Y,
		last.X, last.Y,
		clr)
}
