package carrows

import "math"

type Rect struct {
	X, Y, Width, Height float64
}

type RectSide int

const (
	Top RectSide = iota
	Right
	Bottom
	Left
)

type Vec2 struct {
	X, Y float64
}

func distanceOf(p1, p2 Vec2) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}

func growRect(box Rect, size float64) Rect {
	return Rect{
		X:      box.X - size,
		Y:      box.Y - size,
		Width:  box.Width + 2*size,
		Height: box.Height + 2*size,
	}
}

func isPointInBox(p Vec2, b Rect) bool {
	return p.X > b.X && p.X < b.X+b.Width && p.Y > b.Y && p.Y < b.Y+b.Height
}

func controlPointOf(target, another Vec2, sideOfTarget RectSide, minDistanceToTarget float64) Vec2 {
	if minDistanceToTarget == 0 {
		minDistanceToTarget = 50
	}
	switch sideOfTarget {
	case Top:
		return Vec2{
			X: target.X,
			Y: math.Min((target.Y+another.Y)/2, target.Y-minDistanceToTarget),
		}
	case Bottom:
		return Vec2{
			X: target.X,
			Y: math.Max((target.Y+another.Y)/2, target.Y+minDistanceToTarget),
		}
	case Left:
		return Vec2{
			X: math.Min((target.X+another.X)/2, target.X-minDistanceToTarget),
			Y: target.Y,
		}
	case Right:
		return Vec2{
			X: math.Max((target.X+another.X)/2, target.X+minDistanceToTarget),
			Y: target.Y,
		}
	default:
		panic("unreachable")
	}
}

func angleOf(enteringSide RectSide) int {
	switch enteringSide {
	case Left:
		return 0
	case Top:
		return 90
	case Right:
		return 180
	case Bottom:
		return 270
	default:
		panic("unreachable")
	}
}
