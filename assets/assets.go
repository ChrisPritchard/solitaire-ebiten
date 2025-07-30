package assets

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

// Card images: first index is suit (0 hearts, 1 diamonds, 2 spades, 3 clubs) second is value from 2 to 14 (ace)
var Cards map[int]map[int]*ebiten.Image

// Card back image
var CardBack *ebiten.Image

// Card space image
var CardSpace *ebiten.Image

// Button to reset the game
var ResetBtn *ebiten.Image

// Background image for the board
var Background *ebiten.Image

//go:embed felt_green.jpg
var background []byte

//go:embed card-suites.png
var card_data []byte

var card_size = image.Rect(0, 0, 36, 54)

//go:embed Audio/card-fan-1.ogg
var card_fan_1 []byte

//go:embed Audio/card-place-1.ogg
var card_place_1 []byte

//go:embed Audio/card-slide-1.ogg
var card_slide_1 []byte

//go:embed Audio/card-slide-2.ogg
var card_slide_2 []byte

//go:embed Audio/card-shove-3.ogg
var card_shove_3 []byte

func rel(dx, dy int) image.Rectangle {
	return card_size.Add(image.Pt(dx, dy))
}

var card_indexes = map[int]map[int]image.Rectangle{
	0: { // hearts
		2:  rel(22, 6),
		3:  rel(62, 6),
		4:  rel(102, 6),
		5:  rel(142, 6),
		6:  rel(182, 6),
		7:  rel(62, 64),
		8:  rel(102, 64),
		9:  rel(142, 64),
		10: rel(182, 64),
		11: rel(62, 122),
		12: rel(102, 122),
		13: rel(142, 122),
		14: rel(182, 122),
	},
	1: { // spades
		2:  rel(226, 300),
		3:  rel(266, 300),
		4:  rel(306, 300),
		5:  rel(346, 300),
		6:  rel(386, 300),
		7:  rel(226, 184),
		8:  rel(266, 184),
		9:  rel(306, 184),
		10: rel(346, 184),
		11: rel(226, 242),
		12: rel(266, 242),
		13: rel(306, 242),
		14: rel(346, 242),
	},
	2: { // diamonds
		2:  rel(22, 300),
		3:  rel(62, 300),
		4:  rel(102, 300),
		5:  rel(142, 300),
		6:  rel(182, 300),
		7:  rel(62, 184),
		8:  rel(102, 184),
		9:  rel(142, 184),
		10: rel(182, 184),
		11: rel(62, 242),
		12: rel(102, 242),
		13: rel(142, 242),
		14: rel(182, 242),
	},
	3: { // clubs
		2:  rel(226, 6),
		3:  rel(266, 6),
		4:  rel(306, 6),
		5:  rel(346, 6),
		6:  rel(386, 6),
		7:  rel(226, 64),
		8:  rel(266, 64),
		9:  rel(306, 64),
		10: rel(346, 64),
		11: rel(226, 122),
		12: rel(266, 122),
		13: rel(306, 122),
		14: rel(346, 122),
	},
}

var card_back = rel(12, 242)
var card_space = rel(12, 183)

var reset_btn = image.Rect(0, 0, 36, 15).Add(image.Pt(386, 256))

type SoundsGroup struct {
	NewGame   *audio.Player
	DragStart *audio.Player
	DragDrop  *audio.Player
	DrawDeck  *audio.Player
	Stack     *audio.Player
}

var Sounds SoundsGroup

func init() {
	load_background()
	load_cards()
	load_sounds()
}

func load_background() {
	background_image, _, err := image.Decode(bytes.NewReader(background))
	if err != nil {
		log.Fatal(err)
	}
	Background = ebiten.NewImageFromImage(background_image)
}

func load_cards() {
	cards_image, _, err := image.Decode(bytes.NewReader(card_data))
	if err != nil {
		log.Fatal(err)
	}
	cards_image2 := ebiten.NewImageFromImage(cards_image)

	CardBack = cards_image2.SubImage(card_back).(*ebiten.Image)
	CardSpace = cards_image2.SubImage(card_space).(*ebiten.Image)

	ResetBtn = cards_image2.SubImage(reset_btn).(*ebiten.Image)

	Cards = make(map[int]map[int]*ebiten.Image)
	for suit, v := range card_indexes {
		Cards[suit] = make(map[int]*ebiten.Image)
		for val, rect := range v {
			Cards[suit][val] = cards_image2.SubImage(rect).(*ebiten.Image)
		}
	}
}

func load_sounds() {
	audioContext := audio.NewContext(48000)

	load_sound := func(b []byte) *audio.Player {
		s, err := vorbis.DecodeF32(bytes.NewReader(b))
		if err != nil {
			log.Fatal(err)
		}

		p, err := audioContext.NewPlayerF32(s)
		if err != nil {
			log.Fatal(err)
		}
		return p
	}

	Sounds = SoundsGroup{
		DragStart: load_sound(card_place_1),
		DragDrop:  load_sound(card_slide_1),
		DrawDeck:  load_sound(card_slide_2),
		NewGame:   load_sound(card_fan_1),
		Stack:     load_sound(card_shove_3),
	}
}
