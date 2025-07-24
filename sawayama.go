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

type SawayamaRules struct {
	Cards []Card
}

func Setup() SawayamaRules {
	shuffled := shuffle_deck()
	r := SawayamaRules{}
	r.Cards = initial_deal(shuffled)
	return r
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

func initial_deal(shuffled []Card) []Card {
	c := len(shuffled) - 1
	res := []Card{}

	pile_y := 2 + CUY_per_card
	for i := range 7 {
		pile_x := 1 + i*CUX_per_card + i
		for j := range i + 1 {
			next := shuffled[c]
			next.CUX = pile_x
			next.CUY = pile_y + j
			next.Visible = true
			res = append(res, next)
			c--
		}
	}

	res = append(res, shuffled[:c]...)

	return res
}

func (r *SawayamaRules) DraggableAt(cux, cuy int) []*Card {
	for i := len(r.Cards) - 1; i >= 0; i-- {
		c := &r.Cards[i]
		if cux >= c.CUX && cux <= c.CUX+CUX_per_card {
			if cuy >= c.CUY && cuy <= c.CUY+CUY_per_card {
				possible := []*Card{c}
				next := c
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

func (r *SawayamaRules) DroppableAt(cux, cuy int, suit, value int) (bool, int, int) {
	for i, c := range slices.Backward(r.Cards) {
		if cux >= c.CUX && cux <= c.CUX+CUX_per_card {
			if cuy >= c.CUY && cuy <= c.CUY+CUY_per_card {
				for _, d := range r.Cards[i:] {
					if d.CUX == c.CUX && d.CUY == c.CUY+1 {
						return false, 0, 0 // not the top card, so can't drop here
						// todo, account for foundations and free spaces
					}
				}
				return c.Suit%2 != suit%2 && c.Value == value+1, c.CUX, c.CUY
			}
		}
	}

	return false, 0, 0
}

func (r *SawayamaRules) DrawFromDeck() bool {
	return false
}
