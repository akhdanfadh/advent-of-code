package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	filename        *string
	version         *string
	versionHandlers = map[string]func(string) (string, error){
		"1": partOne,
		"2": partTwo,
	}
)

func validateFlags() error {
	if *filename == "" {
		return fmt.Errorf("input file name is required")
	}
	if _, ok := versionHandlers[*version]; !ok {
		return fmt.Errorf("invalid version: %s", *version)
	}
	return nil
}

func main() {
	// program input
	version = flag.String("v", "1", "logic version: only 1 available")
	filename = flag.String("f", "", "input file name (required)")
	flag.Parse()
	if err := validateFlags(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// main logic
	handler := versionHandlers[*version]
	result, err := handler(*filename)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Print(result)
}

func partOne(filename string) (string, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close() // error not checked as file only for reading

	// find beam origin 'S': ideally on the first line
	beams := NewSet[int]()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		beam := strings.IndexByte(line, 'S')
		if beam >= 0 {
			beams.Add(beam)
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if beams.Size() == 0 {
		return "", fmt.Errorf("no beam origin 'S' found in the input")
	}

	// now split the beam(s) while reading line by line
	splitCount := 0
	for scanner.Scan() {
		line := scanner.Text()

		// get splitters positions
		splitters := getSplitters(line, '^')
		if splitters.Size() == 0 {
			continue
		}

		// split if there are any beam positions at splitters
		for _, beam := range beams.Items() {
			if splitters.Contains(beam) {
				beams.Add(beam - 1)
				beams.Add(beam + 1)
				beams.Remove(beam)
				splitters.Remove(beam)
				splitCount++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return fmt.Sprintf("The beam is split %d times\n", splitCount), nil
}

func getSplitters(s string, b byte) *Set[int] {
	splitters := NewSet[int]()
	for i := 0; i < len(s); {
		// IndexByte for slightly faster search, i read somewhere hehe
		idx := strings.IndexByte(s[i:], b)
		if idx == -1 {
			break
		} // not found
		splitters.Add(i + idx)
		i += idx + 1 // next look after found index
	}
	return splitters
}

func partTwo(filename string) (string, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close() // error not checked as file only for reading

	// find beam origin 'S': ideally on the first line
	var beamOrigin int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		beamOrigin = strings.IndexByte(line, 'S')
		if beamOrigin >= 0 {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if beamOrigin < 0 {
		return "", fmt.Errorf("no beam origin 'S' found in the input")
	}

	// instead of processing line by line, we want some kind of backtracking here
	// so we read all lines first containing splitters
	splittersLines := make([]*Set[int], 0)
	for scanner.Scan() {
		line := scanner.Text()
		splitters := getSplitters(line, '^')
		if splitters.Size() == 0 {
			continue
		}
		splittersLines = append(splittersLines, splitters)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// memoization for backtracking
	var backtrack func(index, beam int) int
	type state struct{ index, beam int }
	memoization := make(map[state]int)
	backtrack = func(index, beam int) int {
		// index tells what splitters line we are processing
		// beam is the current "root" beam position

		// base: reached the end of the lines, count as one valid way
		if index >= len(splittersLines) {
			return 1
		}
		// get value from memo if already computed
		key := state{index, beam}
		if count, exists := memoization[key]; exists {
			return count
		}

		// if we have a split here, count on branch left and right
		var count int
		if splittersLines[index].Contains(beam) {
			count = backtrack(index+1, beam-1) + backtrack(index+1, beam+1)
		} else { // otherwise continue straight
			count = backtrack(index+1, beam)
		}
		memoization[key] = count
		return count
	}

	countTimeline := backtrack(0, beamOrigin)
	return fmt.Sprintf("The beam can be split in %d different ways\n", countTimeline), nil
}
