package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
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
	result, err := process(file, *p2)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("Total output joltage: %d\n", result)
}

func process(file io.Reader, p2 bool) (int, error) {
	// first count lines for goroutine channel buffer size
	lc, err := lineCounter(file)
	if err != nil {
		return 0, fmt.Errorf("failed to count lines: %w", err)
	}
	fmt.Printf("Amount of battery banks: %d\n", lc)
	jolts := make(chan int, lc) // channel to collect results
	var wg sync.WaitGroup       // to synchronize goroutines

	// reset file pointer to beginning
	_, err = file.(io.Seeker).Seek(0, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("failed to seek to beginning of file: %w", err)
	}

	// read line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// now process the line concurrently
		wg.Add(1)
		go func(line string) {
			defer wg.Done()

			// try to visualize this yourself hehe
			left, right := 0, 1
			for i := 1; i < len(line); i++ {
				if line[i] > line[left] && i != len(line)-1 {
					left, right = i, i+1
				} else {
					if line[i] > line[right] {
						right = i
					}
				}
			}

			// TODO: error handling for Atoi inside goroutine
			jolt, _ := strconv.Atoi(string(line[left]) + string(line[right]))
			fmt.Printf("bank=%s, jolt=%d\n", line, jolt)
			jolts <- jolt
		}(line)
	}

	// close channel once all goroutines are done
	go func() {
		wg.Wait()
		close(jolts)
	}()

	// collect and sum
	totalJolt := 0
	for jolt := range jolts {
		totalJolt += jolt
	}
	return totalJolt, nil
}

// lineCounter is faster line counter using bytes.Count to find the newline characters
//
// It's faster because it takes away all the extra logic and buffering required to
// return whole lines, and takes advantage of some assembly optimized functions
// offered by the bytes package to search characters in a byte slice.
//
// Modified version of code from:
// Source - https://stackoverflow.com/a/24563853
// Posted by Mr_Pink, modified by community. See post 'Timeline' for change history
// Retrieved 2025-12-03, License - CC BY-SA 3.0
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lastCharWasNewline := true

	for {
		c, err := r.Read(buf)
		if c > 0 {
			count += bytes.Count(buf[:c], []byte{'\n'})
			lastCharWasNewline = buf[c-1] == '\n'
		}

		if err == io.EOF {
			if !lastCharWasNewline && c > 0 {
				count++
			}
			return count, nil
		}
		if err != nil {
			return count, err
		}
	}
}
