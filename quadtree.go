package spatial

//
// Interfaces
//

type Seekable interface {
	Position() Point
}

type Seeker interface {
	ContainsPoint(Point) bool
	IntersectsRectangle(Rectangle) bool
}

//
// Constructor
//

func MakeQuadtree(depth, max int, bounds Rectangle) *Quadtree {
	return &Quadtree{
		Depth:  depth,
		Max:    max,
		Bounds: bounds,
		Leaf:   true,
	}
}

//
// Quadtree
//

type Quadtree struct {
	Depth    int
	Max      int
	Bounds   Rectangle
	Objects  []Seekable
	Subtrees [4]*Quadtree
	Leaf     bool
}

func (q *Quadtree) Query(sk Seeker) []Seekable {
	if q.Leaf {
		return q.QuerySelf(sk)
	}

	return q.QuerySubtrees(sk)
}

func (q *Quadtree) QuerySelf(sk Seeker) []Seekable {
	//log.Printf("%vQUERY SELF: %v", strings.Repeat("-", q.Depth), q.Bounds)

	//this will be a list of objects that are contained by the sk object
	results := make([]Seekable, 0)

	//check each point to see if it's contained by our sk object
	for _, obj := range q.Objects {
		if !sk.ContainsPoint(obj.Position()) {
			continue
		}

		results = append(results, obj)
	}

	return results
}

func (q *Quadtree) QuerySubtrees(sk Seeker) []Seekable {
	//log.Printf("%vQUERY SUBTREE: %v", strings.Repeat("-", q.Depth), q.Bounds)

	results := make([]Seekable, 0)

	//dispatch the sk to any relevant subtrees
	for _, subtree := range q.Subtrees {
		//our search radius doesn't intersect this subtree, skip it
		if !sk.IntersectsRectangle(subtree.Bounds) {
			//log.Printf("%vIGNORING SUBTREE: %v", strings.Repeat("-", subtree.Depth), subtree.Bounds)
			continue
		}

		//log.Printf("%vINVESTIGATING SUBTREE: %v", strings.Repeat("-", subtree.Depth), subtree.Bounds)
		for _, result := range subtree.Query(sk) {
			results = append(results, result)
		}
	}

	return results
}

func (q *Quadtree) Add(obj Seekable) {
	q.Objects = append(q.Objects, obj)

	if q.Leaf && len(q.Objects) > q.Max {
		q.Divide()
	} else if !q.Leaf {
		for _, subtree := range q.Subtrees {
			if !subtree.Bounds.ContainsPoint(obj.Position()) {
				continue
			}

			subtree.Add(obj)
			break
		}
	}
}

func (q *Quadtree) Divide() {
	//log.Printf("%vDIVIDE: %v", strings.Repeat("-", q.Depth), q.Bounds)

	q.Leaf = false

	width := q.Bounds.Width / 2
	height := q.Bounds.Height / 2
	depth := q.Depth + 1

	//top left
	q.Subtrees[0] = MakeQuadtree(depth, q.Max, Rectangle{
		Width:  width,
		Height: height,
		Center: Point{q.Bounds.Center.X - width/2, q.Bounds.Center.Y + height/2},
	})

	//top right
	q.Subtrees[1] = MakeQuadtree(depth, q.Max, Rectangle{
		Width:  width,
		Height: height,
		Center: Point{q.Bounds.Center.X + width/2, q.Bounds.Center.Y + height/2},
	})

	//bottom left
	q.Subtrees[2] = MakeQuadtree(depth, q.Max, Rectangle{
		Width:  width,
		Height: height,
		Center: Point{q.Bounds.Center.X - width/2, q.Bounds.Center.Y - height/2},
	})

	//bottom right
	q.Subtrees[3] = MakeQuadtree(depth, q.Max, Rectangle{
		Width:  width,
		Height: height,
		Center: Point{q.Bounds.Center.X + width/2, q.Bounds.Center.Y - height/2},
	})

	for _, obj := range q.Objects {
		for _, subtree := range q.Subtrees {
			if !subtree.Bounds.ContainsPoint(obj.Position()) {
				continue
			}

			subtree.Add(obj)
			break
		}
	}
}
