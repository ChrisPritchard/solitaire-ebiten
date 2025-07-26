package main

import (
	"math/rand"
)

type Card struct {
	Suit, Value int
}

func (c Card) Equals(o Card) bool {
	return c.Suit == o.Suit && c.Value == o.Value
}

type CardInfo struct {
	Card
	Pos     Vec2[int]
	Visible bool
}

var CU_per_card = Vec2[int]{4, 4}
var Deck_CU = Vec2[int]{1, 1}
var Pile_CUs = []Vec2[int]{
	Vec2[int]{1, 2 + CU_per_card.Y},
	Vec2[int]{2 + CU_per_card.X, 2 + CU_per_card.Y},
	Vec2[int]{3 + 2*CU_per_card.X, 2 + CU_per_card.Y},
	Vec2[int]{4 + 3*CU_per_card.X, 2 + CU_per_card.Y},
	Vec2[int]{5 + 4*CU_per_card.X, 2 + CU_per_card.Y},
	Vec2[int]{6 + 5*CU_per_card.X, 2 + CU_per_card.Y},
	Vec2[int]{7 + 6*CU_per_card.X, 2 + CU_per_card.Y},
}

type SawayamaRules struct {
	deck        []Card
	deck_space  *Card
	piles       [7][]Card
	foundations [4][]Card
	waste       []Card
}

func Setup() SawayamaRules {
	shuffled := shuffle_deck()
	r := SawayamaRules{}
	r.deck, r.piles = initial_deal(shuffled)
	return r
}

func shuffle_deck() []Card {
	cards := make([]Card, 0)
	for i := range 4 {
		for j := range 13 {
			cards = append(cards, Card{i, j + 2})
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

func initial_deal(shuffled []Card) ([]Card, [7][]Card) {
	c := len(shuffled) - 1

	var piles [7][]Card

	for i := range 7 {
		piles[i] = []Card{}
		for range i + 1 {
			piles[i] = append(piles[i], shuffled[c])
			c--
		}
	}

	deck := shuffled[:c]

	return deck, piles
}

func (r *SawayamaRules) Cards() []CardInfo {
	res := []CardInfo{}

	if len(r.deck) > 0 {
		res = append(res, CardInfo{Pos: Deck_CU, Visible: false})
	}

	for i := range r.piles {
		for j, c := range r.piles[i] {
			res = append(res, CardInfo{
				Card:    c,
				Pos:     Pile_CUs[i].Add(0, j),
				Visible: true})
		}
	}

	return res
}

func (r *SawayamaRules) DraggableAt(point Vec2[int]) []Card {
	if Deck_CU.Contains(point, CU_per_card) {
		return nil
	}

	for i, pile := range r.piles {
		if len(pile) == 0 {
			continue
		}
		if Pile_CUs[i].X <= point.X && Pile_CUs[i].X+CU_per_card.X >= point.X {
			for j := len(pile) - 1; j >= 0; j-- {
				cu := Pile_CUs[i].Add(0, j)
				if cu.Contains(point, CU_per_card) {
					return []Card{pile[j]}
				}
			}
		}
	}

	// for i := len(r.Cards) - 1; i >= 0; i-- {
	// 	c := &r.Cards[i]
	// 	if c.Pos.Contains(point, CU_per_card) {
	// 		possible := []*Card{c}
	// 		next := c
	// 		for _, d := range r.Cards[i:] {
	// 			if d.Pos.Equal(next.Pos.Add(0, 1)) {
	// 				if d.Value == next.Value-1 && d.Suit%2 != next.Suit%2 {
	// 					possible = append(possible, &d)
	// 					next = &d
	// 				} else {
	// 					return nil
	// 				}
	// 			}
	// 		}
	// 		return possible
	// 	}
	// }

	return nil
}

func (r *SawayamaRules) DropAt(point Vec2[int], cards []Card) {
	// for i, c := range slices.Backward(r.Cards) {
	// 	if c.Pos.Contains(point, CU_per_card) {
	// 		for _, d := range r.Cards[i:] {
	// 			if d.Pos.Equal(c.Pos.Add(0, 1)) {
	// 				return // not the top card, so can't drop here
	// 				// todo, account for foundations and free spaces
	// 			}
	// 		}
	// 		if c.Suit%2 != cards[0].Suit%2 && c.Value == cards[0].Value+1 {
	// 			for i, d := range cards {
	// 				d.Pos = c.Pos.Add(0, i+1)
	// 			}
	// 		}
	// 		return
	// 	}
	// }
}

func (r *SawayamaRules) DrawFromDeck() bool {
	return false
}
