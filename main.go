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

	game := Game{rules: &SawayamaRules{}}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
