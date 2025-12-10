package main

import "fmt"

func processV1(filename string) (string, error) {
	// get input
	machines, err := readFile(filename)
	if err != nil {
		return "", err
	}

	// key idea: since pressing a button twice cancels out, we only need to press once or not at all
	// this reduces the problem to finding which subset of buttons to press

	totalMinPresses := 0
	for _, m := range machines {
		// fmt.Printf("Processing machine with %d lights and %d buttons\n", len(m.lightsReq), len(m.buttons))
		numButtons := len(m.buttons)
		minPresses := numButtons + 1

		// in this approach, we will brute force all combinations of button presses
		for mask := 0; mask < (1 << numButtons); mask++ { // << moves 1 binary to left by numButtons
			// which buttons are we pressing in this combination
			// e.g., mask = 5 = 101 -> press buttons 0 and 2
			var pressed []int
			for i := range numButtons {
				if mask&(1<<i) != 0 { // bitwise AND to check if ith bit is set
					// say mask = 5, numButtons = 3
					// i=0: 5 & (1<<0) = 101 & 001 = 1 != 0 -> press button 0
					// i=1: 5 & (1<<1) = 101 & 010 = 0 == 0 -> don't press button 1
					// i=2: 5 & (1<<2) = 101 & 100 = 4 != 0 -> press button 2
					pressed = append(pressed, i)
				}
			}

			// simulate the button presses
			result := m.simulatePresses(pressed)
			if lightsEqual(result, m.lightsReq) && len(pressed) < minPresses {
				minPresses = len(pressed)
				// fmt.Printf("Found solution with %d presses: buttons %v\n", len(pressed), pressed)
			}
		}

		totalMinPresses += minPresses
	}
	return fmt.Sprintf("Total minimum button presses: %d", totalMinPresses), nil
}

func (m *machine) simulatePresses(pressed []int) []bool {
	lights := make([]bool, len(m.lightsReq))
	// start with all off
	for _, btnIdx := range pressed {
		for _, lightIdx := range m.buttons[btnIdx] {
			lights[lightIdx] = !lights[lightIdx] // toggle the light
		}
	}
	return lights
}
