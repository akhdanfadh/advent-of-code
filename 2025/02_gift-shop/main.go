package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// validate command line arguments
	s2 := flag.Bool("s2", false, "enable step 2 logic")
	flag.Parse()        // parse optional
	args := flag.Args() // get positional
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <input file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// open file
	file, err := os.Open(args[0])
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %s", err)
		}
	}()

	// main logic
	result, err := process(file, *s2)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("-------------------------\n")
	fmt.Printf("Total sum of invalid IDs: %d\n", result)
}

func process(file io.Reader, s2 bool) (int, error) {
	// read the first line (expected input format)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	totalSum := 0
	for scope := range strings.SplitSeq(line, ",") {
		// split the range to left and right
		var left, right int
		_, err := fmt.Sscanf(scope, "%d-%d", &left, &right)
		if err != nil {
			return 0, fmt.Errorf("failed to parse range %q: %w", scope, err)
		}

		// now process the range
		fmt.Printf("Processing range: %d - %d\n", left, right)
		for id := left; id <= right; id++ {
			if isMirrored(id) {
				fmt.Printf("  Invalid ID found: %d\n", id)
				totalSum += id
			}
		}
	}
	return totalSum, nil
}

func isMirrored(id int) bool {
	// count digits
	digits := 0
	for temp := id; temp > 0; temp /= 10 {
		digits++
	}
	if digits%2 == 1 {
		return false // skip odd number of digits
	}

	// extract halves with power of 10
	// eg, 1234 / 10*2 = 12, 1234 % 10*2 = 34
	divisor := 1
	for i := 0; i < digits/2; i++ {
		divisor *= 10
	}
	return id/divisor == id%divisor
}
