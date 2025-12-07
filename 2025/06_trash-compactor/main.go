package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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

	// main logic
	process := partOne // or partOneBrute
	if *p2 {
		process = partTwo
	}
	result, err := process(args[0])
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Printf("Total: %d\n", result)
}

func partOne(fname string) (uint64, error) {
	// open file
	file, err := os.Open(fname)
	if err != nil {
		return 0, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// read line by line
	syms := make([]byte, 0)
	nums := make([]uint64, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// handle last line (math symbols)
		if line[0] == '*' || line[0] == '+' {
			for symStr := range strings.FieldsSeq(line) {
				syms = append(syms, symStr[0])
			}
			break
		}

		// handle number line(s)
		for numStr := range strings.FieldsSeq(line) {
			num, err := strconv.ParseUint(numStr, 10, 64)
			if err != nil {
				return 0, err
			}
			nums = append(nums, num)
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	// process numbers and symbols
	result := uint64(0)
	numLines := len(nums) / len(syms)
	for symIdx, sym := range syms {
		var innerResult uint64
		if sym == '*' {
			innerResult = 1
			for i := range numLines {
				innerResult *= nums[i*len(syms)+symIdx]
			}
		} else {
			innerResult = 0
			for i := range numLines {
				innerResult += nums[i*len(syms)+symIdx]
			}
		}
		fmt.Printf("%d %c: %d\n", symIdx, sym, innerResult)
		result += innerResult
	}
	return result, nil
}

func partTwo(fname string) (uint64, error) {
	return 0, nil
}
