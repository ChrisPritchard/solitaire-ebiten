package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/chrispritchard/solitaire-ebiten/assets"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	cardScaling  = 2
	card_width   = 36. * cardScaling
	card_height  = 54. * cardScaling
)

type game_loop struct {
	pressed bool
}

type TouchState struct {
	Pressed bool
	X       float64
	Y       float64
}

var game SawayamaRules
var view_model ViewModel

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sawayama Solitaire")

	game = SawayamaRules{}
	view_model = ViewModel{CardWidth: card_width, CardHeight: card_height}

	if err := ebiten.RunGame(&game_loop{}); err != nil {
		log.Fatal(err)
	}
}

func (gl *game_loop) Update() error {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		gl.pressed = true
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		gl.pressed = false
	}
	x, y := ebiten.CursorPosition()
	touchState := TouchState{gl.pressed, float64(x), float64(y)}

	return view_model.Update(touchState, &game)
}

func (*game_loop) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.Background, nil)
	for _, image_data := range view_model.Transform(game) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(cardScaling, cardScaling)
		op.GeoM.Translate(image_data.X, image_data.Y)

		img := image_data.Image
		screen.DrawImage(img, op)
	}
}

func (*game_loop) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
