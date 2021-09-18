package svg

import (
	"github.com/ByteArena/box2d"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseSvgPath(t *testing.T) {
	verts, err := parsePath("M 324.59604,0.08457427 325.01891,0.16914854 509.13709,108.25506 598.10922,285.86103 507.44561,513.19666 279.77168,590.66669 91.340209,465.49677 1.3531883,290.59718 105.88698,108.59336 Z")
	require.NoError(t, err)
	require.Len(t, verts, 9)

	verts, err = parsePath("M 0,-600 560,120 480,300 140,380 H -140 L -480,300 -560,120 Z")
	require.NoError(t, err)
	require.Len(t, verts, 7)
	require.Equal(t, box2d.MakeB2Vec2(0, -600), verts[0])
	require.Equal(t, box2d.MakeB2Vec2(560, 120), verts[1])
	require.Equal(t, box2d.MakeB2Vec2(480, 300), verts[2])
	require.Equal(t, box2d.MakeB2Vec2(140, 380), verts[3])
	require.Equal(t, box2d.MakeB2Vec2(-140, 380), verts[4])
	require.Equal(t, box2d.MakeB2Vec2(-480, 300), verts[5])
	require.Equal(t, box2d.MakeB2Vec2(-560, 120), verts[6])

	verts, err = parsePath("M 1,2 l 1,1 z")
	require.NoError(t, err)
	require.Len(t, verts, 2)
	require.Equal(t, box2d.MakeB2Vec2(1, 2), verts[0])
	require.Equal(t, box2d.MakeB2Vec2(2, 3), verts[1])

	verts, err = parsePath("M 0,1 H 2")
	require.NoError(t, err)
	require.Len(t, verts, 2)
	require.Equal(t, box2d.MakeB2Vec2(0, 1), verts[0])
	require.Equal(t, box2d.MakeB2Vec2(2, 1), verts[1])

	verts, err = parsePath("M 0,800 560,80 1120,800 1020,980 680,1060 H 440 L 100,980 Z")
	require.NoError(t, err)
	require.Len(t, verts, 7)

	verts, err = parsePath("M 1,2 V 3")
	require.NoError(t, err)
	require.Len(t, verts, 2)
	require.Equal(t, box2d.MakeB2Vec2(1, 2), verts[0])
	require.Equal(t, box2d.MakeB2Vec2(1, 3), verts[1])

	verts, err = parsePath("M 1,2 3,4 5,6 7,8")
	require.NoError(t, err)
	require.Len(t, verts, 4)
	require.Equal(t, box2d.MakeB2Vec2(1, 2), verts[0])
	require.Equal(t, box2d.MakeB2Vec2(3, 4), verts[1])
	require.Equal(t, box2d.MakeB2Vec2(5, 6), verts[2])
	require.Equal(t, box2d.MakeB2Vec2(7, 8), verts[3])

	verts, err = parsePath("M 1,2 3,4 5,6 7,8 Z")
	require.NoError(t, err)
	require.Len(t, verts, 4)
	require.Equal(t, box2d.MakeB2Vec2(1, 2), verts[0])
	require.Equal(t, box2d.MakeB2Vec2(3, 4), verts[1])
	require.Equal(t, box2d.MakeB2Vec2(5, 6), verts[2])
	require.Equal(t, box2d.MakeB2Vec2(7, 8), verts[3])
}
