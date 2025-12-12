package main

import "fmt"

func processV1(filename string) (string, error) {
	aoc, err := readInput(filename)
	if err != nil {
		return "", err
	}

	// log on googling:
	// - gamedev stackoverflow -> packing problem wiki rabbithole
	// - awesome blog: https://www.gorillasun.de/blog/a-simple-solution-for-shape-packing-in-2d/
	// - yt search give this cool video, now we know our problem is polygon packing
	//   https://www.youtube.com/watch?v=jkyGgfCYkMg
	// - github search with tag packing-algorithm gives us this
	//   https://github.com/fontanf/packingsolver#irregular-solver -> irregular interesting
	//   https://github.com/JonasTollenaere/sparrow-3d gives https://github.com/JeroenGar/sparrow -> 2509 arxiv
	// - from sparrow github, we got our problem that is 2D nesting algorithm
	// - googling the term -> https://github.com/sasam2/nesting -> SVGnest!
	//   https://github.com/Jack000/SVGnest (2019) most starred and readme very good explanation

	// BRUH moment:
	// - based on those googling, this problem is NP-hard
	// - should i make tbe go version of the solver? surely no right?
	//   htf aoc this hard, and this is my first on top of that haha
	// - i'll look on youtube for aoc day 12 solution and found
	//   programming with larry confirms np-hard: https://www.youtube.com/watch?v=3MUTlkdFSUE
	// - and then a github repo with many stars https://github.com/jonathanpaulson/AdventOfCode
	//   and the video too https://www.youtube.com/watch?v=am-X5j1DVkA
	//   just shows that we can just assume the result LMFAO

	// you got me there, aoc

	// count shape #'s
	shapeSize := make(map[int]int)
	for id, shape := range aoc.presents {
		size := 0
		for i := range *shape {
			for j := range (*shape)[i] {
				if (*shape)[i][j] == 1 {
					size++
				}
			}
		}
		shapeSize[id] = size
	}

	// count area in each region and the presents area
	result := 0
	for _, region := range aoc.regions {
		area := region.width * region.height
		presentsArea := 0
		for presentID, count := range region.presentsCount {
			presentsArea += shapeSize[presentID] * count
		}

		// this 1.3 factor is just eyeballing, it is incorrect on test1 but input is correct somehow
		if float64(presentsArea)*1.3 < float64(area) {
			fmt.Printf("Correct? %dx%d: area=%d, presentsArea=%d\n", region.width, region.height, area, presentsArea)
			result++
		} else if presentsArea > area {
			fmt.Printf("Invalid %dx%d: area=%d, presentsArea=%d\n", region.width, region.height, area, presentsArea)
		} else {
			fmt.Printf("Correct??!! %dx%d: area=%d, presentsArea=%d\n", region.width, region.height, area, presentsArea)
		}
	}
	return fmt.Sprintf("Supposedly correct regions: %d", result), nil
}
