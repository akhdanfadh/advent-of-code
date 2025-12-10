package main

import (
	"bufio"
	"fmt"
	"os"
)

type line struct {
	p1, p2 tile
	isVert bool
}

func processV2b(filename string) (string, error) {
	// same idea as V2, ray casting
	// what's different:
	// - in V2, due to sampling, we may miss small invalid regions inside large rectangles
	// - in V2, we apply ray casting at the tile level and use sampling to reduce the number of tiles checked
	//   i.e., this tile is OK if it is red, green, or inside the polygon
	// - here, we reason with the rectangle as a whole and only ever touches the polygon edges, no inner tiles
	//   i.e., this rectangle is valid if all corners are inside/boundary and no edge crosses the rectangle interior
	// - here, we exploit the fact that the polygon is axis-aligned (only vertical/horizontal edges),
	//   due to that, we can just check based on endpoints/corners, no need the inner points

	corners, err := getCorners(filename)
	if err != nil {
		return "", err
	}
	edges, err := buildEdges(corners)
	if err != nil {
		return "", err
	}

	maxArea := uint(0)
	for i := 0; i < len(corners); i++ {
		for j := i + 1; j < len(corners); j++ {
			if !isRectangleValid(corners[i], corners[j], edges) {
				continue
			}
			area := calcArea(corners[i], corners[j])
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return fmt.Sprintf("Largest rectangle area: %d", maxArea), nil
}

func getCorners(filename string) ([]tile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	corners := []tile{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		t, err := parseTile(line)
		if err != nil {
			return nil, err
		}
		corners = append(corners, t)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return corners, nil
}

func buildEdges(corners []tile) ([]line, error) {
	n := len(corners)
	edges := make([]line, n)
	for i := range n {
		j := (i + 1) % n // wrap around
		p1, p2 := corners[i], corners[j]
		if p1.x == p2.x { // vertical line
			edges[i] = line{p1: p1, p2: p2, isVert: true}
		} else if p1.y == p2.y { // horizontal line
			edges[i] = line{p1: p1, p2: p2, isVert: false}
		} else {
			return nil, fmt.Errorf("non-axis-aligned edge between %v and %v", p1, p2)
		}
	}
	return edges, nil
}

func isInsideOrBoundary(t tile, edges []line) bool {
	// classic ray casting algorithm:
	// - for a point (x, y), cast a ray to the right (increasing x)
	// - count how many times it intersects polygon edges
	// - if the count is odd, the point is inside, otherwise outside
	// - additionally, if the point lies exactly on an edge, consider it inside
	// also this code assumes the polygon is axis-aligned (only vertical/horizontal edges)
	x, y := t.x, t.y

	// 1. boundary check (on corners or edges)
	//
	// for each edge from (x1, y1) to (x2, y2):
	// - if vertical, (x, y) on this edge if x == x1 and ymin <= y <= ymax
	// - if horizontal, (x, y) on this edge if y == y1 and xmin <= x <= xmax
	for _, e := range edges {
		x1, y1 := e.p1.x, e.p1.y
		x2, y2 := e.p2.x, e.p2.y
		if e.isVert {
			ylow, yhigh := min(y1, y2), max(y1, y2)
			if x == x1 && ylow <= y && y <= yhigh {
				return true
			}
		} else {
			xlow, xhigh := min(x1, x2), max(x1, x2)
			if y == y1 && xlow <= x && x <= xhigh {
				return true
			}
		}
	}

	// 2. ray casting
	//
	// our ray is horizontal line from (x,y) to the right.
	// for each vertical edge at x = xe with ylow to yhigh,
	// ray intersects this edge if:
	// - xe > x (edge is to the right of point), and
	// - ylow < y < yhigh (edge spans the y coordinate of the point)
	crossings := 0
	for _, e := range edges {
		if !e.isVert {
			continue
		} // only vertical edges matter for our horizontal ray castingj

		xe := e.p1.x // == e.p2.x
		ylow, yhigh := min(e.p1.y, e.p2.y), max(e.p1.y, e.p2.y)
		if xe > x && ylow < y && y < yhigh {
			crossings++
		}
	}
	return crossings%2 == 1 // if odd, inside
}

func isRectangleValid(t1, t2 tile, edges []line) bool {
	// this code assumes the polygon is axis-aligned (only vertical/horizontal edges)

	x1, y1 := t1.x, t1.y
	x2, y2 := t2.x, t2.y
	xmin, xmax := min(x1, x2), max(x1, x2)
	ymin, ymax := min(y1, y2), max(y1, y2)

	// case 1: thin horizontal rectangle
	// check that no vertical edge crosses the open segment (xmin, xmax)
	//
	// for each vertical edge at x = xe with ylow to yhigh,
	// if xmin < xe < xmax and ylow < y < yhigh, then this edge intersects
	// the interior of the horizontal segment, meaning we'd toggle from
	// inside to outside somewhere (no need explicit ray casting)
	if y1 == y2 {
		y := y1
		for _, e := range edges {
			if !e.isVert {
				continue
			}
			xe := e.p1.x // == e.p2.x
			ylow, yhigh := min(e.p1.y, e.p2.y), max(e.p1.y, e.p2.y)
			if xmin < xe && xe < xmax && ylow < y && y < yhigh {
				return false
			}
		}
	}

	// case 2: thin vertical rectangle
	// same logic as case 1, but swap x and y
	if x1 == x2 {
		x := x1
		for _, e := range edges {
			if e.isVert {
				continue
			}
			ye := e.p1.y
			xlow, xhigh := min(e.p1.x, e.p2.x), max(e.p1.x, e.p2.x)
			if ymin < ye && ye < ymax && xlow < x && x < xhigh {
				return false
			}
		}
	}

	// case 3: non-thin rectangle
	// check 3a: the other two corners are inside or on the boundary
	c1 := tile{x1, y2}
	c2 := tile{x2, y1}
	if !isInsideOrBoundary(c1, edges) || !isInsideOrBoundary(c2, edges) {
		return false
	}

	// check 3b: no edge intersects the open interior of the rectangle (xmin, xmax) x (ymin, ymax)
	//
	// for each vertical edge at x = xe with ylow to yhigh,
	// the edge intersects the interior of the rectangle if
	// xmin < xe < xmax and max(ylow, ymin) < min(yhigh, ymax)
	//
	// for each horizontal edge at y = ye with xlow to xhigh,
	// the edge intersects the interior of the rectangle if
	// ymin < ye < ymax and max(xlow, xmin) < min(xhigh, xmax)
	for _, e := range edges {
		ex1, ey1 := e.p1.x, e.p1.y
		ex2, ey2 := e.p2.x, e.p2.y

		if e.isVert {
			if !(xmin < ex1 && ex1 < xmax) {
				continue
			}
			ylow, yhigh := min(ey1, ey2), max(ey1, ey2)
			if max(ylow, ymin) < min(yhigh, ymax) {
				return false
			}

		} else {
			if !(ymin < ey1 && ey1 < ymax) {
				continue
			}
			xlow, xhigh := min(ex1, ex2), max(ex1, ex2)
			if max(xlow, xmin) < min(xhigh, xmax) {
				return false
			}
		}
	}

	return true
}
