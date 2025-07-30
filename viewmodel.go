package main

import (
	"math"

	"github.com/chrispritchard/solitaire-ebiten/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
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

type stack_state struct {
	card        Card
	origin      Vec2[int]
	destination Vec2[int]
	progress    float64
}

func (s stack_state) distance(vm *ViewModel) Vec2[float64] {
	origin := vm.card_units_to_pixels(vm.stacking.origin)
	dest := vm.card_units_to_pixels(vm.stacking.destination)
	diff := dest.Subtract2(origin)
	direct := math.Sqrt(diff.X*diff.X + diff.Y*diff.Y)
	progress := direct * vm.stacking.progress
	scale := progress / direct
	curr := diff.Scale(Vec2[float64]{scale, scale})
	return curr
}

var progress_per_dt = (1. / 60.) * 4

type ViewModel struct {
	CardSize      Vec2[float64]
	dragged_cards *drag_state
	cursor        Vec2[float64]
	stacking      *stack_state
}

func NewViewModel(cardSize Vec2[float64]) ViewModel {
	play_sound(assets.Sounds.NewGame)
	return ViewModel{CardSize: cardSize}
}

func (vm *ViewModel) pixels_to_card_units(p Vec2[float64]) Vec2[int] {
	return p.Divide(vm.CardSize).Scale(CU_per_card.ToFloat()).ToInt()
}

func (vm *ViewModel) card_units_to_pixels(c Vec2[int]) Vec2[float64] {
	return c.ToFloat().Divide(CU_per_card.ToFloat()).Scale(vm.CardSize)
}

func play_sound(p *audio.Player) {
	p.Rewind()
	p.Play()
}

func (vm *ViewModel) Update(ts TouchState, game *SawayamaRules) {

	vm.cursor = ts.Pos

	if vm.stacking != nil {
		vm.stacking.progress += progress_per_dt
		if vm.stacking.progress >= 1. {
			game.DropAt(vm.stacking.destination, []Card{vm.stacking.card}, vm.stacking.origin)
			vm.stacking = nil
		}
		return
	} else if stackable := game.NextStackable(); stackable != nil {
		vm.stacking = &stack_state{card: stackable.Card,
			origin:      stackable.Origin,
			destination: stackable.Destination}
		play_sound(assets.Sounds.Stack)
		return
	}

	if ts.Pressed && ts.JustChanged && vm.dragged_cards == nil {
		cu := vm.pixels_to_card_units(ts.Pos)

		if Deck_CU.Contains(cu, CU_per_card) && game.DrawFromDeck() {
			play_sound(assets.Sounds.DrawDeck)
			return
		}

		cards, origin := game.DraggableAt(cu)
		if cards == nil {
			return
		}
		offset := ts.Pos
		vm.dragged_cards = &drag_state{cards, offset, origin}
		play_sound(assets.Sounds.DragStart)
	} else if !ts.Pressed && vm.dragged_cards != nil {
		cu := vm.pixels_to_card_units(ts.Pos)
		game.DropAt(cu, vm.dragged_cards.cards, vm.dragged_cards.origin)
		vm.dragged_cards = nil
		play_sound(assets.Sounds.DragDrop)
	}
}

func (vm *ViewModel) Transform(game SawayamaRules) []ImageData {
	res := []ImageData{}
	moving := []ImageData{}

	for _, c := range game.Cards() {
		p := vm.card_units_to_pixels(c.Pos)
		is_moving := false

		if vm.stacking != nil && vm.stacking.card.Equals(c.Card) {
			p = p.Add2(vm.stacking.distance(vm))
			is_moving = true
		} else if vm.dragged_cards != nil {
			for _, d := range vm.dragged_cards.cards {
				if d.Equals(c.Card) {
					p = p.Add2(vm.cursor.Subtract2(vm.dragged_cards.offset))
					is_moving = true
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

		if is_moving {
			moving = append(moving, ImageData{image, p})
		} else {
			res = append(res, ImageData{image, p})
		}
	}

	res = append(res, moving...) // ensuring they're on top
	return res
}
