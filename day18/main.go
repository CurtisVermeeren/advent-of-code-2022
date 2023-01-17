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

	// lava is a map of all points found from the input file.
	// these are 3D points with x, y, z coordinates.
	// min and max store the minimum and maximum x, y, z values found in the input to create a bounding box around the 3D shape
	lava := map[Point]struct{}{}
	min := Point{math.MaxInt, math.MaxInt, math.MaxInt}
	max := Point{math.MinInt, math.MinInt, math.MinInt}

	// Scan all lines of input. Each input line is a set of x, y, z coordinates indicating a block of "lava droplet"
	// Check for minimum and maximum values of x, y, z
	for fileScanner.Scan() {
		var p Point
		fmt.Sscanf(fileScanner.Text(), "%d,%d,%d", &p.x, &p.y, &p.z)
		lava[p] = struct{}{}

		min = Point{Min(min.x, p.x), Min(min.y, p.y), Min(min.z, p.z)}
		max = Point{Max(max.x, p.x), Max(max.y, p.y), Max(max.z, p.z)}
	}

	// Expand the bounding box by 1 on all sides
	min = min.Add(Point{-1, -1, -1})
	max = max.Add(Point{1, 1, 1})

	// These delta points are the 6 locations that surround a cube
	delta := []Point{
		{-1, 0, 0}, {0, -1, 0}, {0, 0, -1},
		{1, 0, 0}, {0, 1, 0}, {0, 0, 1},
	}

	// Track the total surface area of the droplet with surface
	// For each cube of lava droplet we check all 6 locations next to it.
	// If a location next to a lava droplet is found without another cube then we know this is an exposed face of the cube and we can increment the total surface area.
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

	// The same as partOne
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

	// queue will be a slice of all Points that need to be visited in the search.
	// We don't want to visit the same Point more than once so we track all Points that have been visited.
	// We know min is outside the lava structure as we subtracted 1 from all its coordinates when creating the bounding box
	queue := []Point{min}
	visited := map[Point]struct{}{min: {}}

	surface := 0

	// Start at the minimum point
	// For each point in the queue check all neighbours of the point for an adjacent lava block
	// If am adjacent point is part of the lava defined by the input then increment the surface counter.
	// If any point is not part of the lava, has not been visited already, and is within the bounds of the 3D grid it is added to the queue and marked as visited.
	// Because we started with the min Point thats outside the lava cube structure we traverse through all empty non-lava cubes around the lava structure and increment the surface area when when a lava cube is reached.
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
