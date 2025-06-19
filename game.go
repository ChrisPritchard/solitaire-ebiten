package main

import (
	"sort"

	"github.com/chrispritchard/solitaire-ebiten/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type Card struct {
	X, Y, Z int
	Content CardContent
}

type CardContent interface {
	Image() *ebiten.Image
}

type RuleSet interface {
	Update() error
	Cards() []Card
}

type Game struct {
	rules RuleSet
}

func (game *Game) Update() error {
	return game.rules.Update()
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.Background, nil)

	cards := game.rules.Cards()
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Z < cards[j].Z
	})

	for _, card := range cards {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(3, 3)
		op.GeoM.Translate(float64(card.X), float64(card.Y))

		img := card.Content.Image()
		screen.DrawImage(img, op)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
