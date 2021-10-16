package main

import (
	"github.com/ByteArena/box2d"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
	"image/color"
	"io/ioutil"
	"path"
	"time"
)

const ()

type LevelSelector struct {
	levels    []LevelDef
	fileNames map[int]string
	index     int

	lastControlled time.Time
}

func NewLevelSelector() *LevelSelector {
	files, err := ioutil.ReadDir(LevelsDir)
	checkErr(err)
	levels := make([]LevelDef, len(files))
	fileNames := make(map[int]string)
	for i, file := range files {
		levelDefData, err := ioutil.ReadFile(path.Join(LevelsDir, file.Name()))
		checkErr(err)
		var levelDef LevelDef
		checkErr(yaml.Unmarshal(levelDefData, &levelDef))
		levels[i] = levelDef
		fileNames[i] = file.Name()
	}
	return &LevelSelector{
		levels:         levels,
		fileNames:      fileNames,
		lastControlled: time.Now(),
	}
}

func (ls *LevelSelector) Update() error {
	if time.Since(ls.lastControlled) < time.Second/10 {
		return nil
	}
	ls.lastControlled = time.Now()

	keys := KeysFromSlice(inpututil.AppendPressedKeys(nil))

	if keys.IsKeyPressed(ebiten.KeyEnter) || keys.IsKeyPressed(ebiten.KeyNumpadEnter) {
		return NewLevelSelected(ls.fileNames[ls.index])
	}
	if keys.IsKeyPressed(ebiten.KeyArrowDown) {
		ls.index++
		if ls.index >= len(ls.levels) {
			ls.index = 0
		}
	}
	if keys.IsKeyPressed(ebiten.KeyArrowUp) {
		ls.index--
		if ls.index < 0 {
			ls.index = len(ls.levels) - 1
		}
	}
	return nil
}

func (ls *LevelSelector) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	opts := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}

	levelBoxSize := box2d.MakeB2Vec2(500, 100)
	levelBoxStartPos := box2d.MakeB2Vec2(100, 100)
	levelBoxDist := 30.0
	selectedBorderSize := float32(5.0)

	for i, level := range ls.levels {
		if i == ls.index {
			var path vector.Path
			x := float32(levelBoxStartPos.X)
			y := float32(levelBoxStartPos.Y + float64(i)*levelBoxSize.Y + float64(i)*levelBoxDist)
			path.MoveTo(x-selectedBorderSize, y-selectedBorderSize)
			path.LineTo(x+float32(levelBoxSize.X)+selectedBorderSize, y-selectedBorderSize)
			path.LineTo(x+float32(levelBoxSize.X)+selectedBorderSize, y+float32(levelBoxSize.Y)+selectedBorderSize)
			path.LineTo(x-selectedBorderSize, y+float32(levelBoxSize.Y)+selectedBorderSize)
			vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
			for i := range vs {
				vs[i].SrcX = 1
				vs[i].SrcY = 1
				vs[i].ColorR = 0.5
				vs[i].ColorG = 0.5
				vs[i].ColorB = 0.5
			}
			screen.DrawTriangles(vs, is, emptySubImage, opts)

			text.Draw(
				screen,
				level.Name,
				levelNameFace,
				1200,
				int(levelBoxStartPos.X),
				color.Black)
			text.Draw(
				screen,
				level.Description,
				levelDescFace,
				800,
				int(levelBoxStartPos.X)+100,
				color.Black)
		}

		var path vector.Path
		x := float32(levelBoxStartPos.X)
		y := float32(levelBoxStartPos.Y + float64(i)*levelBoxSize.Y + float64(i)*levelBoxDist)
		path.MoveTo(x, y)
		path.LineTo(x+float32(levelBoxSize.X), y)
		path.LineTo(x+float32(levelBoxSize.X), y+float32(levelBoxSize.Y))
		path.LineTo(x, y+float32(levelBoxSize.Y))

		vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vs {
			vs[i].SrcX = 1
			vs[i].SrcY = 1
			vs[i].ColorR = 0
			vs[i].ColorG = 0
			vs[i].ColorB = 0
		}
		screen.DrawTriangles(vs, is, emptySubImage, opts)

		bounds := text.BoundString(levelNameFace, level.Name)
		text.Draw(
			screen,
			level.Name,
			levelNameFace,
			int(x+float32(levelBoxSize.X)/2)-bounds.Max.X/2,
			int(y+float32(levelBoxSize.Y)/2)-bounds.Max.Y/2,
			color.White)
	}

}

func (ls *LevelSelector) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
