package spatial

import "testing"

func mockQuadtree() *Quadtree {
	bounds := Rectangle{Center: Point{0, 0}, Width: 50, Height: 50}
	return MakeQuadtree(0, 5, bounds)
}

func TestQuadtreeAdd(t *testing.T) {

}

func TestQuadtreeQuery(t *testing.T) {

}

func TestQuadtreeMax(t *testing.T) {

}
