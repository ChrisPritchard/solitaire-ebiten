package main

import (
	"github.com/chrispritchard/solitaire-ebiten/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageData struct {
	Image *ebiten.Image
	X, Y  float64
}

const (
	card_width  = 36. * 3.
	card_height = 54. * 3.
)

func Transform(game SawayamaRules) ([]ImageData, error) {

	res := []ImageData{}

	for _, c := range game.Cards {
		x := float64(c.CUX) * card_width / 2
		y := float64(c.CUY) * card_height / 4 // todo: needs to be even?
		var image *ebiten.Image
		if !c.Visible {
			image = assets.CardBack
		} else {
			image = assets.Cards[c.Suit][c.Value]
		}

		res = append(res, ImageData{
			X:     x,
			Y:     y,
			Image: image,
		})
	}

	return res, nil
}
