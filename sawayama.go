package main

import (
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

type SawayamaRules struct{}

func (r *SawayamaRules) Update() error { return nil }

func (r *SawayamaRules) Cards() []Card {

	cards := make([]Card, 0)
	for i := range 4 {
		for j := range 13 {
			cards = append(cards, Card{
				X: i*120 + j*20, Y: j * 20, Z: j,
				Content: &StandardDeck{Suit: i, Value: j + 2, Visible: j%2 == 0},
			})
		}
	}

	return cards
}
