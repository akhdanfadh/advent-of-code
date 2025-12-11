package main

// how to use: follow instructions at https://github.com/draffensperger/golp
import (
	"fmt"
	"slices"

	"github.com/draffensperger/golp"
)

func processV2a(filename string) (string, error) {
	machines, err := readFile(filename)
	if err != nil {
		return "", err
	}

	totalPresses := 0
	for i, m := range machines {
		presses, solution := m.solveWithGOLP()
		if presses < 0 {
			return "", fmt.Errorf("no solution found for machine %d", i)
		}
		fmt.Printf("Machine %d: total presses = %d, solution = %v\n", i, presses, solution)
		totalPresses += presses
	}

	return fmt.Sprintf("Total button presses for all machines: %d", totalPresses), nil
}

func (m *machine) solveWithGOLP() (int, []int) {
	numButtons, numCounters := len(m.buttons), len(m.joltageReq)

	// create LP problem: numCouters constraints, numButtons variables
	lp := golp.NewLP(numCounters, numButtons)

	// set objective: minimize sum of all buttons presses
	// each button press costs 1, so objective coefficients are all 1
	//
	// what this means mathemically: minimize 1.0*x1 + 1.0*x2 + ... + 1.0*xN
	// if objective := []float64{2.0, 1.0, 1.0, 3.0, 1.0, 1.0}
	// now button 0 "costs" 2 and button 3 "costs" 3, and solver would prefer buttons with lower costs
	objective := make([]float64, numButtons)
	for i := range objective {
		objective[i] = 1.0
	}
	lp.SetObjFn(objective)
	// golp does minimization by default, https://pkg.go.dev/github.com/draffensperger/golp#LP.SetMaximize

	// add constraints: for each counter, sum of buttons affecting it = target
	// what this means mathematically: we are building the equations like x1 + x3 + x5 = target1, x2 + x4 = target2, etc.
	for cntIdx := range numCounters {
		// build sparse constraint for this counter
		entries := []golp.Entry{}

		// check if this button affects this counter
		// if so, since each button press adds 1 to the counter, we add entry (btnIdx, 1.0)
		for btnIdx, btn := range m.buttons {
			if slices.Contains(btn, cntIdx) {
				entries = append(entries, golp.Entry{Col: btnIdx, Val: 1.0})
			}
		}

		// add constraint: sum = target[cntIdx]
		lp.AddConstraintSparse(entries, golp.EQ, float64(m.joltageReq[cntIdx]))
	}

	// set all variables to be non-negative integers
	for i := range numButtons {
		lp.SetInt(i, true)      // make variable integer
		lp.SetBounds(i, 0, 1e6) // non-negative with reasonable upper bound
	}

	// solve the problem
	solutionType := lp.Solve()
	if solutionType != golp.OPTIMAL && solutionType != golp.SUBOPTIMAL {
		// see https://lpsolve.sourceforge.net/5.5/solve.htm
		return -1, nil // solver failed
	}

	// after calling lp.Solve(), we get:
	// - lp.Variables() = the solution (x values), i.e., how many times to press each button
	// - lp.Objective() = the cost (what we get when we plug solution into the objective function), i.e., total button presses
	vars := lp.Variables()
	solution := make([]int, numButtons)
	totalPresses := 0
	for i := range numButtons {
		// even though we called SetInt, the solver returns float64 values, so we need to round them
		solution[i] = int(vars[i] + 0.5) // round to nearest int, e.g., 2.99 -> 3
		totalPresses += solution[i]
	}

	// verify solution
	if !m.verifyGOLPSolution(solution) {
		return -1, nil // invalid solution
	}

	return totalPresses, solution
}

func (m *machine) verifyGOLPSolution(solution []int) bool {
	result := make([]int, len(m.joltageReq))
	for btnIdx, count := range solution {
		for _, cntIdx := range m.buttons[btnIdx] {
			result[cntIdx] += count
		}
	}

	for i := range result {
		if result[i] != m.joltageReq[i] {
			return false
		}
	}
	return true
}
