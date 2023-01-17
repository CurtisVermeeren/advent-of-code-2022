package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"reflect"
)

func solve() (int, error) {

	// Read the input file
	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// Build the grid of Elf positions marked by #
	grid := map[image.Point]struct{}{}
	y := 0
	for fileScanner.Scan() {
		for x, r := range fileScanner.Text() {
			if r == '#' {
				grid[image.Point{x, y}] = struct{}{}
			}
		}
		y++
	}

	// sides is a slice of 4 sets of points
	// These points correspond to the 3 directions to check before moving N, E, S or W
	// These directions are defined in the problem. Ex if there is no Elf in the N, NE, NW the Elf proposes moving north one step.
	sides := [][]image.Point{
		{{0, -1}, {1, -1}, {-1, -1}},
		{{0, 1}, {1, 1}, {-1, 1}},
		{{-1, 0}, {-1, -1}, {-1, 1}},
		{{1, 0}, {1, -1}, {1, 1}},
	}

	//
	for i := 0; ; i++ {

		// prop maps the proposed new position for each point
		prop := map[image.Point]image.Point{}
		// count stores the number of points that are proposed to occupy the same position in the grid
		count := map[image.Point]int{}

		// Iterate over all points on the current grid
		// Check each neighbour for all points. Store the number of neighbouring points in each direction of the neigh map
		for p := range grid {
			neigh := map[int]int{}
			for i := range sides {
				for _, q := range sides[i] {
					if _, ok := grid[p.Add(q)]; ok {
						neigh[i]++
					}
				}
			}

			// If a point has no neighbouring points skip to the next
			if len(neigh) == 0 {
				continue
			}

			// For each point that has at least one neighbour iterate over all possible directions and check if there
			// are neighbouring points in that direction. If a direction is found with no neighbours the loop assigns the points
			// position in the new grid as the first point in the corresponding direction in the sides array.
			// The count map is incremented for that point
			for d := 0; d < len(sides); d++ {
				if dir := (i + d) % len(sides); neigh[dir] == 0 {
					prop[p] = p.Add(sides[dir][0])
					count[prop[p]]++
					break
				}
			}
		}

		// newGrid is used to store the new positions of the points
		// All points in the current grid are then iterated over. If a point has a proposed position in the prop map
		// and the count for that point in the count map is 1 then the new position is assigned for the newGrid
		newGrid := map[image.Point]struct{}{}
		for p := range grid {
			if _, ok := prop[p]; ok && count[prop[p]] == 1 {
				p = prop[p]
			}
			newGrid[p] = struct{}{}
		}

		// If the 10th iteration is reached then calculate the size of the bounding box of the newGrid
		// and the number of points within that bounding box and output the difference to determine the number of empty tiles
		if i == 9 {
			var r image.Rectangle
			for p := range newGrid {
				r = r.Union(image.Rectangle{p, p.Add(image.Point{1, 1})})
			}
			fmt.Println(r.Dx()*r.Dy() - len(newGrid))
		}

		// When the new grid is equal to the current grid then no elves have moved.
		if reflect.DeepEqual(grid, newGrid) {
			fmt.Println(i + 1)
			break
		}

		grid = newGrid
	}

	return -1, nil
}

func main() {
	answer, err := solve()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(answer)
}
