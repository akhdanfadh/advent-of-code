package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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
	ranges, ingredients, err := readFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	// main logic
	process := partOne
	if *p2 {
		process = partTwo
	}
	result, err := process(ranges, ingredients)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fresh ingredients count: %d\n", result)
}

func readFile(fname string) ([][2]int, []int, error) {
	// open file
	file, err := os.Open(fname)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// read ranges
	ranges := make([][2]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		} // break on empty line

		var lo, hi int
		_, err = fmt.Sscanf(line, "%d-%d", &lo, &hi)
		if err != nil {
			return nil, nil, err
		}
		ranges = append(ranges, [2]int{lo, hi})
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	// read ingredients
	ingredients := make([]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		var ing int
		_, err = fmt.Sscanf(line, "%d", &ing)
		if err != nil {
			return nil, nil, err
		}
		ingredients = append(ingredients, ing)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return ranges, ingredients, nil
}

// brute force solution
func partOne(ranges [][2]int, ingredients []int) (int, error) {
	count := 0
	for _, ing := range ingredients {
		for _, r := range ranges {
			if ing >= r[0] && ing <= r[1] {
				// fmt.Printf("Ingredient %d is fresh (in range %d-%d)\n", ing, r[0], r[1])
				count++
				break
			}
		}
	}
	return count, nil
}

func partTwo(ranges [][2]int, ingredients []int) (int, error) {
	return 0, nil
}
