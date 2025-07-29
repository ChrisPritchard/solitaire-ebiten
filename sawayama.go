package main

import (
	"log"
	"math/rand"
	"slices"
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
var Deck_CU = Vec2[int]{2 + CU_per_card.X, 1}

var waste_cu = Deck_CU.Add(CU_per_card.X+1, 0)
var pile_cus = []Vec2[int]{
	{2 + CU_per_card.X, 2 + CU_per_card.Y},
	{3 + 2*CU_per_card.X, 2 + CU_per_card.Y},
	{4 + 3*CU_per_card.X, 2 + CU_per_card.Y},
	{5 + 4*CU_per_card.X, 2 + CU_per_card.Y},
	{6 + 5*CU_per_card.X, 2 + CU_per_card.Y},
	{7 + 6*CU_per_card.X, 2 + CU_per_card.Y},
	{8 + 7*CU_per_card.X, 2 + CU_per_card.Y},
}
var foundation_cus = []Vec2[int]{
	{1, 1},
	{1, 2 + CU_per_card.Y},
	{1, 3 + 2*CU_per_card.Y},
	{1, 4 + 3*CU_per_card.Y},
}

type SawayamaRules struct {
	deck        []Card
	deck_space  *Card
	piles       [7][]Card
	foundations [4][]Card
	waste       []Card
}

func NewGameSetup() SawayamaRules {
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

	deck := shuffled[:c+1]

	return deck, piles
}

func (r *SawayamaRules) Cards() []CardInfo {
	res := []CardInfo{}

	if len(r.deck) > 0 {
		res = append(res, CardInfo{Pos: Deck_CU, Visible: false})
	} else if r.deck_space != nil {
		res = append(res, CardInfo{Card: *r.deck_space, Pos: Deck_CU, Visible: true})
	} else {
		res = append(res, CardInfo{Card: Card{}, Pos: Deck_CU, Visible: true})
	}

	for i, c := range r.waste {
		res = append(res, CardInfo{Card: c, Pos: waste_cu.Add(i, 0), Visible: true})
	}

	for i := range r.piles {
		if len(r.piles[i]) == 0 {
			res = append(res, CardInfo{Card: Card{}, Pos: pile_cus[i], Visible: true})
		} else {
			for j, c := range r.piles[i] {
				res = append(res, CardInfo{Card: c, Pos: pile_cus[i].Add(0, j), Visible: true})
			}
		}
	}

	for i, f := range r.foundations {
		if len(f) != 0 {
			res = append(res, CardInfo{Card: f[len(f)-1], Pos: foundation_cus[i], Visible: true})
		} else {
			res = append(res, CardInfo{Card: Card{}, Pos: foundation_cus[i], Visible: true})
		}
	}

	return res
}

func (r *SawayamaRules) DraggableAt(point Vec2[int]) ([]Card, Vec2[int]) {

	if len(r.deck) == 0 && r.deck_space != nil && Deck_CU.Contains(point, CU_per_card) {
		return []Card{*r.deck_space}, Deck_CU
	}

	if len(r.waste) > 0 && waste_cu.Add(len(r.waste)-1, 0).Contains(point, CU_per_card) {
		return []Card{r.waste[len(r.waste)-1]}, waste_cu.Add(len(r.waste)-1, 0)
	}

	for i, pile := range r.piles {
		if len(pile) == 0 {
			continue
		}
		if pile_cus[i].X <= point.X && pile_cus[i].X+CU_per_card.X >= point.X {
			for j := len(pile) - 1; j >= 0; j-- {
				if j != len(pile)-1 && (pile[j].Value != pile[j+1].Value+1 || pile[j].Suit%2 == pile[j+1].Suit%2) {
					return nil, Vec2[int]{} // stack isn't valid
				}
				cu := pile_cus[i].Add(0, j)
				if cu.Contains(point, CU_per_card) {
					return pile[j:], pile_cus[i].Add(0, j)
				}
			}
		}
	}

	return nil, Vec2[int]{}
}

func (r *SawayamaRules) DropAt(point Vec2[int], cards []Card, origin_cu Vec2[int]) {
	// deck space
	if len(cards) == 1 && len(r.deck) == 0 && r.deck_space == nil && Deck_CU.Contains(point, CU_per_card) {
		r.deck_space = &cards[0]
		r.remove_from_origin(cards, origin_cu)
		return
	}

	// foundations
	if len(cards) == 1 {
		for i, f := range foundation_cus {
			if f.Contains(point, CU_per_card) {
				var top_card *Card = nil
				if len(r.foundations[i]) != 0 {
					top_card = &(r.foundations[i][len(r.foundations[i])-1])
				}
				if top_card == nil && cards[0].Value == 14 { // ace
					r.foundations[i] = []Card{cards[0]}
					r.remove_from_origin(cards, origin_cu)
					return
				} else if top_card != nil && cards[0].Suit == top_card.Suit {
					if (top_card.Value == 14 && cards[0].Value == 2) || (top_card.Value == cards[0].Value-1) {
						r.foundations[i] = append(r.foundations[i], cards[0])
						r.remove_from_origin(cards, origin_cu)
						return
					}
				}
			}
		}
	}

	// piles
	for i, p := range pile_cus {
		if len(r.piles[i]) == 0 && p.Contains(point, CU_per_card) {
			r.piles[i] = cards
			r.remove_from_origin(cards, origin_cu)
			return
		}
		if p.Add(0, len(r.piles[i])-1).Contains(point, CU_per_card) {
			top_card := r.piles[i][len(r.piles[i])-1]
			if cards[0].Value == top_card.Value-1 && cards[0].Suit%2 != top_card.Suit%2 {
				r.piles[i] = append(r.piles[i], cards...)
				r.remove_from_origin(cards, origin_cu)
				return
			}
		}
	}
}

func (r *SawayamaRules) remove_from_origin(cards []Card, origin_cu Vec2[int]) {
	if origin_cu.Equal(Deck_CU) {
		r.deck_space = nil
		return
	}

	if origin_cu.Equal(waste_cu.Add(len(r.waste)-1, 0)) {
		r.waste = r.waste[:len(r.waste)-1]
		return
	}

	for i := range pile_cus {
		if origin_cu.X == pile_cus[i].X {
			r.piles[i] = r.piles[i][:len(r.piles[i])-len(cards)]
			return
		}
	}

	log.Fatal("reached this point, shouldn't be possible")
}

func (r *SawayamaRules) DrawFromDeck() bool {
	if len(r.deck) == 0 {
		return false
	}

	if len(r.deck) < 3 {
		slices.Reverse(r.deck)
		r.waste = append(r.waste, r.deck...)
		r.deck = nil
		return true
	}

	cards := r.deck[len(r.deck)-3:]
	slices.Reverse(cards)
	r.waste = append(r.waste, cards...)
	r.deck = r.deck[:len(r.deck)-3]
	return true
}

type Stackable struct {
	Card        Card
	Origin      Vec2[int]
	Destination Vec2[int]
}

func (r *SawayamaRules) NextStackable() *Stackable {
	var lfv = 14
	for _, f := range r.foundations {
		if len(f) > 0 {
			tc := f[len(f)-1]
			v := tc.Value
			if v == 14 {
				v = 1
			}
			if v < lfv {
				lfv = tc.Value
			}
		} else {
			lfv = 0
		}
	}
	lfv += 1
	if lfv == 1 {
		lfv = 14
	}

	for i, p := range r.piles {
		if len(p) > 0 {
			bc := p[len(p)-1]
			if bc.Value == lfv {
				for j, f := range r.foundations {
					if len(f) > 0 {
						tc := f[len(f)-1]
						if tc.Suit == bc.Suit && (tc.Value == bc.Value-1 || tc.Value == 14 && bc.Value == 2) {
							return &Stackable{bc, pile_cus[i].Add(0, len(p)-1), foundation_cus[j]}
						}
					} else if bc.Value == 14 {
						return &Stackable{bc, pile_cus[i].Add(0, len(p)-1), foundation_cus[j]}
					}
				}
			}
		}
	}

	if len(r.waste) > 0 && r.waste[len(r.waste)-1].Value == lfv {
		wc := r.waste[len(r.waste)-1]
		for j, f := range r.foundations {
			if len(f) > 0 {
				tc := f[len(f)-1]
				if tc.Suit == wc.Suit && (tc.Value == wc.Value-1 || tc.Value == 14 && wc.Value == 2) {
					return &Stackable{wc, waste_cu.Add(len(r.waste)-1, 0), foundation_cus[j]}
				}
			} else if wc.Value == 14 {
				return &Stackable{wc, waste_cu.Add(len(r.waste)-1, 0), foundation_cus[j]}
			}
		}
	}

	return nil
}
