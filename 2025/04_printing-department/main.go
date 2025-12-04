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
	for r := range len(grid) {
		for c := range len(grid[r]) {
			if !isRoll(grid[r][c]) {
				continue
			}

			adjacentRolls := 0
			if r-1 >= 0 && isRoll(grid[r-1][c]) {
				adjacentRolls++
			} // up
			if r-1 >= 0 && c+1 < len(grid[r]) && isRoll(grid[r-1][c+1]) {
				adjacentRolls++
			} // up-right
			if c+1 < len(grid[r]) && isRoll(grid[r][c+1]) {
				adjacentRolls++
			} // right
			if r+1 < len(grid) && c+1 < len(grid[r]) && isRoll(grid[r+1][c+1]) {
				adjacentRolls++
			} // down-right
			if r+1 < len(grid) && isRoll(grid[r+1][c]) {
				adjacentRolls++
			} // down
			if r+1 < len(grid) && c-1 >= 0 && isRoll(grid[r+1][c-1]) {
				adjacentRolls++
			} // down-left
			if c-1 >= 0 && isRoll(grid[r][c-1]) {
				adjacentRolls++
			} // left
			if r-1 >= 0 && c-1 >= 0 && isRoll(grid[r-1][c-1]) {
				adjacentRolls++
			} // up-left

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
