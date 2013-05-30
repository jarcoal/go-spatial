package spatial

//
// Interfaces
//

type Seekable interface {
	Position() Point
	PreviousPosition() Point
}

type Seeker interface {
	ContainsPoint(Point) bool
	IntersectsRectangle(Rectangle) bool
}

//
// Constructor
//

func MakeQuadtree(depth, min, max int, bounds Rectangle) *Quadtree {
	return &Quadtree{
		Depth:   depth,
		Min:     min,
		Max:     max,
		Bounds:  bounds,
		Leaf:    true,
		Objects: make([]Seekable, 0),
	}
}

//
// Quadtree
//

type Quadtree struct {
	Depth    int
	Min, Max int
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

func (q *Quadtree) Remove(sk Seekable) {
	//delete the object from our local list
	for idx, qo := range q.Objects {
		if sk != qo {
			continue
		}

		//do the delete
		q.Objects = append(q.Objects[:idx], q.Objects[idx+1:]...)
		break
	}

	//one thing to keep in mind here is that a quadtree shouldn't
	//ever collapse itself because it can't know if it's siblings are
	//also ready to collapse.  only a parent can collapse it's children.

	if !q.Leaf {
		for _, subtree := range q.Subtrees {
			if !subtree.Bounds.ContainsPoint(sk.Position()) {
				continue
			}

			subtree.Remove(sk)

			//TODO: a way to allow the root node to collapse it's children.

			//if the subtree that we just removed an object from has 
			//dropped below the minimum, and also has children, collapse it
			if !subtree.Leaf && subtree.Min > len(subtree.Objects) {
				subtree.Collapse()
			}
		}
	}
}

func (q *Quadtree) Update(sk Seekable) {
	//we shouldn't reach the leaf nodes for this
	//operation, unless the root node is a leaf
	//in which case we wouldn't do anything
	if q.Leaf {
		return
	}

	//we're going to check to see if at any point the
	//new position diverges from the old one.
	var currentSubtree, nextSubtree *Quadtree

	for _, subtree := range q.Subtrees {
		//current
		if subtree.Bounds.ContainsPoint(sk.Position()) {
			currentSubtree = subtree
		}

		//next
		if subtree.Bounds.ContainsPoint(sk.PreviousPosition()) {
			nextSubtree = subtree
		}

		if currentSubtree != nil && nextSubtree != nil {
			break
		}
	}

	//if they've diverged, do a remove/add, otherwise
	//go down the tree further (assuming we're not at the bottom)
	if currentSubtree != nextSubtree {
		currentSubtree.Remove(sk)
		nextSubtree.Add(sk)
	} else if !currentSubtree.Leaf {
		currentSubtree.Update(sk)
	}

	//turns out the new position is in the
	//same leaf as before.  we do nothing.
}

func (q *Quadtree) Divide() {
	//log.Printf("%vDIVIDE: %v", strings.Repeat("-", q.Depth), q.Bounds)

	q.Leaf = false

	width := q.Bounds.Width / 2
	height := q.Bounds.Height / 2
	depth := q.Depth + 1

	//top left
	q.Subtrees[0] = MakeQuadtree(depth, q.Min, q.Max, Rectangle{
		Width:  width,
		Height: height,
		Center: Point{q.Bounds.Center.X - width/2, q.Bounds.Center.Y + height/2},
	})

	//top right
	q.Subtrees[1] = MakeQuadtree(depth, q.Min, q.Max, Rectangle{
		Width:  width,
		Height: height,
		Center: Point{q.Bounds.Center.X + width/2, q.Bounds.Center.Y + height/2},
	})

	//bottom left
	q.Subtrees[2] = MakeQuadtree(depth, q.Min, q.Max, Rectangle{
		Width:  width,
		Height: height,
		Center: Point{q.Bounds.Center.X - width/2, q.Bounds.Center.Y - height/2},
	})

	//bottom right
	q.Subtrees[3] = MakeQuadtree(depth, q.Min, q.Max, Rectangle{
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

func (q *Quadtree) Collapse() {
	q.Leaf = true

	for idx, _ := range q.Subtrees {
		q.Subtrees[idx] = nil
	}
}
