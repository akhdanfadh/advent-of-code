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
		if *connection < 1 {
			fmt.Fprintf(os.Stderr, "Error: connection must be >= 1, got %d\n", *connection)
			os.Exit(1)
		}
		result, err = processV1(*filename, *connection)
	case "1a":
		result, err = processV1a(*filename, *connection)
	case "2":
		result, err = processV2(*filename)
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
	// read points from file
	points, err := readPointsFromFile(filename)
	if err != nil {
		return "", err
	}

	// heapify while calculating distances
	// why heap? note that we don't need exactly sorted list, just N shortest for now
	pairs := buildPairHeap(points, connection)

	// now create the circuit from the pairs
	mapCircuitToPoints := make(map[int][]int, pairs.len())
	mapPointsToCircuit := make(map[int]int, pairs.len())
	circuitID := 0
	for _, pair := range *pairs {
		p1id, p2id := pair.p1.id, pair.p2.id
		c1id, p1used := mapPointsToCircuit[p1id]
		c2id, p2used := mapPointsToCircuit[p2id]

		// if both points used...
		if p1used && p2used {
			if c1id == c2id {
				continue // same circuit, skip
			}

			// different circuit, merge:
			// delete circuit c2id, move all points to c1id
			movingPoints := mapCircuitToPoints[c2id]
			for _, mp := range movingPoints {
				mapPointsToCircuit[mp] = c1id
			}
			mapCircuitToPoints[c1id] = append(mapCircuitToPoints[c1id], movingPoints...)
			delete(mapCircuitToPoints, c2id)
			continue
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
		mapPointsToCircuit[p1id] = circuitID
		mapPointsToCircuit[p2id] = circuitID
		mapCircuitToPoints[circuitID] = []int{p1id, p2id}
		circuitID++
	}

	// count how many points with no pairing
	singlePoints := 0
	for _, p := range points {
		if _, used := mapPointsToCircuit[p.id]; !used {
			singlePoints++
		}
	}

	// now compute sizes of all circuits
	sizes := make([]int, 0, len(mapCircuitToPoints))
	for _, points := range mapCircuitToPoints {
		sizes = append(sizes, len(points))
	}
	// and sort descending
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})
	// and add single points as size 1 circuits
	for i := 0; i < singlePoints; i++ {
		sizes = append(sizes, 1)
	}

	result := sizes[0] * sizes[1] * sizes[2]
	return fmt.Sprintf("%d", result), nil
}

func processV1a(filename string, connection int) (string, error) {
	// read points from file
	points, err := readPointsFromFile(filename)
	if err != nil {
		return "", err
	}

	// validate input
	if len(points) < 3 {
		return "", fmt.Errorf("need minimum 3 points, got %d", len(points))
	}
	maxPairs := len(points) * (len(points) - 1) / 2
	if connection > maxPairs {
		connection = maxPairs
	}

	// heapify while calculating distances
	// why heap? note that we don't need exactly sorted list, just N shortest for now
	pairs := buildPairHeap(points, connection)

	// build circuits with disjoint set
	circuits := initCircuits(len(points))
	for _, pair := range *pairs {
		circuits.union(pair.p1.id, pair.p2.id)
	}

	// get each circuit sizes and sort descending
	sizes := circuits.getRootSizes()
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	result := sizes[0] * sizes[1] * sizes[2]
	return fmt.Sprintf("%d", result), nil
}

func processV2(filename string) (string, error) {
	// get points from file
	points, err := readPointsFromFile(filename)
	if err != nil {
		return "", err
	}

	// build pairs as heap but no min size
	pairs := &pairMinHeap{}
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dist := calcDist(points[i], points[j])
			pairs.push(pair{p1: points[i], p2: points[j], dist: dist})
		}
	}

	// process shortest pair one by one until all points connected
	circuits := initCircuits(len(points)) // disjoint set
	allConnected := false
	var x1, x2 int
	for pairs.len() > 0 {
		shortest := pairs.pop()
		fmt.Printf("Connecting point %d and %d with distance %.2f\n", shortest.p1.id, shortest.p2.id, shortest.dist)
		circuits.union(shortest.p1.id, shortest.p2.id)

		// check if all connected: all points have the same root
		root := circuits.find(0)
		allConnected = true
		for id := 1; id < len(circuits.parent); id++ {
			if circuits.find(id) != root {
				allConnected = false
				break
			}
		}

		if allConnected {
			x1 = shortest.p1.x
			x2 = shortest.p2.x
			break
		}
	}

	if allConnected {
		return fmt.Sprintf("Multiplying the X coordinates of the last two junction boxes got %d", x1*x2), nil
	}
	return "", fmt.Errorf("could not connect all points")
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
)

func readPointsFromFile(filename string) ([]point, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		points = append(points, point{id: i, x: x, y: y, z: z})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return points, nil
}

func buildPairHeap(points []point, connection int) *pairHeap {
	pairs := &pairHeap{}
	for i := 0; i < len(points); i++ {
		for j := i + 1; j < len(points); j++ {
			dist := calcDist(points[i], points[j])
			pairs.push(pair{p1: points[i], p2: points[j], dist: dist})
			if pairs.len() > connection {
				pairs.pop()
			}
		}
	}
	return pairs
}

func calcDist(a, b point) float64 {
	// note that this is not the actual distance, but squared distance just for comparison
	return float64((a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y) + (a.z-b.z)*(a.z-b.z))
}
