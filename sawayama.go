package main

import (
	"math/rand"
	"slices"
	"sort"
)

type Card struct {
	Pos         Vec2[int]
	Suit, Value int
	Visible     bool
}

var CU_per_card = Vec2[int]{4, 4}
var Deck_CU = Vec2[int]{1, 1}

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
				Pos:  Deck_CU,
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

	pile_y := 2 + CU_per_card.Y
	for i := range 7 {
		pile_x := 1 + i*CU_per_card.X + i
		for j := range i + 1 {
			next := shuffled[c]
			next.Pos = Vec2[int]{pile_x, pile_y + j}
			next.Visible = true
			res = append(res, next)
			c--
		}
	}

	res = append(res, shuffled[:c]...)

	return res
}

func (r *SawayamaRules) Sort() {
	sort.Slice(r.Cards, func(i, j int) bool {
		return r.Cards[i].Pos.Compare(r.Cards[j].Pos)
	})
}

func (r *SawayamaRules) DraggableAt(point Vec2[int]) []*Card {
	if Deck_CU.Contains(point, CU_per_card) {
		return nil
	}

	for i := len(r.Cards) - 1; i >= 0; i-- {
		c := &r.Cards[i]
		if c.Pos.Contains(point, CU_per_card) {
			possible := []*Card{c}
			next := c
			for _, d := range r.Cards[i:] {
				if d.Pos.Equal(next.Pos.Add(0, 1)) {
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

	return nil
}

func (r *SawayamaRules) DropAt(point Vec2[int], cards []*Card) {
	for i, c := range slices.Backward(r.Cards) {
		if c.Pos.Contains(point, CU_per_card) {
			for _, d := range r.Cards[i:] {
				if d.Pos.Equal(c.Pos.Add(0, 1)) {
					return // not the top card, so can't drop here
					// todo, account for foundations and free spaces
				}
			}
			if c.Suit%2 != cards[0].Suit%2 && c.Value == cards[0].Value+1 {
				for i, d := range cards {
					d.Pos = c.Pos.Add(0, i+1)
				}
			}
			return
		}
	}
}

func (r *SawayamaRules) DrawFromDeck() bool {
	return false
}
