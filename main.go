package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"sort"

	"github.com/chrispritchard/solitaire-ebiten/assets"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sawayama Solitaire")

	game := NewGame(&SawayamaRules{})

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

type Card struct {
	X, Y, Z       int
	Width, Height int
	Content       CardContent
}

type CardContent interface {
	Image() *ebiten.Image
}

type TouchState struct {
	Pressed bool
	X       int
	Y       int
}

type RuleSet interface {
	Update(TouchState) error
	Cards() []Card
}

type Game struct {
	rules   RuleSet
	pressed bool
}

func NewGame(rules RuleSet) Game {
	return Game{
		rules:   rules,
		pressed: false,
	}
}

func (game *Game) Update() error {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		game.pressed = true
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		game.pressed = false
	}
	x, y := ebiten.CursorPosition()
	touchState := TouchState{game.pressed, x, y}

	return game.rules.Update(touchState)
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.Background, nil)

	cards := game.rules.Cards()
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Z < cards[j].Z
	})

	for _, card := range cards {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(3, 3)
		op.GeoM.Translate(float64(card.X), float64(card.Y))

		img := card.Content.Image()
		screen.DrawImage(img, op)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
