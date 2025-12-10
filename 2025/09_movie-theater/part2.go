package main

import (
	"bufio"
	"fmt"
	"os"
)

func processV2(filename string, sampleSize int) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// idea: ray casting -> i'll admit i'm asking youtube for this :'(
	// https://www.youtube.com/watch?v=RyLuE5xFLxw
	// - first get the red tiles (they form a closed loop in order)
	// - for all possible rectangles, check validity by using ray casting
	//   to check if points are inside the polygon formed by red tiles
	// - to optimize large sparse rectangles, use sampling instead of checking every tile

	// read red tiles, they form a loop (next is always adjacent to previous)
	polygonCorners := []tile{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		t, err := parseTile(line)
		if err != nil {
			return "", err
		}
		polygonCorners = append(polygonCorners, t)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	isRedTile := make(map[tile]bool)
	for _, t := range polygonCorners {
		isRedTile[t] = true
	}

	// get grren tile boundaries
	isGreenTile := make(map[tile]bool)
	for i := 0; i < len(polygonCorners); i++ {
		t1 := polygonCorners[i]
		t2 := polygonCorners[(i+1)%len(polygonCorners)] // wrap around

		// add all tiles between t1 and t2 (exclusive)
		if t1.x == t2.x { // same column, fill vertically
			minY, maxY := min(t1.y, t2.y), max(t1.y, t2.y)
			for y := minY; y <= maxY; y++ {
				t := tile{t1.x, y}
				if !isRedTile[t] {
					isGreenTile[t] = true
				}
			}
		} else if t1.y == t2.y { // same row, fill horizontally
			minX, maxX := min(t1.x, t2.x), max(t1.x, t2.x)
			for x := minX; x <= maxX; x++ {
				t := tile{x, t1.y}
				if !isRedTile[t] {
					isGreenTile[t] = true
				}
			}
		}
	}

	// check all pairs of red tiles
	largestArea := uint(0)
	for i := 0; i < len(polygonCorners); i++ {
		for j := i + 1; j < len(polygonCorners); j++ {
			t1, t2 := polygonCorners[i], polygonCorners[j]
			area := calcArea(t1, t2)

			// early skip
			if area <= largestArea {
				continue
			}
			fmt.Printf("Checking tiles %v and %v with area %d\n", t1, t2, area)

			// check if all tiles in rectangle are either red, green, or inside polygon
			// for large rectangles, use sampling to avoid checking millions of points inside it
			minX, maxX := min(t1.x, t2.x), max(t1.x, t2.x)
			maxY, minY := max(t1.y, t2.y), min(t1.y, t2.y)
			width, height := maxX-minX+1, maxY-minY+1
			if width > 1000 || height > 1000 {
				if isValidRectangleSampled(uint(sampleSize), minX, maxX, minY, maxY, polygonCorners, isRedTile, isGreenTile) {
					largestArea = area
				}
			} else {
				valid := true
				for x := minX; x <= maxX && valid; x++ {
					for y := minY; y <= maxY; y++ {
						if !isTileValid(tile{x, y}, polygonCorners, isRedTile, isGreenTile) {
							valid = false
							fmt.Printf("  Invalid tile found at %d,%d\n", x, y)
							break
						}
					}
				}
				if valid {
					largestArea = area
				}

			}
		}
	}

	return fmt.Sprintf("Largest rectangle area: %d", largestArea), nil
}

func isTileValid(t tile, polygonCorners []tile, isRedTile, isGreenTile map[tile]bool) bool {
	return isRedTile[t] || isGreenTile[t] || isInsidePolygon(t, polygonCorners)
}

func isValidRectangleSampled(sampleSize, minX, maxX, minY, maxY uint, polygonCorners []tile, isRedTile, isGreenTile map[tile]bool) bool {
	// sample at intervals
	width, height := maxX-minX+1, maxY-minY+1
	stepX := max(width/sampleSize, 1)
	stepY := max(height/sampleSize, 1)

	// check all four corners first
	corners := []tile{
		{minX, minY}, {maxX, minY}, {minX, maxY}, {maxX, maxY},
	}
	for _, t := range corners {
		if !isTileValid(t, polygonCorners, isRedTile, isGreenTile) {
			return false
		}
	}

	// now check all tiles at intervals
	for x := minX; x <= maxX; x += stepX {
		for y := minY; y <= maxY; y += stepY {
			if !isTileValid(tile{x, y}, polygonCorners, isRedTile, isGreenTile) {
				return false
			}
		}
	}
	return true
}

// ray casting algorithm to determine if point is inside polygon
// it casts a horizontal ray from the point to the right and counts edge crossings
// odd crossings = inside, even crossings = outside
func isInsidePolygon(t tile, polygonCorners []tile) bool {
	isInside := false

	// iterate over each edge build by consecutive red tiles
	for i := 0; i < len(polygonCorners); i++ {
		c1 := polygonCorners[i]
		c2 := polygonCorners[(i+1)%len(polygonCorners)] // wrap around

		// base 1: point lies exactly on vertical edge
		if c1.x == c2.x && c1.x == t.x { // c1 and c2 and t share same x
			if t.y >= min(c1.y, c2.y) && t.y <= max(c1.y, c2.y) {
				return true
			}
		}
		// base 2: point lies exactly on horizontal edge
		if c1.y == c2.y && c1.y == t.y { // c1 and c2 and t share same y
			if t.x >= min(c1.x, c2.x) && t.x <= max(c1.x, c2.x) {
				return true
			}
		}

		// ray casting logic: check if horizontal ray from tile (going right) intersects edge
		// use > for min and <= for max to handle vertex crossings consistently
		// cond. 1: ray's y-coordinate must be within edge's y-range
		if t.y > min(c1.y, c2.y) && t.y <= max(c1.y, c2.y) {
			// cond. 2: tile must be to the left of (or at) the rightmost edge
			if t.x <= max(c1.x, c2.x) {
				if c1.y != c2.y { // skip horizontal edges (already handled above)
					// get x where ray intersect edge
					// using line equation: x = x1 + (y - y1) * (x2 - x1) / (y2 - y1)
					xIntersect := float64(int64(t.y)-int64(c1.y))*float64(int64(c2.x)-int64(c1.x))/float64(int64(c2.y)-int64(c1.y)) + float64(c1.x)

					// if tile is to the left of intersection (or on vertical edge), ray crosses this edge
					if c1.x == c2.x || float64(t.x) <= xIntersect {
						isInside = !isInside // toggle inside/outside
					}
				}
			}
		}
	}

	return isInside
}
