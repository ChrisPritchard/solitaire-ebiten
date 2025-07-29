package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/chrispritchard/solitaire-ebiten/assets"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 768
	screenHeight = 560
	cardScaling  = 2
)

var card_size = Vec2[float64]{36. * cardScaling, 54. * cardScaling}

var reset_btn_size = Vec2[float64]{
	float64(assets.ResetBtn.Bounds().Dx() * cardScaling),
	float64(assets.ResetBtn.Bounds().Dy() * cardScaling),
}
var reset_btn_loc = Vec2[float64]{screenWidth - reset_btn_size.X - 20, 20}

type game_loop struct {
	pressed bool
}

type TouchState struct {
	Pressed     bool
	JustChanged bool
	Pos         Vec2[float64]
}

var game SawayamaRules
var view_model ViewModel

var last_pressed = false

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sawayama Solitaire")

	game = NewGameSetup()
	view_model = NewViewModel(card_size)

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
	touchState := TouchState{gl.pressed, gl.pressed != last_pressed, Vec2[int]{x, y}.ToFloat()}
	last_pressed = gl.pressed

	if touchState.Pressed && touchState.JustChanged && reset_btn_loc.Contains(touchState.Pos, reset_btn_size) {
		game = NewGameSetup()
		view_model = NewViewModel(card_size)
	} else {
		view_model.Update(touchState, &game)
	}

	return nil
}

func (*game_loop) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.Background, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(cardScaling, cardScaling)
	op.GeoM.Translate(reset_btn_loc.X, reset_btn_loc.Y)
	screen.DrawImage(assets.ResetBtn, op)

	for _, image_data := range view_model.Transform(game) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(cardScaling, cardScaling)
		op.GeoM.Translate(image_data.Pos.X, image_data.Pos.Y)

		img := image_data.Image
		screen.DrawImage(img, op)
	}
}

func (*game_loop) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
