package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// Card images: first index is suit (0 hearts, 1 diamonds, 2 spades, 3 clubs) second is value from 2 to 14 (ace)
var Cards map[int]map[int]*ebiten.Image

// Card back image
var CardBack *ebiten.Image

// Background image for the board
var Background *ebiten.Image

//go:embed felt_green.jpg
var background []byte

//go:embed card-suites.png
var card_data []byte

var card_size = image.Rect(0, 0, 36, 54)

var card_indexes = map[int]map[int]image.Rectangle{
	0: map[int]image.Rectangle{ // hearts
		2:  card_size.Add(image.Pt(22, 6)),
		3:  card_size.Add(image.Pt(62, 6)),
		4:  card_size.Add(image.Pt(102, 6)),
		5:  card_size.Add(image.Pt(142, 6)),
		6:  card_size.Add(image.Pt(182, 6)),
		7:  card_size.Add(image.Pt(62, 64)),
		8:  card_size.Add(image.Pt(102, 64)),
		9:  card_size.Add(image.Pt(142, 64)),
		10: card_size.Add(image.Pt(182, 64)),
		11: card_size.Add(image.Pt(62, 122)),
		12: card_size.Add(image.Pt(102, 122)),
		13: card_size.Add(image.Pt(142, 122)),
		14: card_size.Add(image.Pt(182, 122)),
	},
	1: map[int]image.Rectangle{ // diamonds
		2:  card_size.Add(image.Pt(22, 300)),
		3:  card_size.Add(image.Pt(62, 300)),
		4:  card_size.Add(image.Pt(102, 300)),
		5:  card_size.Add(image.Pt(142, 300)),
		6:  card_size.Add(image.Pt(182, 300)),
		7:  card_size.Add(image.Pt(62, 184)),
		8:  card_size.Add(image.Pt(102, 184)),
		9:  card_size.Add(image.Pt(142, 184)),
		10: card_size.Add(image.Pt(182, 184)),
		11: card_size.Add(image.Pt(62, 242)),
		12: card_size.Add(image.Pt(102, 242)),
		13: card_size.Add(image.Pt(142, 242)),
		14: card_size.Add(image.Pt(182, 242)),
	},
	2: map[int]image.Rectangle{ // spades
		2:  card_size.Add(image.Pt(226, 300)),
		3:  card_size.Add(image.Pt(266, 300)),
		4:  card_size.Add(image.Pt(306, 300)),
		5:  card_size.Add(image.Pt(346, 300)),
		6:  card_size.Add(image.Pt(386, 300)),
		7:  card_size.Add(image.Pt(226, 184)),
		8:  card_size.Add(image.Pt(266, 184)),
		9:  card_size.Add(image.Pt(306, 184)),
		10: card_size.Add(image.Pt(346, 184)),
		11: card_size.Add(image.Pt(62, 242)),
		12: card_size.Add(image.Pt(102, 242)),
		13: card_size.Add(image.Pt(142, 242)),
		14: card_size.Add(image.Pt(182, 242)),
	},
	3: map[int]image.Rectangle{ // clubs
		2:  card_size.Add(image.Pt(226, 6)),
		3:  card_size.Add(image.Pt(266, 6)),
		4:  card_size.Add(image.Pt(306, 6)),
		5:  card_size.Add(image.Pt(346, 6)),
		6:  card_size.Add(image.Pt(386, 6)),
		7:  card_size.Add(image.Pt(226, 64)),
		8:  card_size.Add(image.Pt(266, 64)),
		9:  card_size.Add(image.Pt(306, 64)),
		10: card_size.Add(image.Pt(346, 64)),
		11: card_size.Add(image.Pt(62, 122)),
		12: card_size.Add(image.Pt(102, 122)),
		13: card_size.Add(image.Pt(142, 122)),
		14: card_size.Add(image.Pt(182, 122)),
	},
}

var card_back = card_size.Add(image.Pt(12, 242))

func init() {
	background_image, _, err := image.Decode(bytes.NewReader(background))
	if err != nil {
		log.Fatal(err)
	}
	Background = ebiten.NewImageFromImage(background_image)

	cards_image, _, err := image.Decode(bytes.NewReader(card_data))
	if err != nil {
		log.Fatal(err)
	}
	cards_image2 := ebiten.NewImageFromImage(cards_image)

	CardBack = cards_image2.SubImage(card_back).(*ebiten.Image)

	Cards = make(map[int]map[int]*ebiten.Image)
	for suit, v := range card_indexes {
		Cards[suit] = make(map[int]*ebiten.Image)
		for val, rect := range v {
			Cards[suit][val] = cards_image2.SubImage(rect).(*ebiten.Image)
		}
	}
}
