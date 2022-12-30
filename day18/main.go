package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

// Point represtents a 3d point
type Point struct {
	x, y, z int
}

// Add point p and q together
func (p Point) Add(q Point) Point {
	return Point{p.x + q.x, p.y + q.y, p.z + q.z}
}

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	lava := map[Point]struct{}{}
	min := Point{math.MaxInt, math.MaxInt, math.MaxInt}
	max := Point{math.MinInt, math.MinInt, math.MinInt}

	for fileScanner.Scan() {
		var p Point
		fmt.Sscanf(fileScanner.Text(), "%d,%d,%d", &p.x, &p.y, &p.z)
		lava[p] = struct{}{}

		min = Point{Min(min.x, p.x), Min(min.y, p.y), Min(min.z, p.z)}
		max = Point{Max(max.x, p.x), Max(max.y, p.y), Max(max.z, p.z)}
	}

	min = min.Add(Point{-1, -1, -1})
	max = max.Add(Point{1, 1, 1})

	delta := []Point{
		{-1, 0, 0}, {0, -1, 0}, {0, 0, -1},
		{1, 0, 0}, {0, 1, 0}, {0, 0, 1},
	}

	surface := 0
	for p := range lava {
		for _, d := range delta {
			if _, ok := lava[p.Add(d)]; !ok {
				surface++
			}
		}
	}

	return surface, nil
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	lava := map[Point]struct{}{}
	min := Point{math.MaxInt, math.MaxInt, math.MaxInt}
	max := Point{math.MinInt, math.MinInt, math.MinInt}

	for fileScanner.Scan() {
		var p Point
		fmt.Sscanf(fileScanner.Text(), "%d,%d,%d", &p.x, &p.y, &p.z)
		lava[p] = struct{}{}

		min = Point{Min(min.x, p.x), Min(min.y, p.y), Min(min.z, p.z)}
		max = Point{Max(max.x, p.x), Max(max.y, p.y), Max(max.z, p.z)}
	}

	min = min.Add(Point{-1, -1, -1})
	max = max.Add(Point{1, 1, 1})

	delta := []Point{
		{-1, 0, 0}, {0, -1, 0}, {0, 0, -1},
		{1, 0, 0}, {0, 1, 0}, {0, 0, 1},
	}

	queue := []Point{min}
	visited := map[Point]struct{}{min: {}}

	surface := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, d := range delta {
			next := current.Add(d)

			if _, ok := lava[next]; ok {
				surface++
			} else if _, ok := visited[next]; !ok && next.x >= min.x && next.x <= max.x && next.y >= min.y && next.y <= max.y && next.z >= min.z && next.z <= max.z {
				visited[next] = struct{}{}
				queue = append(queue, next)
			}
		}
	}

	return surface, nil
}

// Min returns the smaller value of a and b
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger value of a and b
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	surface, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(surface)

	surface, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(surface)
}
