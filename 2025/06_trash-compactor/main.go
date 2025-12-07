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
			syms = parseSymLine(line)
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
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	// separate number lines and symbol line
	numLines := lines[:len(lines)-1]
	if !verifyNumLines(numLines) {
		return 0, fmt.Errorf("inconsistent line lengths in number lines")
	}
	symLine := lines[len(lines)-1]
	syms := parseSymLine(symLine)

	// process numbers by right-to-left one column at a time
	lenColumn := len(numLines[0]) // how many chars to process
	curNums := make([]uint64, 0)  // store current numbers for symbol operation
	curSymsIdx := len(syms) - 1   // which symbol to use for curNums
	grandResult := uint64(0)      // final result
	for i := lenColumn - 1; i >= -1; i-- {
		var curColumn []int

		// handle character column, i == -1 means we are done with digits (and process last symbol)
		if i >= 0 {
			for j := range numLines {
				char := numLines[j][i]
				num := byteToDigit(char)
				if num >= 0 {
					curColumn = append(curColumn, num)
				}
			}
		}

		// if curColumn is empty, then we math with symbols
		if len(curColumn) == 0 {
			sym := syms[curSymsIdx]
			var innerResult uint64
			if sym == '*' {
				innerResult = 1
				for _, num := range curNums {
					innerResult *= num
				}
			} else {
				innerResult = 0
				for _, num := range curNums {
					innerResult += num
				}
			}
			fmt.Printf("%d %c: %d\n", curSymsIdx, sym, innerResult)
			grandResult += innerResult
			curNums = curNums[:0] // reset for next column
			curSymsIdx--

			// otherwise, we collect the numbers
		} else {
			num := digitsToUint64(curColumn)
			curNums = append(curNums, num)
			fmt.Printf("col %d: %v -> %d\n", i, curColumn, num)
		}
	}

	return grandResult, nil
}

func parseSymLine(line string) (syms []byte) {
	for symStr := range strings.FieldsSeq(line) {
		syms = append(syms, symStr[0])
	}
	return syms
}

func verifyNumLines(lines []string) bool {
	if len(lines) == 0 {
		return true
	}
	charCount := len(lines[0])
	for i := 1; i < len(lines); i++ {
		if len(lines[i]) != charCount {
			return false
		}
	}
	return true
}

func byteToDigit(b byte) int {
	if b == ' ' {
		return -1
	}
	return int(b - '0')
}

func digitsToUint64(digits []int) (result uint64) {
	for _, d := range digits {
		result = result*10 + uint64(d)
	}
	return result
}
