package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	// validate command line arguments
	p2 := flag.Bool("p2", false, "enable step 2 logic")
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
	fmt.Printf("Total sum of invalid IDs: %d\n", result)
}

func process(ctx context.Context, file io.Reader, p2 bool) (int, error) {
	// create cancellable context from parent
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// read the first line (expected input format)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	scopes := strings.Split(line, ",")
	results := make(chan int, len(scopes)) // channel to collect results
	var wg sync.WaitGroup                  // to synchronize goroutines

	for _, scope := range scopes {
		// split the range to left and right
		var left, right int
		_, err := fmt.Sscanf(scope, "%d-%d", &left, &right)
		if err != nil {
			cancel() // signal workers to stop
			return 0, fmt.Errorf("failed to parse range %q: %w", scope, err)
		}

		// now process the range concurrently
		wg.Add(1)
		go func(left, right int) {
			defer wg.Done()
			sum := 0
			check := isMirrored
			if p2 {
				check = isRepeated
			}
			for id := left; id <= right; id++ {
				select {
				case <-ctx.Done():
					return // exit early if context is cancelled
				default:
					if check(id) {
						sum += id
					}
				}
			}
			results <- sum
		}(left, right)
	}

	// close channel once all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// collect and sum
	totalSum := 0
	for sum := range results {
		totalSum += sum
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

var divisors = make(map[int][]int)

func init() {
	// from pre-build divisors for lengths 1 to 20 (hard coded)
	// TODO: optimization, use thread-safe divisors map
	for n := 1; n <= 20; n++ {
		divisors[n] = buildDivisors(n)
	}
}

func isRepeated(id int) bool {
	// logic is to get the divisors of num digits of id
	// eg, for 12121212 (num digits 8), divisors are [1,2,4]
	// then for each divisor, check if the pattern repeats
	// eg, div 1 is 1->2 so go to next divisor, div 2 12->12->12->12 return true
	digits := strconv.Itoa(id)
	for _, patternLen := range divisors[len(digits)] {
		matches := true
		pattern := digits[:patternLen] // to match against
		for i := patternLen; i < len(digits); i += patternLen {
			if digits[i:i+patternLen] != pattern {
				matches = false
				break
			}
		}
		if matches {
			return true
		} // found a repeating pattern, direct return
	}
	return false
}

func buildDivisors(n int) []int {
	// eg, n=24 => [1,2,3,4,6,8,12]
	var res []int
	for i := 1; i <= n/2; i++ {
		if n%i == 0 {
			res = append(res, i)
		}
	}
	return res
}
