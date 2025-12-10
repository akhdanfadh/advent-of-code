package main

import (
	"bufio"
	"fmt"
	"os"
)

func processV1(filename string) (string, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close() // error ignored (file only for reading)

	// read file line by line
	tiles := []tile{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tile, err := parseTile(line)
		if err != nil {
			return "", err
		}
		tiles = append(tiles, tile)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// check largest area
	largestArea := uint(0)
	for i := 0; i < len(tiles); i++ {
		for j := i + 1; j < len(tiles); j++ {
			area := calcArea(tiles[i], tiles[j])
			if area > largestArea {
				largestArea = area
			}
		}
	}
	return fmt.Sprintf("Largest rectangle area: %d", largestArea), nil
}

func processV1a(filename string) (string, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close() // error ignored (file only for reading)

	// key idea: working row-by-row or column-by-column
	// for a fixed pair of rows/columns, we only need the leftmost and rightmost (topmost and bottommost) tiles
	// so it will be O(n+R^2) or O(n+C^2) instead of O(n^2),
	// where n is number of tiles (when reading), R is number of rows, C is number of columns

	type mm struct {
		min, max uint
	}
	rowMap := make(map[uint]mm)
	colMap := make(map[uint]mm)

	// read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tile, err := parseTile(line)
		if err != nil {
			return "", err
		}

		// check for row (y)
		if r, ok := rowMap[tile.y]; !ok {
			rowMap[tile.y] = mm{tile.x, tile.x} // first entry
		} else {
			if tile.x < r.min {
				r.min = tile.x
			}
			if tile.x > r.max {
				r.max = tile.x
			}
			rowMap[tile.y] = r
		}

		// check for column (x)
		if c, ok := colMap[tile.x]; !ok {
			colMap[tile.x] = mm{tile.y, tile.y} // first entry
		} else {
			if tile.y < c.min {
				c.min = tile.y
			}
			if tile.y > c.max {
				c.max = tile.y
			}
			colMap[tile.x] = c
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// check which is smaller (for efficiency)
	toCheck := rowMap
	if len(colMap) < len(rowMap) {
		toCheck = colMap
	}

	type slice struct {
		n        uint /// x or y
		min, max uint
	}

	// convert map to slice to interate with indices
	slices := make([]slice, 0, len(toCheck))
	for n, m := range toCheck {
		slices = append(slices, slice{n, m.min, m.max})
	}

	largestArea := uint(0)

	// case 1: check for same row/column
	for _, s := range slices {
		if s.max > s.min {
			area := (s.max - s.min + 1) * 1
			if area > largestArea {
				largestArea = area
			}
		}
	}

	// case 2: check for different rows/columns
	for i := 0; i < len(slices); i++ {
		si := slices[i]
		for j := i + 1; j < len(slices); j++ {
			sj := slices[j]

			// slice distance
			ds := absDiff(si.n, sj.n) + 1
			// max min distance
			d1 := absDiff(si.max, sj.min) + 1
			d2 := absDiff(sj.max, si.min) + 1
			dm := max(d1, d2)

			area := ds * dm
			if area > largestArea {
				largestArea = area
			}
		}
	}

	return fmt.Sprintf("Largest rectangle area: %d", largestArea), nil
}
