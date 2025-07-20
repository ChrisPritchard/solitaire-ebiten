package main

type Card struct {
	CUX, CUY    int
	Suit, Value int
	Visible     bool
}

type gamestate int

const (
	shuffling gamestate = iota
	dealing
)

type dragState struct {
	card    *Card
	offsetX int
	offsetY int
}

type SawayamaRules struct {
	Cards     []Card
	state     gamestate
	dragState dragState
}

func (r *SawayamaRules) Update(ts TouchState) error {
	switch r.state {
	case shuffling:
		r.Cards = shuffle_deck()
		r.state = dealing
	default:
		// if ts.Pressed && r.dragState.card == nil {
		// 	// detect if card under cursor
		// 	var card *Card = nil
		// 	for i := range r.cards {
		// 		c := &r.cards[i]
		// 		if c.X <= ts.X && c.Y <= ts.Y && c.X+c.Width >= ts.X && c.Y+c.Height >= ts.Y && (card == nil || card.Z < c.Z) {
		// 			card = c
		// 		}
		// 	}
		// 	if card != nil {
		// 		card.Z = 999999
		// 		r.dragState.card = card
		// 		r.dragState.offsetX = ts.X - card.X
		// 		r.dragState.offsetY = ts.Y - card.Y
		// 		card.Content.(*StandardDeck).Visible = true
		// 	}
		// } else if ts.Pressed {
		// 	r.dragState.card.X = ts.X - r.dragState.offsetX
		// 	r.dragState.card.Y = ts.Y - r.dragState.offsetY
		// } else if r.dragState.card != nil {
		// 	max := 0
		// 	for _, c := range r.cards {
		// 		if c.Z > max {
		// 			max = c.Z
		// 		}
		// 	}
		// 	r.dragState.card.Z = max + 1
		// 	// todo: place or return card
		// 	r.dragState.card = nil
		// }
	}
	return nil
}

func shuffle_deck() []Card {
	cards := make([]Card, 0)
	for i := range 4 {
		for j := range 13 {
			cards = append(cards, Card{
				CUX: i * 3, CUY: j,
				Suit: i, Value: j + 2,
				Visible: true,
			})
		}
	}

	// shuffled := make([]Card, 0)
	// for {
	// 	next := rand.Intn(len(cards))
	// 	shuffled = append(shuffled, cards[next])
	// 	if next == 0 {
	// 		cards = cards[1:]
	// 	} else {
	// 		cards = append(cards[:next-1], cards[next+1:]...)
	// 	}
	// 	if len(cards) == 0 {
	// 		break
	// 	}
	// }

	// return shuffled

	return cards
}
