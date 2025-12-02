package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	s2 := flag.Bool("s2", false, "enable step 2 password logic")
	flag.Parse() // important!

	args := flag.Args() // get positional arguments
	if len(args) != 1 {
		// error msg and usage info should go to stderr
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <input file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	pass, err := getPassword(args[0], *s2)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("The password is: %d\n", pass)
}

func getPassword(fname string, s2 bool) (int, error) {
	// open file
	file, err := os.Open(fname)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %s", err)
		}
	}()

	dialPos := 50
	zeroCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // read line by line
		line := scanner.Text()
		if len(line) < 2 {
			return 0, fmt.Errorf("invalid line: %s", line)
		} // prevent panic on slicing

		numSteps, err := strconv.Atoi(line[1:]) // string to int
		if err != nil {
			return 0, fmt.Errorf("failed to convert steps to integer: %w", err)
		}

		switch line[0] {
		case 'R':
			dialPos = (dialPos + numSteps) % 100 // simple modulo logic to wrap around
		case 'L':
			dialPos = ((dialPos-numSteps)%100 + 100) % 100 // tricky: go mod of negative number
		default:
			return 0, fmt.Errorf("invalid direction: %q", line[0])
		}

		if dialPos == 0 {
			zeroCount++
		} // password is number of times dial stop at 0
	}

	return zeroCount, nil
}
