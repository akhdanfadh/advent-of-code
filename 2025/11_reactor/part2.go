package main

import "fmt"

func processV2(filename string) (string, error) {
	connections, err := readFile(filename)
	if err != nil {
		return "", err
	}

	// key idea: just like v1, but we check the path

	// use state for memoization: 0 = neither, 1 = dac only, 2 = fft only, 3 = both
	// memo[device][state] = num of paths from device to "out" with that state
	memo := make(map[string]map[int]int)

	var dfs func(string, int) int
	dfs = func(from string, state int) int {
		if from == "out" {
			if state == 3 {
				return 1
			}
			return 0
		} // base case

		if memo[from] != nil {
			if count, exists := memo[from][state]; exists {
				return count // memoization check
			}
		}

		// update state based on current device
		newState := state
		if from == "dac" {
			// turn on bit 0; 00 | 01 = 01, 01 | 01 = 01, 10 | 01 = 11, 11 | 01 = 11
			newState |= 1
		}
		if from == "fft" {
			// turn on bit 1; 00 | 10 = 10, 01 | 10 = 11, 10 | 10 = 10, 11 | 10 = 11
			newState |= 2
		}

		var total int
		for _, to := range connections[from] {
			total += dfs(to, newState) // try all outgoing paths
		}

		if memo[from] == nil {
			memo[from] = make(map[int]int)
		}
		memo[from][state] = total
		return total
	}

	return fmt.Sprintf("Total possible path is %d", dfs("svr", 0)), nil
}
