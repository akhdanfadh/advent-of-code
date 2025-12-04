package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// validate command line arguments
	p2 := flag.Bool("p2", false, "enable part two logic")
	flag.Parse()        // parse optional
	args := flag.Args() // get positional
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <input file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// read file into memory (variable)
	grid, err := readFileAs2DGrid(args[0])
	if err != nil {
		log.Fatal(err)
	}

	// main logic
	process := partOne
	if *p2 {
		process = partTwo
	}
	result, err := process(grid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Accessible paper rolls amount: %d\n", result)
}

func readFileAs2DGrid(fname string) ([][]byte, error) {
	// open file
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()

	// read line by line
	var grid [][]byte // byte for efficiency (also we already know it is ASCII)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []byte(line)
		grid = append(grid, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return grid, nil
}

func isRoll(char byte) bool {
	return char == '@'
}

var directions = [8][2]int{
	{-1, 0},  // up
	{-1, 1},  // up-right
	{0, 1},   // right
	{1, 1},   // down-right
	{1, 0},   // down
	{1, -1},  // down-left
	{0, -1},  // left
	{-1, -1}, // up-left
}

func countAdjacentRolls(grid [][]byte, r, c int) int {
	count := 0
	for _, dir := range directions {
		nr, nc := r+dir[0], c+dir[1]
		if nr >= 0 && nr < len(grid) &&
			nc >= 0 && nc < len(grid[r]) &&
			isRoll(grid[nr][nc]) {
			count++
			if count >= 4 {
				break // stop early
			}
		}
	}
	return count
}

// brute force to the rescue haha
func partOne(grid [][]byte) (int, error) {
	resultChan := make(chan int, len(grid))

	for r := range len(grid) {
		go func(r int) {
			result := 0
			for c := range len(grid[r]) {
				if !isRoll(grid[r][c]) {
					continue
				}
				if countAdjacentRolls(grid, r, c) < 4 {
					result++
				}
			}
			resultChan <- result
		}(r)
	}

	// collect results
	total := 0
	for range grid {
		total += <-resultChan
	}
	return total, nil
}

func partTwo(grid [][]byte) (int, error) {
	// read and store where rolls are (initially)
	rolls := make([][2]int, 0)
	for r := range len(grid) {
		for c := range len(grid[r]) {
			if isRoll(grid[r][c]) {
				rolls = append(rolls, [2]int{r, c})
			}
		}
	}

	// while loop until there is no roll to remove
	result := 0
	removed := make([][2]int, len(rolls)) // preallocate
	stayed := make([][2]int, len(rolls))
	for {
		removed, stayed = removed[:0], stayed[:0] // reset slices

		// check all rolls in current iteration
		for _, pos := range rolls {
			if countAdjacentRolls(grid, pos[0], pos[1]) < 4 {
				removed = append(removed, pos)
			} else {
				stayed = append(stayed, pos)
			}
		}

		// break loop as there are no roll to remove
		if len(removed) == 0 {
			break
		}

		// update rolls and grid for next iteration
		fmt.Printf("Removed %d from %d rolls\n", len(removed), len(rolls))
		for _, pos := range removed {
			grid[pos[0]][pos[1]] = '.' // mark as removed
		}
		result += len(removed)
		rolls, stayed = stayed, rolls
		// we do swap here because stayed will be reset in next slices
		// if we do rolls = stayed, then rolls itself will be lost
	}
	return result, nil
}
