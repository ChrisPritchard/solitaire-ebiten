package main

import (
	"math/rand"
	"slices"
)

type Card struct {
	CUX, CUY    int
	Suit, Value int
	Visible     bool
}

const (
	CUX_per_card = 4
	CUY_per_card = 4
)

type gamestate int

const (
	shuffling gamestate = iota
	dealing
	ready
)

type SawayamaRules struct {
	Cards []Card
	state gamestate
}

func (r *SawayamaRules) Update() error {
	switch r.state {
	case shuffling:
		r.Cards = shuffle_deck()
		r.state = dealing
	case dealing:
		r.Cards = initial_deal(r.Cards)
		r.state = ready
	default:

	}
	return nil
}

func shuffle_deck() []Card {
	cards := make([]Card, 0)
	for i := range 4 {
		for j := range 13 {
			cards = append(cards, Card{
				CUX: 1, CUY: 1, // top-left deck position
				Suit: i, Value: j + 2,
				Visible: false,
			})
		}
	}

	shuffled := make([]Card, 0)
	for {
		next := rand.Intn(len(cards))
		shuffled = append(shuffled, cards[next])
		if next == 0 {
			cards = cards[1:]
		} else {
			cards = append(cards[:next], cards[next+1:]...)
		}
		if len(cards) == 0 {
			break
		}
	}

	return shuffled
}

func initial_deal(deck []Card) []Card {
	c := len(deck) - 1
	res := []Card{}

	pile_y := 2 + CUY_per_card
	for i := range 7 {
		pile_x := 1 + i*CUX_per_card + i
		for j := range i + 1 {
			next := deck[c]
			next.CUX = pile_x
			next.CUY = pile_y + j
			next.Visible = true
			res = append(res, next)
			c--
		}
	}

	res = append(res, deck[:c]...)

	return res
}

func (r *SawayamaRules) DraggableAt(cux, cuy int) []*Card {
	if r.state != ready {
		return nil // unlikely to ever happen
	}

	for i, c := range slices.Backward(r.Cards) {
		if cux >= c.CUX && cux <= c.CUX+CUX_per_card {
			if cuy >= c.CUY && cuy <= c.CUY+CUY_per_card {
				possible := []*Card{&c}
				next := &c
				for _, d := range r.Cards[i:] {
					if d.CUX == c.CUX && d.CUY == next.CUY+1 {
						if d.Value == next.Value-1 && d.Suit%2 != next.Suit%2 {
							possible = append(possible, &d)
							next = &d
						} else {
							return nil
						}
					}
				}
				return possible
			}
		}
	}

	return nil
}
