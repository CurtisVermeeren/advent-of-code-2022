package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"
)

func partOne() (int, error) {

	// read the input file
	file, err := os.ReadFile("./input.txt")
	if err != nil {
		return -1, err
	}
	jets := strings.TrimSpace(string(file))

	// rocks is a slice which contains slices of points.
	// each slive of points represents a rock shape
	rocks := [][]image.Point{
		{{0, 0}, {1, 0}, {2, 0}, {3, 0}},         // Straight horizontal line
		{{1, 2}, {0, 1}, {1, 1}, {2, 1}, {1, 0}}, // Plus shape
		{{2, 2}, {2, 1}, {0, 0}, {1, 0}, {2, 0}}, // Backwards L shape
		{{0, 3}, {0, 2}, {0, 1}, {0, 0}},         // Straight vertical line
		{{0, 1}, {1, 1}, {0, 0}, {1, 0}},         // Square of size 2
	}

	// grid holds all positions of rocks
	grid := map[image.Point]struct{}{}

	// func move is used to move a slice of points by a given delta value
	// returns false if the new position is already occupied or outside the bounds of the grid.
	// returns true if the position was updated and changes the value of rock
	move := func(rock []image.Point, delta image.Point) bool {
		newRock := make([]image.Point, len(rock))
		for i, p := range rock {
			// Shift the rock by delta
			p = p.Add(delta)
			// If the new position is already in the grid or is outside the bounds return false
			if _, ok := grid[p]; ok || p.X < 0 || p.X >= 7 || p.Y < 0 {
				return false
			}
			newRock[i] = p
		}
		copy(rock, newRock)
		return true
	}

	// The current height of the tower and the index of the jet
	height, jet := 0, 0
	// calculate 2022 rocks falling as stated in the question
	for i := 0; i <= 2022; i++ {

		// Get the current rock shape that should be falling and place it at the current starting point
		rock := []image.Point{}
		for _, p := range rocks[i%len(rocks)] {
			rock = append(rock, p.Add(image.Point{2, height + 3}))
		}

		// Move the rock until it's in a resting position
		for {

			// Move the rock according to the jet direction
			move(rock, image.Point{int(jets[jet]) - int('='), 0})
			jet = (jet + 1) % len(jets)

			// Attempt to move the rock down 1 position
			// If the move is unsuccessful then the rock has reached a resting point
			// update all positions of the rock in the grid and update the height of the tower
			if !move(rock, image.Point{0, -1}) {
				for _, p := range rock {
					grid[p] = struct{}{}
					if p.Y+1 > height {
						height = p.Y + 1
					}
				}
				break
			}
		}
	}

	return height, nil
}

func partTwo() (int, error) {
	file, err := os.ReadFile("./input.txt")
	if err != nil {
		return -1, err
	}
	jets := strings.TrimSpace(string(file))

	rocks := [][]image.Point{
		{{0, 0}, {1, 0}, {2, 0}, {3, 0}},
		{{1, 2}, {0, 1}, {1, 1}, {2, 1}, {1, 0}},
		{{2, 2}, {2, 1}, {0, 0}, {1, 0}, {2, 0}},
		{{0, 3}, {0, 2}, {0, 1}, {0, 0}},
		{{0, 1}, {1, 1}, {0, 0}, {1, 0}},
	}

	grid := map[image.Point]struct{}{}
	move := func(rock []image.Point, delta image.Point) bool {
		newRock := make([]image.Point, len(rock))
		for i, p := range rock {
			p = p.Add(delta)
			if _, ok := grid[p]; ok || p.X < 0 || p.X >= 7 || p.Y < 0 {
				return false
			}
			newRock[i] = p
		}
		copy(rock, newRock)
		return true
	}

	// cache maps two pairs of integers to eachother
	// The first pair is the current iteration and the current index in jets
	// The second pair represents the iteration and the height
	cache := map[[2]int][2]int{}

	height, jet := 0, 0
	// Changed to 1000000000000 rocks for the second part
	for i := 0; i < 1000000000000; i++ {

		// k is used to check the cache. It is the curret iteration mod length of rocks, and the current index of jets
		k := [2]int{i % len(rocks), jet}
		// If the key is found in the cache calculate the difference between the current iteration and the iteration stored in the cache
		// If the difference is a multiple of the number of iterations between the current and the cached a pattern is found
		// Return the height at the current iteration plus
		// the difference divided by the number of iterations between the current and the cached iteration multiplied by the difference in height between the current iteration and the cached iteration
		if c, ok := cache[k]; ok {
			if n, d := 1000000000000-i, i-c[0]; n%d == 0 {
				return height + n/d*(height-c[1]), nil

			}
		}

		// If the key was not found. Cache the current iteration and continue to move rocks like part one
		cache[k] = [2]int{i, height}

		rock := []image.Point{}
		for _, p := range rocks[i%len(rocks)] {
			rock = append(rock, p.Add(image.Point{2, height + 3}))
		}

		for {
			move(rock, image.Point{int(jets[jet]) - int('='), 0})
			jet = (jet + 1) % len(jets)

			if !move(rock, image.Point{0, -1}) {
				for _, p := range rock {
					grid[p] = struct{}{}
					if p.Y+1 > height {
						height = p.Y + 1
					}
				}
				break
			}
		}
	}

	return -1, nil
}

func main() {
	height, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(height)

	height, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(height)
}
