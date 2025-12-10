package main

import (
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
	sampleSize := flag.Int("s", 0, "sample size for version 2")
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
		if *sampleSize < 1 {
			fmt.Fprintf(os.Stderr, "Error: sample size must be > 0\n")
			os.Exit(1)
		}
		result, err = processV2(*filename, *sampleSize)
	case "2a":
		result, err = processV2a(*filename)
	case "2b":
		result, err = processV2b(*filename)
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

type tile struct {
	x, y uint
}

func parseTile(s string) (tile, error) {
	sep := strings.IndexByte(s, ',')
	if sep == -1 {
		return tile{}, fmt.Errorf("invalid tile format: %s", s)
	}
	x, err := strconv.ParseUint(s[:sep], 10, 32)
	if err != nil {
		return tile{}, err
	}
	y, err := strconv.ParseUint(s[sep+1:], 10, 32)
	if err != nil {
		return tile{}, err
	}
	return tile{uint(x), uint(y)}, nil
}

func calcArea(t1, t2 tile) uint {
	x := absDiff(t1.x, t2.x) + 1
	y := absDiff(t1.y, t2.y) + 1
	return x * y
}

func absDiff(a, b uint) uint {
	if a > b {
		return a - b
	}
	return b - a
}
