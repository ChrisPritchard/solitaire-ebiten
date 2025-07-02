package main

import (
	"math/rand"

	"github.com/chrispritchard/solitaire-ebiten/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

type StandardDeck struct {
	Suit, Value int
	Visible     bool
}

func (c *StandardDeck) Image() *ebiten.Image {
	if !c.Visible {
		return assets.CardBack
	}

	return assets.Cards[c.Suit][c.Value]
}

type gamestate int

const (
	shuffling gamestate = iota
	dealing
)

type SawayamaRules struct {
	cards   []Card
	state   gamestate
	dragged *Card
}

func (r *SawayamaRules) Update(ts TouchState) error {
	switch r.state {
	case shuffling:
		r.cards = shuffle_deck()
		r.state = dealing
	default:
		if ts.Pressed && r.dragged == nil {
			// detect if card under cursor
			var card *Card = nil
			for i := range r.cards {
				c := &r.cards[i]
				if c.X <= ts.X && c.Y <= ts.Y && c.X+c.Width >= ts.X && c.Y+c.Height >= ts.Y {
					if card == nil || card.Z < c.Z {
						card = c
					}
				}
			}
			r.dragged = card
			card.Content.(*StandardDeck).Visible = true
		} else if ts.Pressed {
			r.dragged.X = ts.X
			r.dragged.Y = ts.Y // todo: should adjust for offset
		} else if r.dragged != nil {
			// place or return card
			r.dragged = nil
		}
	}
	return nil
}

func (r *SawayamaRules) Cards() []Card {

	return r.cards
}

func shuffle_deck() []Card {
	cards := make([]Card, 0)
	for i := range 4 {
		for j := range 13 {
			cards = append(cards, Card{
				X: 0, Y: 0, Z: 0,
				Width: 36 * 3, Height: 54 * 3, // todo, fix
				Content: &StandardDeck{Suit: i, Value: j + 2, Visible: false},
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
			cards = append(cards[:next-1], cards[next+1:]...)
		}
		if len(cards) == 0 {
			break
		}
	}

	return shuffled
}
