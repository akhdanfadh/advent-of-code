package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
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

type (
	shape  [][]int
	region struct {
		width, height int
		presentsCount []int
	}
	aoc struct {
		presents map[int]*shape
		regions  []*region
	}
)

var (
	regionPattern = regexp.MustCompile(`^(\d+)x(\d+):\s*(.*)$`)
	posNumPattern = regexp.MustCompile(`\d+`)
)

func readInput(filename string) (*aoc, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := &aoc{
		presents: make(map[int]*shape),
		regions:  []*region{},
	}

	// we read the file in two sections: shapes and regions
	// since we read line by line, we need to track the current shape being read
	var curShapeBuf [][]int
	curShapeID := -1
	inShapes := true
	flushShape := func() {
		if curShapeID >= 0 {
			s := shape(curShapeBuf)
			result.presents[curShapeID] = &s
		}
		curShapeID = -1
		curShapeBuf = nil // we initialize it when we read a new shape
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // skip empty line
		}

		// if region section detected, switch mode
		if inShapes && strings.Contains(line, "x") && strings.Contains(line, ":") {
			inShapes = false
			flushShape()
		}

		if inShapes {
			// shape header: "id:"
			if strings.HasSuffix(line, ":") {
				flushShape()

				idStr := strings.TrimSuffix(line, ":")
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return nil, err
				}
				curShapeID = id
				curShapeBuf = [][]int{} // reset buffer
				continue
			}

			// shape row
			row := make([]int, len(line))
			for i, ch := range line {
				if ch == '#' {
					row[i] = 1
				} else {
					row[i] = 0
				}
			}
			curShapeBuf = append(curShapeBuf, row)

		} else {
			// at this point, we are in region section
			matches := regionPattern.FindStringSubmatch(line)
			if matches == nil || len(matches) != 4 {
				return nil, fmt.Errorf("invalid region line: %q", line)
			}

			width, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, err
			}
			height, err := strconv.Atoi(matches[2])
			if err != nil {
				return nil, err
			}

			countStrs := posNumPattern.FindAllString(matches[3], -1)
			presentsCount := make([]int, len(countStrs))
			for i, cs := range countStrs {
				v, err := strconv.Atoi(cs)
				if err != nil {
					return nil, err
				}
				presentsCount[i] = v
			}

			result.regions = append(result.regions, &region{width, height, presentsCount})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
