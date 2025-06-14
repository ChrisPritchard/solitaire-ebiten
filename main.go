package main

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

//go:embed assets/card-suites.png
var card_data []byte

//go:embed assets/felt_green.jpg
var background []byte

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sawayama Solitaire")

	if err := ebiten.RunGame(NewGame(card_data, background)); err != nil {
		log.Fatal(err)
	}
}
