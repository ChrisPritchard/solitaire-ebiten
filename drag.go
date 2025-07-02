package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// code taken from https://github.com/hajimehoshi/ebiten/blob/main/examples/drag/main.go with some small modifications

func Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if sp := g.spriteAt(ebiten.CursorPosition()); sp != nil {
			s := NewStroke(&MouseStrokeSource{}, sp)
			g.strokes[s] = struct{}{}
			g.moveSpriteToFront(sp)
		}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		if sp := g.spriteAt(ebiten.TouchPosition(id)); sp != nil {
			s := NewStroke(&TouchStrokeSource{id}, sp)
			g.strokes[s] = struct{}{}
			g.moveSpriteToFront(sp)
		}
	}

	for s := range g.strokes {
		s.Update()
		if !s.sprite.dragged {
			delete(g.strokes, s)
		}
	}
}

type Movable interface {
	SetPosition(int, int)
	GetPosition() (int, int)
	GetDragState() bool
	SetDragState(bool)
}

// StrokeSource represents a input device to provide strokes.
type StrokeSource interface {
	Position() (int, int)
	IsJustReleased() bool
}

// MouseStrokeSource is a StrokeSource implementation of mouse.
type MouseStrokeSource struct{}

func (m *MouseStrokeSource) Position() (int, int) {
	return ebiten.CursorPosition()
}

func (m *MouseStrokeSource) IsJustReleased() bool {
	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
}

// TouchStrokeSource is a StrokeSource implementation of touch.
type TouchStrokeSource struct {
	ID ebiten.TouchID
}

func (t *TouchStrokeSource) Position() (int, int) {
	return ebiten.TouchPosition(t.ID)
}

func (t *TouchStrokeSource) IsJustReleased() bool {
	return inpututil.IsTouchJustReleased(t.ID)
}

// Stroke manages the current drag state by mouse.
type Stroke struct {
	source StrokeSource

	// offsetX and offsetY represents a relative value from the sprite's upper-left position to the cursor position.
	offsetX int
	offsetY int

	target *Movable
}

func NewStroke(source StrokeSource, target *Movable) *Stroke {
	(*target).SetDragState(true)
	sx, sy := (*target).GetPosition()
	x, y := source.Position()
	return &Stroke{
		source:  source,
		offsetX: x - sx,
		offsetY: y - sy,
		target:  target,
	}
}

func (s *Stroke) Update() {
	if !(*s.target).GetDragState() {
		return
	}
	if s.source.IsJustReleased() {
		(*s.target).SetDragState(false)
		return
	}

	x, y := s.source.Position()
	x -= s.offsetX
	y -= s.offsetY
	(*s.target).SetPosition(x, y)
}
