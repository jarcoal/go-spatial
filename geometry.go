package spatial

import "math"

type Point struct {
	X, Y float64
}

type Rectangle struct {
	Center        Point
	Width, Height float64
}

func (r *Rectangle) ContainsPoint(p Point) bool {
	return p.X > r.Center.X-r.Height/2 && p.X < r.Center.X+r.Height/2 && p.Y > r.Center.Y-r.Height/2 && p.Y < r.Center.Y+r.Height/2
}

type Circle struct {
	Center Point
	Radius float64
}

func (c Circle) ContainsPoint(p Point) bool {
	return math.Pow(c.Center.X-p.X, 2)+math.Pow(c.Center.Y-p.Y, 2) <= math.Pow(c.Radius, 2)
}

func (c Circle) IntersectsRectangle(r Rectangle) bool {
	circleDistanceX := math.Abs(c.Center.X - r.Center.X)
	circleDistanceY := math.Abs(c.Center.Y - r.Center.Y)

	if circleDistanceX > r.Width/2+c.Radius || circleDistanceY > r.Height/2+c.Radius {
		return false
	}

	if circleDistanceX <= r.Width/2 || circleDistanceY <= r.Height/2 {
		return true
	}

	return math.Pow(circleDistanceX-r.Width/2, 2)+math.Pow(circleDistanceY-r.Height/2, 2) <= math.Pow(c.Radius, 2)
}
