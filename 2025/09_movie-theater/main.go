package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	// program input
	filename := flag.String("f", "", "input file name (required)")
	version := flag.String("v", "1", "logic version")
	flag.Parse()

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
	var x, y uint
	tiles := []tile{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		_, err := fmt.Sscanf(line, "%d,%d", &x, &y)
		if err != nil {
			return "", err
		}
		tiles = append(tiles, tile{x, y})
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

func calcArea(t1, t2 tile) uint {
	var x, y uint
	if t1.x > t2.x {
		x = t1.x - t2.x + 1
	} else {
		x = t2.x - t1.x + 1
	}
	if t1.y > t2.y {
		y = t1.y - t2.y + 1
	} else {
		y = t2.y - t1.y + 1
	}
	return x * y
}
