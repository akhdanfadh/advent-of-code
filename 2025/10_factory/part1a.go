package main

import "fmt"

// processV1a uses Gaussian Elimination to solve the button-light toggle problem.
//
// This is much faster than brute force in processV1 for larger number of buttons.
// Complexity is O(numButtons^2 x numLights) compared to O(2^numButtons x numLights).
func processV1a(filename string) (string, error) {
	machines, err := readFile(filename)
	if err != nil {
		return "", err
	}

	totalMinPresses := 0
	for _, m := range machines {
		// step 1: build augmented matrix
		matrix := m.buildAugmentedMatrix()

		// step 2: perform Gaussian elimination
		freeVars := m.gaussianElimination(matrix)

		// step 3: check for inconsistencies (no solution exists)
		if !m.isConsistent(matrix) {
			return "", fmt.Errorf("no solution exists for this machine configuration")
		}

		// step 4: find the solution with back substitution
		minPresses := m.findMinButtonPresses(matrix, freeVars)
		totalMinPresses += minPresses
	}

	return fmt.Sprintf("Total minimum button presses: %d", totalMinPresses), nil
}

// buildMatrix creates the augmented matrix [A | b].
// matrix[i][j] = 1 if button j toggles light i, else 0.
// The last column is the target state for each light.
func (m *machine) buildAugmentedMatrix() [][]int {
	numButtons, numLights := len(m.buttons), len(m.lightsReq)
	matrix := make([][]int, numLights)
	for i := range matrix {
		matrix[i] = make([]int, numButtons+1) // +1 for target state
	}

	// fill the matrix: mark which buttons affect which lights
	for btnIdx, lights := range m.buttons {
		for _, lightIdx := range lights {
			matrix[lightIdx][btnIdx] = 1
		}
	}

	// add target column (what each light should be)
	for i, state := range m.lightsReq {
		if state {
			matrix[i][numButtons] = 1
		}
	}

	return matrix
}

// gaussianElimination transforms matrix to row echelon form (a triangular shape) in GF(2) (binary field).
// Returns a list of "free variables" (buttons that can be 0 or 1).
func (m *machine) gaussianElimination(matrix [][]int) []int {
	numButtons, numLights := len(m.buttons), len(m.lightsReq)
	pivotRow := 0 // track which row we're working on; work our way down
	freeVars := []int{}

	// for each column, we need to find a "pivot": a row that has a 1 in this column
	// - we only look at rows from pivotRow downwards because rows above are already processed
	// - when we find a 1, we swap that row with current pivotRow to put it in the right position
	// - why swap? imagine we're solving equations on paper: we'd naturally arrange them so
	//   the first equation has the first variable, the second equation has the second variable, etc
	// this is like moving an equation to the top of our list so we can use it to eliminate variables in other equations
	//
	// we iterate through each column (button), and stop either when we've processed all buttons or all rows
	// double condition is important bcs we may have underdetermined system (more buttons than lights) or otherwise
	for col := 0; col < numButtons && pivotRow < numLights; col++ {
		// find a row with a 1 in this column (at or below current pivotRow)
		foundPivot := false
		for row := pivotRow; row < numLights; row++ {
			if matrix[row][col] == 1 {
				// swap this row with pivotRow
				matrix[row], matrix[pivotRow] = matrix[pivotRow], matrix[row]
				foundPivot = true
				break
			}
		}

		// if no pivot found (we can't find any row with a 1 in this column),
		// it means this button (variable) doesn't have a unique value determined by the system.
		// it is "free": we can set it to either 0 or 1 without affecting consistency.
		// this happens when the system has multiple solutions (underdetermined).
		// we record this free variable and move to the next column.
		//
		if !foundPivot {
			freeVars = append(freeVars, col)
			continue
		}

		// now comes the elimination step:
		// - for every OTHER row (not the pivotRow) that has a 1 in this column, we XOR it with the pivotRow
		// - this "eliminates" the variable from that row (sets it to 0). for example,
		//   if equation 1 is `B0 + B2 = 1` and equation 2 is `B0 + B1 = 1`,
		//   we can eliminate B0 from equation 2 by XORing them: `(B0 + B1) XOR (B0 + B2) = B1 + B2`.
		//   The B0 cancels out because `B0 XOR B0 = 0`.
		// before elimination:       after eliminating col 0:
		// [1 0 1 | 1]               [1 0 1 | 1]
		// [1 1 0 | 1]         --->  [0 1 1 | 0]  <- XORed with row 0
		// [0 1 1 | 0]               [0 1 1 | 0]
		//
		for row := range numLights {
			if row != pivotRow && matrix[row][col] == 1 {
				// XOR this row with pivotRow, this sets matrix[row][col] to 0
				for c := 0; c <= numButtons; c++ {
					matrix[row][c] ^= matrix[pivotRow][c] // ^= means XOR assignment
				}
			}
		}

		// after processing this column, move to next row
		pivotRow++
	}
	return freeVars
}

// isConsistent checks if the system has a solution.
// Returns false if there's a row like [0 0 0 | 1] (contradiction).
func (m *machine) isConsistent(matrix [][]int) bool {
	numButtons := len(m.buttons)
	for _, row := range matrix {
		allZero := true
		for col := range numButtons {
			if row[col] != 0 {
				allZero = false
				break
			}
		}

		// this row is saying "no buttons affect this light, but we need it to be ON": impossible
		// if the target were 0 instead, it would be fine because lights start OFF by default
		if allZero && row[numButtons] == 1 {
			return false
		}
	}
	return true
}

// findMinButtonPresses finds which combination of free variables gives us the minimum button presses.
func (m *machine) findMinButtonPresses(matrix [][]int, freeVars []int) int {
	numButtons := len(m.buttons)

	// if there are no free variables, there's only one solution
	// so we just compute it and count how many buttons we need to press
	if len(freeVars) == 0 {
		solution := m.backSubstitute(matrix, make([]int, numButtons))
		return countOnes(solution)
	}

	// if there are free variables, try all combinations
	// this is still efficient because typically the number of free variables is small
	minPresses := numButtons + 1          // start with a large number
	numCombinations := 1 << len(freeVars) // 2^(number of free variables)
	for combo := range numCombinations {
		// set free variables according to this combination
		// this work like masking bits in combo to decide which free vars are 0 or 1 (similar to processV1)
		buttonStates := make([]int, numButtons)
		for i, varIdx := range freeVars {
			if combo&(1<<i) != 0 {
				buttonStates[varIdx] = 1
			}
		}

		// solve for the other variables using back substitution and count how many buttons are pressed
		solution := m.backSubstitute(matrix, buttonStates)
		presses := countOnes(solution)
		if presses < minPresses {
			minPresses = presses
		}
	}

	return minPresses
}

// backSubstitute solves for the non-free variables given the free variable values.
func (m *machine) backSubstitute(matrix [][]int, freeVarValues []int) []int {
	numButtons := len(m.buttons)
	solution := make([]int, numButtons)
	copy(solution, freeVarValues)

	// work backwards through the matrix (back substitution)
	// - why backwards? because after Gaussian elimination, the matrix is in row echelon form,
	//   where the bottom rows have fewer variables.
	// - the last row might have just one variable, the second last row might have two, etc.
	// - by working backwards, we can solve for variables one at a time without needing the
	//   values of variables we haven't solved yet.
	// for example:
	// Row 0: B0 XOR B2 = 1        <- needs B2 to solve for B0
	// Row 1: B1 XOR B2 = 0        <- needs B2 to solve for B1
	// Row 2: nothing (all zeros)
	// when we go bottom-up: B2 is free (set it first), then solve B1, then B0
	//
	for row := len(matrix) - 1; row >= 0; row-- {
		// find the pivot column (leftmost 1) in this row
		// this is the variable we're going to solve for in this row
		pivotCol := -1
		for col := range numButtons {
			if matrix[row][col] == 1 {
				pivotCol = col
				break
			}
		}

		// skip rows with no pivot (all zeros)
		// this is one of those "0 = 0" rows that don't give us any info
		if pivotCol == -1 {
			continue
		}

		// now we solve for the pivot variable: `pivot XOR other_vars = target`.
		// - to solve for pivot, we XOR both sides with the other_vars: `pivot = target XOR other_vars`.
		// - we start with `val = target`, then XOR it with all the other variable values that appear
		//   in this row (where the matrix has a 1).
		// - the columns after `pivotCol` are the "other_vars" in this row: we've already solved them
		//   because we're working backwards.
		// for example:
		// row: [0 1 1 | 0]  means  B1 XOR B2 = 0
		// start with val = 0 (target)
		// if B2 = 1 (already known), then val = val XOR solution[2] = 0 XOR 1 = 1, so B1 = 1
		//
		val := matrix[row][numButtons] // start with target value
		for col := pivotCol + 1; col < numButtons; col++ {
			if matrix[row][col] == 1 {
				val ^= solution[col] // XOR with known button value
			}
		}
		solution[pivotCol] = val
	}

	return solution
}

// countOnes counts how many 1s are in the array
func countOnes(arr []int) int {
	count := 0
	for _, v := range arr {
		count += v
	}
	return count
}
