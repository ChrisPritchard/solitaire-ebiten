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
	offset_x float64
	offset_y float64
}

type ViewModel struct {
	CardWidth, CardHeight float64
	dragged_cards         []drag_state
	cursor_x, cursor_y    float64
}

func (vm *ViewModel) pixels_to_card_units(x, y float64) (int, int) {
	cux := int(x / vm.CardWidth * CUX_per_card)
	cuy := int(y / vm.CardHeight * CUY_per_card)
	return cux, cuy
}

func (vm *ViewModel) card_units_to_pixels(cux, cuy int) (float64, float64) {
	x := float64(cux) / CUX_per_card * vm.CardWidth
	y := float64(cuy) / CUY_per_card * vm.CardHeight
	return x, y
}

func (vm *ViewModel) Update(ts TouchState, game *SawayamaRules) error {
	game.Update()
	vm.cursor_x = ts.X
	vm.cursor_y = ts.Y

	if ts.Pressed && vm.dragged_cards == nil {
		cux, cuy := vm.pixels_to_card_units(ts.X, ts.Y)

		cards := game.DraggableAt(cux, cuy)
		vm.dragged_cards = []drag_state{}
		for _, c := range cards {
			vm.dragged_cards = append(vm.dragged_cards, drag_state{
				card:     c,
				offset_x: ts.X,
				offset_y: ts.Y,
			})
		}
	} else if !ts.Pressed && vm.dragged_cards != nil {
		cux, cuy := vm.pixels_to_card_units(ts.X, ts.Y)
		head_card := vm.dragged_cards[0]
		can_be_dropped, cux, cuy := game.DroppableAt(cux, cuy, head_card.card.Suit, head_card.card.Value)
		if can_be_dropped {
			for i, d := range vm.dragged_cards {
				(*d.card).CUX = cux
				(*d.card).CUY = cuy + i + 1
			}
		}
		vm.dragged_cards = nil
	}

	return nil
}

func (vm *ViewModel) Transform(game SawayamaRules) []ImageData {

	res := []ImageData{}
	dragged := []ImageData{}

	for _, c := range game.Cards {
		x, y := vm.card_units_to_pixels(c.CUX, c.CUY)
		is_dragged := false

		for _, d := range vm.dragged_cards {
			if *d.card == c {
				x += float64(vm.cursor_x - d.offset_x)
				y += float64(vm.cursor_y - d.offset_y)
				is_dragged = true
				break
			}
		}

		var image *ebiten.Image
		if !c.Visible {
			image = assets.CardBack
		} else {
			image = assets.Cards[c.Suit][c.Value]
		}

		if is_dragged {
			dragged = append(dragged, ImageData{
				X:     x,
				Y:     y,
				Image: image,
			})
		} else {
			res = append(res, ImageData{
				X:     x,
				Y:     y,
				Image: image,
			})
		}
	}

	res = append(res, dragged...) // ensuring they're on top
	return res
}
