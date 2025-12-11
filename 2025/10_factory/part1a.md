# processV1a: XOR + Gaussian Elimination

## Rewriting the problem as XOR linear equations

We can actually represent the button and light problem using linear algebra.
Let's write out what happens to Light 0:
`Light 0 final state = (Light 0 initial state + (does Button 0 affect it?) x (is Button 0 pressed?) + (does Button 1 affect it?) x (is Button 1 pressed?) + ...) mod 2`.
Since all lights start OFF, the initial state is 0, so we can ignore that term.
Then if, for example, button 0 and 2 affect Light 0, while the other buttons do not, we have:
`L0 = (1xB0 + 0xB1 + 1xB2 + ...) mod 2 = (B0 + B2) mod 2` where `Bx` is 1 if button x is pressed, and 0 otherwise.

Now here's the trick: in mod 2 arithmetic, addition is the just the same as XOR operation.
Observe this behavior:

```
(0 + 0) mod 2 = 0 = 0 XOR 0
(0 + 1) mod 2 = 1 = 0 XOR 1
(1 + 0) mod 2 = 1 = 1 XOR 0
(1 + 1) mod 2 = 0 = 1 XOR 1
```

So when we write `(B0 + B2) mod 2`, we can just as well write `B0 XOR B2`.

### XOR Basics

XOR has several important properties that make it useful for our problem.
It is commutative, meaning `A XOR B = B XOR A`, and associative, meaning `(A XOR B) XOR C = A XOR (B XOR C)`.
The identity element is 0, so `A XOR 0 = A`.
Most importantly, XOR is its own inverse: `A XOR A = 0`.

This self-inverse property means that `A XOR B XOR B = A`, which is like subtraction in regular math where `A + B - B = A`.
In XOR math, there's no difference between "adding" and "subtracting" because XORing with the same value twice cancels out (just like toggling a light switch twice returns it to its original state).

## Solving linear equations with Gaussian elimination

When we have a system of linear equations, we can use a method called Gaussian elimination to solve them.
For a reminder, it's a system where you can solve `2x + 3y = 5` and `4x + 5y = 11` by manipulating the equations to isolate one variable at a time.
Also, we can represent the system as a matrix and perform row operations to simplify it (well, I learn this in university courses)

Okay, let's say we have 3 lights and 3 buttons, with the following relationships:

- Light 0 is affected by Button 0 and Button 2, and we want it ON
- Light 1 is affected by Button 0 and Button 1, and we want it ON
- Light 2 is affected by Button 1 and Button 2, and we want it OFF

We write the system as an augmented matrix [A | b], where each row represents a light and each column (except the last) represents a button.
The value at position `(i, j)` is 1 if button `j` affects light `i`, and 0 otherwise.
The last column contains our target state for each light.

```
          B0  B1  B2 | Target
Light 0:  1   0   1  |   1
Light 1:  1   1   0  |   1
Light 2:  0   1   1  |   0
```

Our goal is to transform the matrix into row echelon form, which has a triangular shape.
We'll do this by eliminating variables systematically.

For our first operation, we want to eliminate `B0` from row 1.
We need to remove the 1 at position (row 1, col 0).
We do this by XORing row 1 with row 0:

```
Row 1:     1   1   0  |  1
Row 0:     1   0   1  |  1
XOR them:  0   1   1  |  0  ← new Row 1
```

After this operation, our matrix looks like:

```
          B0  B1  B2 | Target
Light 0:  1   0   1  |   1
Light 1:  0   1   1  |   0  ← changed
Light 2:  0   1   1  |   0
```

For our second operation, we want to eliminate `B1` from row 2:

```
Row 2:     0   1   1  |  0
Row 1:     0   1   1  |  0
XOR them:  0   0   0  |  0  ← new Row 2
```

Our final matrix in row echelon form is:

```
          B0  B1  B2 | Target
Light 0:  1   0   1  |   1
Light 1:  0   1   1  |   0
Light 2:  0   0   0  |   0
```

From our matrix, we can extract the relationships between variables.
Row 1 tells us that `B1 XOR B2 = 0`, therefore `B1 = B2`.
Row 0 tells us that `B0 XOR B2 = 1`, therefore `B0 = 1 XOR B2`.

Row 2 is `[0 0 0 | 0]`, which means `0 = 0`.
This is always true and tells us that `B2` is a free variable.
It can be either `0` or `1`, and we'll need to try both values to find which gives us the minimum number of button presses.
In case 1, if `B2 = 0`, then `B1 = B2 = 0` and `B0 = 1 XOR 0 = 1`.
This gives us the solution: press button 0 only, which is 1 press.
In case 2, if `B2 = 1`, then `B1 = B2 = 1` and `B0 = 1 XOR 1 = 0`.
This gives us the solution: press buttons 1 and 2, which is 2 presses.

## Now the coding implementation

The algorithm follows a clear structure.
First, we build an augmented matrix `[A | b]` where rows represent lights, columns represent buttons, and the last column holds the target states.
Second, we perform Gaussian elimination to find pivots in each column, eliminate those columns in other rows, and track which variables are free.
Third, we check for consistency by looking for contradictions like `[0 0 0 | 1]`.
Finally, we find the minimum solution by trying all combinations of free variables, using back substitution for other variables, and returning the solution with the fewest button presses.

```go
func solveButtonToggle(buttons [][]int, target []bool) int {
  // step 1: build matrix
  matrix := buildMatrix(buttons, target)

  // step 2: gaussian elimination
  freeVars := gaussianEliminate(matrix)

  // step 3: check consistency
  if !isConsistent(matrix) {
    return -1 // no solution marker
  }

  // step 4: find minimum
  minPresses := findMinSolution(matrix, len(buttons), freeVars)
  return minPresses
}
```

When eliminating a column, we perform a row XOR operation.
This is done by XORing each element of the rows together using the bitwise XOR operator.
For back substitution, we work through equations like `B0 XOR B2 = 1` where `B2` is already known.
We start with the target value and XOR it with all the other button states that appear in the equation, giving us the value for the pivot button.

See the code implementation for details.
