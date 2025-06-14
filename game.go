package main

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Card struct {
	X, Y  int
	Image *ebiten.Image
}

type Game struct {
	cards []Card
}

func (game *Game) Update() error {
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	for _, card := range game.cards {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(card.X), float64(card.Y))
		screen.DrawImage(card.Image, op)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	cards := make([]Card, 0)

	file, err := os.Open("./assets/clubs-ace.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	img2 := ebiten.NewImageFromImage(img)
	card := Card{X: 100, Y: 100, Image: img2}
	cards = append(cards, card)

	return &Game{cards}
}
