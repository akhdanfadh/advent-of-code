package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
)

func main() {
	// program input
	filename := flag.String("f", "", "input file name (required)")
	version := flag.String("v", "1", "logic version")
	connection := flag.Int("c", 10, "number of connections for the logic")
	flag.Parse()

	// main logic
	var result string
	var err error
	switch *version {
	case "1":
		result, err = processV1(*filename, *connection)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown version %s\n", *version)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}

func processV1(filename string, connection int) (string, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close() // error ignored (file only for reading)

	// read points
	var x, y, z int
	points := []point{}
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		_, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		if err != nil {
			return "", err
		}
		points = append(points, point{id: i, x: x, y: y, z: z})
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	fmt.Printf("Read %d points\n", len(points))

	// heapify while calculating distances
	// why heap? note that we don't need exactly sorted list, just N shortest for now
	nthShortest := &pairHeap{}
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dist := calcDist(points[i], points[j])
			nthShortest.push(pair{p1: points[i], p2: points[j], dist: dist})
			if nthShortest.len() > connection {
				nthShortest.pop()
			}
		}
	}
	fmt.Printf("Calculated distances, kept %d shortest\n", nthShortest.len())

	// now create the circuit from the pairs
	mapCircuitToPoints := make(map[int][]int, nthShortest.len())
	mapPointsToCircuit := make(map[int]int, nthShortest.len())
	circuitID := 0
	for _, pair := range *nthShortest {
		p1id, p2id := pair.p1.id, pair.p2.id
		c1id, p1used := mapPointsToCircuit[p1id]
		c2id, p2used := mapPointsToCircuit[p2id]

		// if both points used...
		if p1used && p2used {
			if c1id == c2id {
				continue // same circuit, skip
			} else { // otherwise merge circuits
				// delete circuit c2id, move all points to c1id
				movingPoints := mapCircuitToPoints[c2id]
				for _, mp := range movingPoints {
					mapPointsToCircuit[mp] = c1id
				}
				mapCircuitToPoints[c1id] = append(mapCircuitToPoints[c1id], movingPoints...)
				delete(mapCircuitToPoints, c2id)
				continue

			}
		}

		// if one of the points used...
		if p1used {
			mapPointsToCircuit[p2id] = c1id
			mapCircuitToPoints[c1id] = append(mapCircuitToPoints[c1id], p2id)
			continue
		}
		if p2used {
			mapPointsToCircuit[p1id] = c2id
			mapCircuitToPoints[c2id] = append(mapCircuitToPoints[c2id], p1id)
			continue
		}

		// if none used, create new circuit
		if !p1used && !p2used {
			mapPointsToCircuit[p1id] = circuitID
			mapPointsToCircuit[p2id] = circuitID
			mapCircuitToPoints[circuitID] = []int{p1id, p2id}
			circuitID++
		}
	}

	// now count how many points are in circuits
	countSet := make(map[int]struct{}, len(mapCircuitToPoints))
	for _, points := range mapCircuitToPoints {
		countSet[len(points)] = struct{}{}
	}

	// sort the counts
	counts := make([]int, 0, len(countSet))
	for k := range countSet {
		counts = append(counts, k)
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i] > counts[j]
	})

	result := counts[0] * counts[1] * counts[2]
	return fmt.Sprintf("%d", result), nil
}

type (
	point struct {
		id      int
		x, y, z int
	}
	pair struct {
		p1, p2 point
		dist   float64
	}
	pairHeap []pair
)

func (h *pairHeap) len() int {
	return len(*h)
}

func (h *pairHeap) push(p pair) {
	// add new element at the end
	*h = append(*h, p)
	// percolate up while larger than parent
	i := h.len() - 1
	par := (i - 1) / 2
	for i > 0 && (*h)[i].dist > (*h)[par].dist {
		(*h)[i], (*h)[par] = (*h)[par], (*h)[i]
		i = par
		par = (i - 1) / 2
	}
}

func (h *pairHeap) pop() pair {
	if h.len() == 0 {
		return pair{}
	}
	p := (*h)[0]         // pop root
	size := h.len() - 1  // new size
	(*h)[0] = (*h)[size] // move last to root
	*h = (*h)[:size]     // shrink slice
	// percolate down while smaller than children
	i := 0
	for 2*i+1 < size { // while there is at least one child (left)
		par, lef, rig := i, 2*i+1, 2*i+2
		if lef < size && (*h)[par].dist < (*h)[lef].dist {
			par = lef
		}
		if rig < size && (*h)[par].dist < (*h)[rig].dist {
			par = rig
		}
		if par == i { // no swap happened
			break
		}
		(*h)[i], (*h)[par] = (*h)[par], (*h)[i] // now swap
		i = par
	}
	return p
}

func calcDist(a, b point) float64 {
	return float64((a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y) + (a.z-b.z)*(a.z-b.z))
}
