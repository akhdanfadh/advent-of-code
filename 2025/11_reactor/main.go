package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
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
	case "2":
		result, err = processV2(*filename)
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

type connections map[string][]string

var (
	linePattern   = regexp.MustCompile(`^([a-z]{3}):\s*(.*)$`)
	devicePattern = regexp.MustCompile(`[a-z]{3}`)
)

func readFile(filename string) (connections, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close() // error handling omitted for brevity

	// read line by line
	connections := make(connections)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		from, to, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		connections[from] = to
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return connections, nil
}

func parseLine(line string) (string, []string, error) {
	matches := linePattern.FindStringSubmatch(line)
	if matches == nil || len(matches) != 3 {
		return "", nil, fmt.Errorf("invalid line format")
	}
	from := matches[1]
	to := devicePattern.FindAllString(matches[2], -1)
	return from, to, nil
}
