package main

import (
	"bufio"
	"fmt"
	"os"
)

func processV2(filename string) (string, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close() // error ignored (file only for reading)

	// the idea:
	// - manually build where are the red and green tiles
	//   note that the input red tiles are in order, i.e.,
	//   next line is always adjacent to previous line
	// - to build the green tiles, we build the boundary first from red tiles,
	//   then flood fill the grid surrounding it with BFS
	// - for every possible rectangle from red tiles, we check all tiles inside
	//   whehter it belongs to green tiles or red tiles

	// get all red tiles
	redTiles := []tile{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tile, err := parseTile(line)
		if err != nil {
			return "", err
		}
		redTiles = append(redTiles, tile)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	isRedTile := make(map[tile]bool)
	for _, t := range redTiles {
		isRedTile[t] = true
	}

	// get all green tiles
	isGreenTileBoundary := getGreenTilesBoundary(redTiles)
	isGreenTileInterior := getGreenTilesInterior(redTiles, isGreenTileBoundary)
	isGreenTile := make(map[tile]bool)
	for t := range isGreenTileBoundary {
		isGreenTile[t] = true
	}
	for t := range isGreenTileInterior {
		isGreenTile[t] = true
	}

	// check all pairs of red tiles as opposite corners
	largestArea := uint(0)
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			t1, t2 := redTiles[i], redTiles[j]

			// check if all tiles in the rectangle are red or green
			if isValidRectangle(t1, t2, isRedTile, isGreenTile) {
				area := calcArea(redTiles[i], redTiles[j])
				if area > largestArea {
					largestArea = area
				}
			}
		}
	}
	return fmt.Sprintf("Largest rectangle area: %d", largestArea), nil
}

func isValidRectangle(t1, t2 tile, isRedTile, isGreenTile map[tile]bool) bool {
	minX, maxX := min(t1.x, t2.x), max(t1.x, t2.x)
	minY, maxY := min(t1.y, t2.y), max(t1.y, t2.y)

	// check all tiles in the rectangle
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			t := tile{x, y}
			if !isRedTile[t] && !isGreenTile[t] {
				return false
			}
		}
	}
	return true
}

func getGreenTilesBoundary(redTiles []tile) map[tile]bool {
	isGreenTile := make(map[tile]bool)

	// connect consecutive red tiles
	// this is poeeible as problem mention the consecutive line is adjacent
	for i := 0; i < len(redTiles); i++ {
		t1 := redTiles[i]
		t2 := redTiles[(i+1)%len(redTiles)] // wrap around
		fmt.Println("Connecting", t1, "to", t2)

		// add all tiles between t1 and t2 (exclusive)
		if t1.x == t2.x { // same column, fill vertically
			minY, maxY := min(t1.y, t2.y), max(t1.y, t2.y)
			for y := minY; y <= maxY; y++ {
				isGreenTile[tile{t1.x, y}] = true
			}
		} else if t1.y == t2.y { // same row, fill horizontally
			minX, maxX := min(t1.x, t2.x), max(t1.x, t2.x)
			for x := minX; x <= maxX; x++ {
				isGreenTile[tile{x, t1.y}] = true
			}
		}
	}
	return isGreenTile
}

func getGreenTilesInterior(redTiles []tile, isGreenTileBoundary map[tile]bool) map[tile]bool {
	// Before flood fill (? is padding)
	// ????????????????
	// ?..............?
	// ?.......#XXX#..?
	// ?.......X...X..?
	// ?..#XXXX#...X..?
	// ?..X........X..?
	// ?..#XXXXXX#.X..?
	// ?.........X.X..?
	// ?.........#X#..?
	// ?..............?
	// ????????????????
	//
	// After flood fill (? is also E):
	// ????????????????
	// ?EEEEEEEEEEEEEE?
	// ?EEEEEEE#XXX#EE?
	// ?EEEEEEEXIIIXEE?
	// ?EE#XXXX#IIIXEE?
	// ?EEXIIIIIIIIXEE?
	// ?EE#XXXXXX#IXEE?
	// ?EEEEEEEEEXIXEE?
	// ?EEEEEEEEE#X#EE?
	// ?EEEEEEEEEEEEEE?
	// ????????????????

	// get bounding box from red tiles for the flood fill
	minX, maxX := redTiles[0].x, redTiles[0].x
	minY, maxY := redTiles[0].y, redTiles[0].y
	for _, t := range redTiles {
		minX = min(minX, t.x)
		maxX = max(maxX, t.x)
		minY = min(minY, t.y)
		maxY = max(maxY, t.y)
	}

	// add padding of 1 each side so we can flood fill from outside
	minX--
	maxX++
	minY--
	maxY++

	// bfs flood fill to mark all exterior tiles
	isExteriorTile := make(map[tile]bool)
	queue := []tile{{minX, minY}}
	isExteriorTile[queue[0]] = true // this is true
	i := 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:] // pop front
		fmt.Printf("%d Visiting exterior tile %v\n", i, current)
		i++

		neighbors := []tile{
			{current.x + 1, current.y},
			{current.x - 1, current.y},
			{current.x, current.y + 1},
			{current.x, current.y - 1},
		}
		for _, next := range neighbors {
			if next.x < minX || next.x > maxX || next.y < minY || next.y > maxY {
				continue // out of bounds
			}
			if isExteriorTile[next] {
				continue // already marked
			}
			if isGreenTileBoundary[next] {
				continue // can't cross the green wall
			}
			isExteriorTile[next] = true
			queue = append(queue, next)
		}
	}

	// then everything not outside and not boundary is interior green tile
	isGreenTileInterior := make(map[tile]bool)
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			t := tile{x, y}
			if !isExteriorTile[t] && !isGreenTileBoundary[t] {
				isGreenTileInterior[t] = true
			}
		}
	}
	return isGreenTileInterior
}
