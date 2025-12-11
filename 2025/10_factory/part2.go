package main

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

	return "not implemented", nil
}
