package main

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Card struct {
	X, Y  int
	Image *ebiten.Image
}

type Game struct {
	cards      []Card
	background *ebiten.Image
}

func (game *Game) Update() error {
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(game.background, nil)

	for _, card := range game.cards {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(card.X), float64(card.Y))
		screen.DrawImage(card.Image, op)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame(card_data []byte, background_data []byte) *Game {

	background_image, _, err := image.Decode(bytes.NewReader(background_data))
	if err != nil {
		log.Fatal(err)
	}
	background := ebiten.NewImageFromImage(background_image)

	cards_image, _, err := image.Decode(bytes.NewReader(card_data))
	if err != nil {
		log.Fatal(err)
	}
	cards_image2 := ebiten.NewImageFromImage(cards_image)

	card := Card{X: 100, Y: 100, Image: cards_image2}
	cards := []Card{card}

	return &Game{cards, background}
}
