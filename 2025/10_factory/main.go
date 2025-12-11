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
	case "1a":
		result, err = processV1a(*filename)
	case "2":
		result, err = processV2(*filename)
	case "2a":
		result, err = processV2a(*filename)
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

type machine struct {
	lightsReq  []bool
	joltageReq []int
	buttons    [][]int
}

func lightsEqual(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func readFile(filename string) ([]machine, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close() // error handling omitted for brevity

	// read line by line
	machines := []machine{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		machine, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		machines = append(machines, machine)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return machines, nil
}

func parseLine(line string) (machine, error) {
	// we assume line is under the format:
	// `[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}`
	// lightsConfig is under square brackets []
	// buttons are all under parentheses
	// joltage are under curly braces {}

	// get lights configuration
	lights := []bool{}
	leftBracket := strings.Index(line, "[")
	rightBracket := strings.Index(line, "]")
	if leftBracket == -1 || rightBracket == -1 || rightBracket < leftBracket {
		return machine{}, fmt.Errorf("invalid lights configuration")
	}
	for c := leftBracket + 1; c < rightBracket; c++ {
		switch line[c] {
		case '#':
			lights = append(lights, true)
		case '.':
			lights = append(lights, false)
		default:
			return machine{}, fmt.Errorf("invalid character in lights configuration: %c", line[c])
		}
	}

	// get joltage configuration
	joltage := []int{}
	leftBrace := strings.Index(line, "{")
	rightBrace := strings.Index(line, "}")
	if leftBrace == -1 || rightBrace == -1 || rightBrace < leftBrace {
		return machine{}, fmt.Errorf("invalid joltage configuration")
	}
	for numStr := range strings.SplitSeq(line[leftBrace+1:rightBrace], ",") {
		num, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil {
			return machine{}, fmt.Errorf("invalid joltage number: %s", numStr)
		}
		joltage = append(joltage, num)
	}

	// assert lights and joltage lengths match
	if len(lights) != len(joltage) {
		return machine{}, fmt.Errorf("mismatched lights and joltage lengths")
	}

	// get buttons configuration
	buttons := make([][]int, 0)
	for btnStr := range strings.FieldsSeq(line[rightBracket+1 : leftBrace]) {
		trimmed := strings.Trim(btnStr, "()")
		splitted := strings.Split(trimmed, ",")
		if len(splitted) == 1 {
			btnNum, err := strconv.Atoi(splitted[0])
			if err != nil {
				return machine{}, fmt.Errorf("invalid button number: %s", splitted[0])
			}
			buttons = append(buttons, []int{btnNum})
		} else if len(splitted) > 1 {
			innerBtns := []int{}
			for _, numStr := range splitted {
				btnNum, err := strconv.Atoi(numStr)
				if err != nil {
					return machine{}, fmt.Errorf("invalid button number: %s", numStr)
				}
				innerBtns = append(innerBtns, btnNum)
			}
			buttons = append(buttons, innerBtns)
		} else {
			return machine{}, fmt.Errorf("invalid button configuration: %s", btnStr)
		}
	}

	return machine{lights, joltage, buttons}, nil
}
