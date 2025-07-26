package main

type Vec2[T int | float64] struct {
	X, Y T
}

func (c Vec2[T]) Compare(o Vec2[T]) bool {
	if c.X == o.X {
		return c.Y < o.Y
	}
	return c.X < c.Y
}

func (c Vec2[T]) Equal(o Vec2[T]) bool {
	return c.X == o.X && c.Y == o.Y
}

func (c Vec2[T]) Scale(scale Vec2[T]) Vec2[T] {
	return Vec2[T]{c.X * scale.X, c.Y * scale.Y}
}

func (c Vec2[T]) Divide(scale Vec2[T]) Vec2[T] {
	return Vec2[T]{c.X / scale.X, c.Y / scale.Y}
}

func (c Vec2[T]) Contains(o Vec2[T], size Vec2[T]) bool {
	return c.X <= o.X && c.X+size.X >= o.X && c.Y <= o.Y && c.Y+size.Y >= o.Y
}

func (c Vec2[T]) Add(x, y T) Vec2[T] {
	return Vec2[T]{c.X + x, c.Y + y}
}

func (c Vec2[T]) Add2(o Vec2[T]) Vec2[T] {
	return Vec2[T]{c.X + o.X, c.Y + o.Y}
}

func (c Vec2[T]) Subtract(x, y T) Vec2[T] {
	return Vec2[T]{c.X - x, c.Y - y}
}

func (c Vec2[T]) Subtract2(o Vec2[T]) Vec2[T] {
	return Vec2[T]{c.X - o.X, c.Y - o.Y}
}

func (c Vec2[int]) ToFloat() Vec2[float64] {
	return Vec2[float64]{float64(c.X), float64(c.Y)}
}

func (c Vec2[float64]) ToInt() Vec2[int] {
	return Vec2[int]{int(c.X), int(c.Y)}
}
