package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
)

// State represents a point on a grid and the the time T needed to reach that point
type State struct {
	P image.Point
	T int
}

func solve() (int, error) {

	// Read input file
	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// Create a map of points to their value
	// This can be clear ground "." or a blizzard up (^), down (v), left (<), or right (>).
	vall := map[image.Point]rune{}
	y := 0
	for fileScanner.Scan() {
		for x, r := range fileScanner.Text() {
			vall[image.Point{x, y}] = r
		}
		y++
	}

	// Find bounding rectables for all runes
	var bliz image.Rectangle
	for p := range vall {
		bliz = bliz.Union(image.Rectangle{p, p.Add(image.Point{1, 1})})
	}
	bliz.Min, bliz.Max = bliz.Min.Add(image.Point{1, 1}), bliz.Max.Sub(image.Point{1, 1})

	// bfs takes a start and end point to be traverses
	// A time that starts at 0 for part one to track time to traverse
	// returns the time to reach the end point
	bfs := func(start image.Point, end image.Point, time int) int {

		// The movement for each type of rune
		delta := map[rune]image.Point{
			'#': {0, 0}, '^': {0, -1}, '>': {1, 0}, 'v': {0, 1}, '<': {-1, 0},
		}

		// queue all possible states starting witht the start position
		queue := []State{{start, time}}
		// seen will track all points that have already been visited.
		seen := map[State]struct{}{queue[0]: {}}

		// Continue to search while there are still elements in the queue
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]

		loop:
			// Determine the next point on the grid by adding the movement delta to the point
			for _, d := range delta {
				next := State{cur.P.Add(d), cur.T + 1}
				// If the end is reached return the amount of time
				if next.P == end {
					return next.T
				}

				// If already seen skip this state
				if _, ok := seen[next]; ok {
					continue
				}

				// Check if point is outside the bounds or a wall
				if r, ok := vall[next.P]; !ok || r == '#' {
					continue
				}

				// Loop through all possible movement again to check if the next point to move to is a wall or not.
				// If not a wall then it is a new point that needs to be explored and is added to the queue.
				if next.P.In(bliz) {
					for r, d := range delta {
						if vall[next.P.Sub(d.Mul(next.T)).Mod(bliz)] == r {
							continue loop
						}
					}
				}

				seen[next] = struct{}{}
				queue = append(queue, next)
			}
		}
		return -1
	}

	// Start at the top-left corner of the bounding box.
	// End at the bottom-right corner
	start, end := bliz.Min.Sub(image.Point{0, 1}), bliz.Max.Sub(image.Point{1, 0})
	fmt.Println(bfs(start, end, 0))
	// For part 2 go from start to end then back to start and then to end again.
	fmt.Println(bfs(start, end, bfs(end, start, bfs(start, end, 0))))

	return -1, nil
}

func main() {
	_, err := solve()
	if err != nil {
		log.Fatal(err)
	}
}
