package main

import (
	"github.com/chrispritchard/solitaire-ebiten/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageData struct {
	Image *ebiten.Image
	Pos   Vec2[float64]
}

type drag_state struct {
	card   *Card
	offset Vec2[float64]
}

type ViewModel struct {
	CardSize      Vec2[float64]
	dragged_cards []drag_state
	cursor        Vec2[float64]
}

func (vm *ViewModel) pixels_to_card_units(p Vec2[float64]) Vec2[int] {
	return p.Divide(vm.CardSize).Scale(CU_per_card.ToFloat()).ToInt()
}

func (vm *ViewModel) card_units_to_pixels(c Vec2[int]) Vec2[float64] {
	return c.ToFloat().Divide(CU_per_card.ToFloat()).Scale(vm.CardSize)
}

func (vm *ViewModel) Update(ts TouchState, game *SawayamaRules) error {

	vm.cursor = ts.Pos

	if ts.Pressed && ts.JustChanged && vm.dragged_cards == nil {
		cu := vm.pixels_to_card_units(ts.Pos)

		cards := game.DraggableAt(cu)
		if cards == nil {
			return nil
		}
		vm.dragged_cards = []drag_state{}
		for _, c := range cards {
			vm.dragged_cards = append(vm.dragged_cards, drag_state{
				card:   c,
				offset: ts.Pos,
			})
		}
	} else if !ts.Pressed && vm.dragged_cards != nil {
		cu := vm.pixels_to_card_units(ts.Pos)
		cards := []*Card{}
		for _, c := range vm.dragged_cards {
			cards = append(cards, c.card)
		}
		game.DropAt(cu, cards)
		vm.dragged_cards = nil
	}

	return nil
}

func (vm *ViewModel) Transform(game SawayamaRules) []ImageData {

	res := []ImageData{}
	dragged := []ImageData{}

	for _, c := range game.Cards {
		p := vm.card_units_to_pixels(c.Pos)
		is_dragged := false

		for _, d := range vm.dragged_cards {
			if *d.card == c {
				p = p.Add2(vm.cursor.Subtract2(d.offset))
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
			dragged = append(dragged, ImageData{image, p})
		} else {
			res = append(res, ImageData{image, p})
		}
	}

	res = append(res, dragged...) // ensuring they're on top
	return res
}
