package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sawayama Solitaire")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
