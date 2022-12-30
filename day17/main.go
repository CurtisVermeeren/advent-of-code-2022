package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"
)

func partOne() (int, error) {

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

	//cache := map[[2]int][2]int{}

	height, jet := 0, 0
	for i := 0; i < 1000000000000; i++ {
		if i == 2022 {
			return height, nil
		}

		/*
			k := [2]int{i % len(rocks), jet}
			if c, ok := cache[k]; ok {
				if n, d := 1000000000000-i, i-c[0]; n%d == 0 {
					return height + n/d*(height-c[1]), nil

				}
			}

		*/

		//cache[k] = [2]int{i, height}

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

	cache := map[[2]int][2]int{}

	height, jet := 0, 0
	for i := 0; i < 1000000000000; i++ {

		k := [2]int{i % len(rocks), jet}
		if c, ok := cache[k]; ok {
			if n, d := 1000000000000-i, i-c[0]; n%d == 0 {
				return height + n/d*(height-c[1]), nil

			}
		}

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
