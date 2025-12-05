package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
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
	process := partOne // or partOneBrute
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

func sortRanges(ranges [][2]int) {
	// this built-in sort function will sort by:
	// if return negative, a should come before b
	// if return positive, a should come after b
	// if return zero, a and b are equal in sort order
	slices.SortFunc(ranges, func(a, b [2]int) int {
		return a[0] - b[0]
	})
}

func mergeRanges(ranges [][2]int) [][2]int {
	// this assumes ranges are already sorted by start
	// cur: [1-5], next: [3-8] -> overlap
	// cur: [1-5], next: [6-8] -> adjacent
	// cur: [1-5], next: [7-8] -> separate
	merged := make([][2]int, len(ranges))
	cur := ranges[0]
	for i := 1; i < len(ranges); i++ {
		next := ranges[i]
		if next[0] <= cur[1]+1 { // overlap or adjacent
			if next[1] > cur[1] {
				cur[1] = next[1]
			}
		} else { // separate
			merged = append(merged, cur)
			cur = next
		}
	}
	merged = append(merged, cur) // last range
	return merged
}

func isFreshBinarySearch(id int, ranges [][2]int) bool {
	left, right := 0, len(ranges)-1
	for left <= right {
		mid := left + (right-left)/2
		r := ranges[mid]
		if id < r[0] { // before this range, search left
			right = mid - 1
		} else if id > r[1] { // after this range, search right
			left = mid + 1
		} else { // in this range
			return true
		}
	}
	return false
}

// binary search solution
func partOne(ranges [][2]int, ingredients []int) (int, error) {
	// preprocess ranges: sort and merged
	sortRanges(ranges)
	ranges = mergeRanges(ranges)

	// process ingredients: binary search
	count := 0
	for _, ing := range ingredients {
		if isFreshBinarySearch(ing, ranges) {
			count++
		}
	}
	return count, nil
}

// brute force solution
func partOneBrute(ranges [][2]int, ingredients []int) (int, error) {
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
