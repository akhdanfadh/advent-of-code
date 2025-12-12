package main

import "fmt"

func processV1(filename string) (string, error) {
	deviceMap, err := readFile(filename)
	if err != nil {
		return "", err
	}

	// key idea: when talking about exploring all possibilities, we are talking about DFS
	// thing is the input graph actually have cycles AND
	// the problem does not explicitly specify whether we can revisit nodes or not
	// we here just make an assumption that we cannot revisit nodes, otherwise it would be infinite
	//
	// improvement: use memoization

	start := "you"
	target := "out"
	visited := make(map[string]bool)
	memo := make(map[string]int) // memo[c] = computed number of paths from c to target

	var dfs func(from string) int
	dfs = func(from string) int {
		if from == target {
			return 1
		} // base case, found 1 path
		if visited[from] {
			return memo[from] // already computed in current path, avoid cycle implicitly
		}

		visited[from] = true
		var total int
		for _, to := range deviceMap[from] { // try all outgoing paths
			total += dfs(to)
		}

		memo[from] = total
		return total
	}

	return fmt.Sprintf("Total possible path is %d", dfs(start)), nil
}
