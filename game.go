package main

import (
	"github.com/chrispritchard/solitaire-ebiten/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Card struct {
	X, Y        int
	Suit, Value int
}

type Game struct {
	cards []Card
}

func (game *Game) Update() error {
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.Background, nil)

	for _, card := range game.cards {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(card.X), float64(card.Y))
		img := assets.Cards[card.Suit][card.Value]
		screen.DrawImage(img, op)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {

	cards := make([]Card, 0)
	for i := 0; i < 4; i++ {
		for j := 2; j <= 14; j++ {
			card := Card{X: 10 + i*50 + j*5, Y: 10 + j*10, Suit: i, Value: j}
			cards = append(cards, card)
		}
	}

	return &Game{cards}
}
