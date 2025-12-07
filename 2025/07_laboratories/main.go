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
	versionHandlers = map[string]func(string) (int, error){
		"1": partOne,
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
	fmt.Printf("The beam is splitted %d times\n", result)
}

func partOne(filename string) (int, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
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

	return splitCount, nil
}

func getSplitters(s string, b byte) *Set[int] {
	splitters := NewSet[int]()
	for i := 0; i < len(s); {
		idx := strings.IndexByte(s[i:], b)
		if idx == -1 {
			break
		} // not founc
		splitters.Add(i + idx)
		i += idx + 1 // next look after found index
	}
	return splitters
}
