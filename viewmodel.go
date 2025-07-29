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
	cards  []Card
	offset Vec2[float64]
	origin Vec2[int] // used to remove a dropped card from where it came from, if valid
}

type ViewModel struct {
	CardSize       Vec2[float64]
	dragged_cards  *drag_state
	cursor         Vec2[float64]
	stacking       *Stackable
	stack_progress float64
}

func (vm *ViewModel) pixels_to_card_units(p Vec2[float64]) Vec2[int] {
	return p.Divide(vm.CardSize).Scale(CU_per_card.ToFloat()).ToInt()
}

func (vm *ViewModel) card_units_to_pixels(c Vec2[int]) Vec2[float64] {
	return c.ToFloat().Divide(CU_per_card.ToFloat()).Scale(vm.CardSize)
}

func play_sound(index int) {
	assets.Sounds[index].Rewind()
	assets.Sounds[index].Play()
}

func (vm *ViewModel) Update(ts TouchState, game *SawayamaRules) {

	vm.cursor = ts.Pos

	if vm.stacking != nil {

		return
	} else if stackable := game.NextStackable(); stackable != nil {
		vm.stacking = stackable
		vm.stack_progress = 0
		return
	}

	if ts.Pressed && ts.JustChanged && vm.dragged_cards == nil {
		cu := vm.pixels_to_card_units(ts.Pos)

		if Deck_CU.Contains(cu, CU_per_card) && game.DrawFromDeck() {
			play_sound(0)
			return
		}

		cards, origin := game.DraggableAt(cu)
		if cards == nil {
			return
		}
		play_sound(1)
		offset := ts.Pos
		vm.dragged_cards = &drag_state{cards, offset, origin}
	} else if !ts.Pressed && vm.dragged_cards != nil {
		cu := vm.pixels_to_card_units(ts.Pos)
		play_sound(2)
		game.DropAt(cu, vm.dragged_cards.cards, vm.dragged_cards.origin)
		vm.dragged_cards = nil
	}
}

func (vm *ViewModel) Transform(game SawayamaRules) []ImageData {

	res := []ImageData{}
	dragged := []ImageData{}

	for _, c := range game.Cards() {
		p := vm.card_units_to_pixels(c.Pos)
		is_dragged := false

		if vm.dragged_cards != nil {
			for _, d := range vm.dragged_cards.cards {
				if d.Equals(c.Card) {
					p = p.Add2(vm.cursor.Subtract2(vm.dragged_cards.offset))
					is_dragged = true
					break
				}
			}
		}

		var image *ebiten.Image
		if !c.Visible {
			image = assets.CardBack
		} else if c.Value > 0 {
			image = assets.Cards[c.Suit][c.Value]
		} else {
			image = assets.CardSpace // hackity hack hack
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
