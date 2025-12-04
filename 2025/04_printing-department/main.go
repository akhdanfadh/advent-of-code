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

// brute force to the rescue haha
func partOne(grid [][]byte) (int, error) {
	result := 0
	directions := [8][2]int{
		{-1, 0},  // up
		{-1, 1},  // up-right
		{0, 1},   // right
		{1, 1},   // down-right
		{1, 0},   // down
		{1, -1},  // down-left
		{0, -1},  // left
		{-1, -1}, // up-left
	}

	for r := range len(grid) {
		for c := range len(grid[r]) {
			if !isRoll(grid[r][c]) {
				continue
			}

			adjacentRolls := 0
			for _, dir := range directions {
				nr, nc := r+dir[0], c+dir[1]
				if nr >= 0 && nr < len(grid) &&
					nc >= 0 && nc < len(grid[r]) &&
					isRoll(grid[nr][nc]) {
					adjacentRolls++
					if adjacentRolls >= 4 {
						break // stop early
					}
				}
			}

			if adjacentRolls < 4 {
				result++
			}
		}
	}
	return result, nil
}

func partTwo(grid [][]byte) (int, error) {
	return 0, nil
}
