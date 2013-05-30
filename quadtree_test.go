package spatial

import (
	"math/rand"
	"testing"
)

//
// mock objects
//

type mockLocation struct {
	point         Point
	previousPoint Point
}

func (m *mockLocation) Position() Point {
	return m.point
}

func (m *mockLocation) PreviousPosition() Point {
	return m.previousPoint
}

func mockQuadtree() *Quadtree {
	bounds := Rectangle{Center: Point{0, 0}, Width: 50, Height: 50}
	return MakeQuadtree(0, 1, 5, bounds)
}

func randLocation() *mockLocation {
	return &mockLocation{
		point:         Point{rand.Float64(), rand.Float64()},
		previousPoint: Point{rand.Float64(), rand.Float64()},
	}
}

//
// tests
//

func TestQuadtreeAdd(t *testing.T) {
	qt := mockQuadtree()
	n := 5

	for i := 0; i < n; i++ {
		qt.Add(randLocation())
	}

	if len(qt.Objects) != n {
		t.FailNow()
	}
}

func TestQuadtreeRemove(t *testing.T) {
	qt := mockQuadtree()

	l := randLocation()
	qt.Add(l)
	qt.Remove(l)

	if len(qt.Objects) != 0 {
		t.FailNow()
	}
}

func TestQuadtreeUpdate(t *testing.T) {

}

func TestQuadtreeQuery(t *testing.T) {

}

func TestQuadtreeDivide(t *testing.T) {

}

func TestQuadtreeCollapse(t *testing.T) {

}
