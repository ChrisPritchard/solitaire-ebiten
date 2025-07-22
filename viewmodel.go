package main

import (
	"github.com/chrispritchard/solitaire-ebiten/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageData struct {
	Image *ebiten.Image
	X, Y  float64
}

type drag_state struct {
	card     *Card
	offset_x int
	offset_y int
}

type ViewModel struct {
	CardWidth, CardHeight float64
	dragged_cards         []drag_state
}

func (vm *ViewModel) pixels_to_card_units(x, y float64) (int, int) {
	cux := x / int(vm.CardWidth/2)
	cuy := y / int(vm.CardHeight/4)
	return cux, cuy
}

func (vm *ViewModel) card_units_to_pixels(cux, cuy int) (float64, float64) {
	x := float64(cux) * vm.CardWidth / 2
	y := float64(cuy) * vm.CardHeight / 4
	return x, y
}

func (vm *ViewModel) Update(ts TouchState, game SawayamaRules) error {
	if ts.Pressed && vm.dragged_cards == nil {
		// test if card can be picked up
		cux, cuy := vm.pixels_to_card_units(float64(ts.X), float64(ts.Y))

		cards := game.DraggableAt(cux, cuy)
		vm.dragged_cards = []drag_state{}
		for _, c := range cards {
			x, y := vm.card_units_to_pixels(c.CUX, c.CUY)
			vm.dragged_cards = append(vm.dragged_cards, drag_state{
				card:     c,
				offset_x: x - float64(ts.X),
				offset_y: y - float64(ts.Y),
			})
		}
	} else if !ts.Pressed && vm.dragged_cards != nil {
		// test if card can be dropped
	} else if vm.dragged_cards != nil {
		// update offsets
	}

	return nil
}

func (vm *ViewModel) Transform(game SawayamaRules) ([]ImageData, error) {

	res := []ImageData{}

	for _, c := range game.Cards {
		x, y := vm.card_units_to_pixels(c.CUX, c.CUY)

		for _, d := range vm.dragged_cards {
			if *d.card == c {
				x += float64(d.offset_x)
				y += float64(d.offset_y)
				break
			}
		}

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
