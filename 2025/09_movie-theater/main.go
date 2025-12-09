package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// program input
	filename := flag.String("f", "", "input file name (required)")
	version := flag.String("v", "1", "logic version")
	flag.Parse()

	now := time.Now()

	// main logic
	var result string
	var err error
	switch *version {
	case "1":
		result, err = processV1(*filename)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown version %s\n", *version)
		os.Exit(1)
	}

	diff := time.Since(now)
	fmt.Printf("Time taken: %v\n", diff)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}

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

type tile struct {
	x, y uint
}

func parseTile(s string) (tile, error) {
	sep := strings.IndexByte(s, ',')
	if sep == -1 {
		return tile{}, fmt.Errorf("invalid tile format: %s", s)
	}
	x, err := strconv.ParseUint(s[:sep], 10, 32)
	if err != nil {
		return tile{}, err
	}
	y, err := strconv.ParseUint(s[sep+1:], 10, 32)
	if err != nil {
		return tile{}, err
	}
	return tile{uint(x), uint(y)}, nil
}

func calcArea(t1, t2 tile) uint {
	x := absDiff(t1.x, t2.x) + 1
	y := absDiff(t1.y, t2.y) + 1
	return x * y
}

func absDiff(a, b uint) uint {
	if a > b {
		return a - b
	}
	return b - a
}
