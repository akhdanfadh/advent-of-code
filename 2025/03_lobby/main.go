package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
	result, err := process(context.Background(), file, *p2)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Printf("Total output joltage: %d\n", result)
}

func process(ctx context.Context, file io.ReadSeeker, p2 bool) (int, error) {
	// create cancellable context from parent
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// first count lines for goroutine channel buffer size
	lc, err := lineCounter(file)
	if err != nil {
		return 0, fmt.Errorf("failed to count lines: %w", err)
	}
	fmt.Printf("Amount of battery banks: %d\n", lc)
	jolts := make(chan int, lc) // channel to collect results
	var wg sync.WaitGroup       // to synchronize goroutines
	// why not just unbuffered channel? "true parallelism"
	// with unbuffered, workers that finish early will just wait (eg line 1 goroutine finishes
	// before main goroutine start receiving, ie, the loop scanner.Scan hasn't done)
	// with buffered, early finishers can deposit results and exit, allowing more concurrency
	// BUT this comes with tradeoff where we read the file first, so tradeoff speed needs to be actually

	// reset file pointer to beginning
	// file must be also io.Seeker to make sure that file is seekable (if not, will panic)
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("failed to seek to beginning of file: %w", err)
	}

	// read line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		digits := scanner.Text()

		// now process the line concurrently
		wg.Add(1)
		go func(digits string) {
			defer wg.Done()

			joltFunc := joltOne
			if p2 {
				joltFunc = joltTwo
			}
			jolt := joltFunc(digits)
			fmt.Printf("bank=%s, jolt=%d\n", digits, jolt)

			select {
			case jolts <- jolt:
				// send succeeds, do nothing special (not fallthrough like C)
			case <-ctx.Done():
				return // abort if context cancelled
			}
		}(digits)
	}

	// if scanner somehow fails mid-way, prevent goroutine leaks
	if err := scanner.Err(); err != nil {
		cancel() // signal all goroutines to stop
		return 0, fmt.Errorf("failed to scan file: %w", err)
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

func joltOne(digits string) int {
	// try to visualize this yourself hehe
	// this logic makes our algorithm O(n*m) if not parallelized,
	// where n is length of digits, m is how many line of digits
	idxBestLeft, idxBestRight := 0, 1
	for idxPointer := 1; idxPointer < len(digits); idxPointer++ {
		// if current digit > best left digit and not the last digit,
		// make that current digit the best left digit, and the next digit the best right digit
		if digits[idxPointer] > digits[idxBestLeft] && idxPointer != len(digits)-1 {
			idxBestLeft, idxBestRight = idxPointer, idxPointer+1
		} else {
			// otherwise, check if current digit > best right digit,
			// if yes, make current digit the best right digit
			if digits[idxPointer] > digits[idxBestRight] {
				idxBestRight = idxPointer
			}
		}
	}

	// when we index a string, we get byte (value of 50 ASCII for '2')
	// so byte offset '0' (48 ASCII) to get actual digit value (in byte)
	digitLeft := digits[idxBestLeft] - '0'
	digitRight := digits[idxBestRight] - '0'
	jolt := int(digitLeft*10 + digitRight)
	return jolt
}

func joltTwo(digits string) int {
	// the idea is kind of sliding window, let's say
	// digits = 2357809, len(digits) = 7, need to find 4 digits to make largest
	// we need to 'somehow' iterate the digits 4 times to fill in the slot result
	// i.e. find the max in 2357, then find max in 3578, then 5780, then 7809 (with tricks!!)
	// see right there? haha greedy clever and this way it will lexicographically largest too
	const k = 12
	n := len(digits)
	result := 0

	idxStart := 0
	for i := range k {
		idxEnd := n - k + i // start to end inclusive

		// find the largest digits in this range
		idxBest := idxStart
		for j := idxStart; j <= idxEnd; j++ {
			if digits[j] > digits[idxBest] {
				idxBest = j
			}
		}

		// add to result
		digit := digits[idxBest] - '0'
		result = result*10 + int(digit)

		// TRICKY here, next search starts after the best found
		idxStart = idxBest + 1
	}

	return result
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
