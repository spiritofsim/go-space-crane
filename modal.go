package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

type Modal struct {
	Title     string
	Text      string
	CloseKeys Keys
	isClosed  bool
}

func NewModal(title string, text string, closeKeys Keys) *Modal {
	return &Modal{
		Title:     title,
		Text:      text,
		CloseKeys: closeKeys,
	}
}

func (m *Modal) Update(keys Keys) {
	found := false
	for key := range m.CloseKeys {
		if keys.IsPressed(key) {
			found = true
			break
		}
	}
	if found {
		m.isClosed = true
	}
}

func (m *Modal) Draw(screen *ebiten.Image) {
	if m.isClosed {
		return
	}

	screen.DrawImage(modalImg, nil)

	titleBounds := text.BoundString(modalTitleFace, m.Title)
	text.Draw(screen, m.Title, modalTitleFace, ScreenWidth/2-titleBounds.Max.X/2, 280, color.White)

	textBounds := text.BoundString(modalTextFace, m.Text)
	text.Draw(screen, m.Text, modalTextFace, ScreenWidth/2-textBounds.Max.X/2, 500, color.White)

}
