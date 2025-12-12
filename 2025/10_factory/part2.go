package main

import (
	"fmt"
	"strconv"
	"strings"
)

func processV2(filename string) (string, error) {
	// approach brainstorming 1:
	// - we can see that, similar to our v1a, we can approach this problem using linear equations, and solve the augmented matrix Ax=b just like before
	// - but since we now are using natural numbers and not just booleans that we can do XOR with, we have more complexity
	// - thing is solving Ax=b may gives us non-integer or even negative solutions, which are invalid as we want natural number counts.
	// - even if it gives one natural number solution, there may be other solutions with fewer button presses
	// - we need a sophisticated method or "solver" that respects our constraints,
	// - this what those "Integer Linear Programming" (ILP) solvers do such as https://github.com/draffensperger/golp
	// let's try this simple library import in v2a.
	//
	// approach brainstorming 2:
	// - we make an unpleasant assumption here: for AoC-style puzzles, the test inputs are usually small enough that "maybe" shortest-part search is feasible
	// - talking about short-path search, we are reminded of BFS or Dijkstra's algorithm
	// - so instead of thinking of "solving equations", we think of this as a graph to "build up the counters":
	//   - we start from all counters at 0: (0, ..., 0)
	//   - each time we press a button, we add its effect to the current counter state
	//     if button 0 affects counters (0, 2), that means we add (1, 0, 1, 0, ..., 0) to the current state
	//   - we are allowed to press buttons in any order, any number of times, BUT
	//     we never want to exceed the target in any counter (remember counters only increase, so its impossible to "undo" a press)
	// - so shortest path in this state graph from (0,...,0) to the target state is our solution
	// NOOOO, TURNS OUT THIS IS IMPRACTICAL FOR OUR INPUT haha
	// - let b_j = target value of counter j, assume rough upper bound B = max(b_j)
	// - then total number of states is roughly B^n where n = number of counters
	// - just see machine 3 in input, {10,187,228,38,28,192,33,218} -> 228^8 ~= 7.3e18 states, way too large

	machines, err := readFile(filename)
	if err != nil {
		return "", err
	}

	totalPresses := 0
	for i, m := range machines {
		presses := m.solveBFS()
		if presses < 0 {
			return "", fmt.Errorf("no solution found for machine %d", i)
		}
		fmt.Printf("Machine %d: total presses = %d\n", i, presses)
		totalPresses += presses
	}

	return fmt.Sprintf("Total button presses for all machines: %d", totalPresses), nil
}

func (m *machine) solveBFS() int {
	numButtons, numCounters := len(m.buttons), len(m.joltageReq)

	// convert each button into a fixed-length effect array of size numCounters
	// effect[i][j] = 1 if button i affects counter j, else 0
	effects := make([][]int, numButtons)
	for btnIdx, btn := range m.buttons {
		effect := make([]int, numCounters)
		for _, cntIdx := range btn {
			effect[cntIdx] = 1
		}
		effects[btnIdx] = effect
	}

	// BFS in "counter space"
	type node struct {
		state []int // current counter values
		dist  int   // total button presses to reach this state
	}

	startState := make([]int, numCounters) // all zeros
	goalState := m.joltageReq
	visited := map[string]bool{encodeState(startState): true}

	queue := []node{{state: startState, dist: 0}}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:] // pop front

		// if we've reached the goal, return distance directly, this is shortest by BFS logic
		if isArrayEqual(curr.state, goalState) {
			return curr.dist
		}

		// otherwise, expand neighbors by pressing each button once
		for _, effect := range effects {
			nextState := make([]int, numCounters)
			valid := true

			// get next state by adding effect
			for i := range numCounters {
				nextStateVal := curr.state[i] + effect[i]
				if nextStateVal > goalState[i] {
					// reject any move that overshoots the goal in any counter,
					// as conters only increase and we can't "undo" presses, guarantee no solution from here
					valid = false
					break
				}
				nextState[i] = nextStateVal
			}
			if !valid { // this button press is not allowed (overshoot), skip
				continue
			}

			nextKey := encodeState(nextState)
			if visited[nextKey] { // we've already seen this state via a shorter or equal path, skip
				continue
			}
			visited[nextKey] = true

			queue = append(queue, node{state: nextState, dist: curr.dist + 1})
		}
	}

	return -1 // no solution found
}

func encodeState(state []int) string {
	var sb strings.Builder
	for i, x := range state {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(x))
	}
	return sb.String()
}

func isArrayEqual(a, b []int) bool {
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
