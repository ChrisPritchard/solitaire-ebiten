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
	cards []Card
	state gamestate
}

func (r *SawayamaRules) Update(touchState TouchState) error {
	switch r.state {
	case shuffling:
		r.cards = shuffle_deck()
		r.state = dealing
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
