package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

// Pair represents a pair of x and y coordinates
type Pair struct {
	x, y int
}

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// heightmap represents and x y grid of elevation values
	heightmap := make([][]rune, 0)

	// start, end mark the start position and the end destination
	var start, end Pair

	// Create the heightmap grid
	for fileScanner.Scan() {
		// currentRow is a row in the heightmap
		var currentRow []rune
		for index, elevation := range fileScanner.Text() {

			// When the start location is found set that pair. Adjust the elevation defined in the problem
			// len(heightmap) is also the current row number
			if elevation == 'S' {
				start = Pair{index, len(heightmap)}
				elevation = 'a'
			}

			// When the end location is found set that pair. Adjust the elevation defined in the problem
			if elevation == 'E' {
				end = Pair{index, len(heightmap)}
				elevation = 'z'
			}

			currentRow = append(currentRow, elevation)
		}

		heightmap = append(heightmap, currentRow)
	}

	// toVisit holds the locations to be visited
	toVisit := []Pair{start}
	// visited is a map of all locations that have already been visited
	visited := make(map[Pair]bool)
	// distanceTravelled maps each Pair of x,y coordinates to the distance from the starting point
	distanceTravelled := map[Pair]int{start: 0}

	heightmapHeight := len(heightmap)
	heightmapWidth := len(heightmap[0])

	for {
		// The current location is the one that is visited next
		// After the start location all currentLocation will be the location with the shortest distanceTravelled
		currentLocation := toVisit[0]
		// mark the current location as visited
		visited[currentLocation] = true
		// remove the current location from the locations to be visited
		toVisit = toVisit[1:]

		// If the end location is found return the distance from the start to the end
		if currentLocation == end {
			return distanceTravelled[end], nil
		}

		// For each location check the 4 neighbour locations to that point
		neighbours := []Pair{{1, 0}, {0, -1}, {-1, 0}, {0, 1}}
		for _, nearby := range neighbours {
			moveX, moveY := nearby.x, nearby.y
			nextLocation := Pair{currentLocation.x + moveX, currentLocation.y + moveY}

			// Check that the nextLocation hasn't been visited, is within the heightmap bounds, And that the elevation of the destination is no more than 1 greater than the current location so that it is a valid possible location
			if !visited[nextLocation] && nextLocation.x >= 0 && nextLocation.y >= 0 && nextLocation.x < heightmapWidth && nextLocation.y < heightmapHeight && (heightmap[nextLocation.y][nextLocation.x]-heightmap[currentLocation.y][currentLocation.x] <= 1) {

				// If the nextLocation has a distance of 0 from the start it hasn't been visited. Set its distanceTravelled from the start to one greater than the previous location
				// Add this neighbour as a location that needs to be visited
				if distanceTravelled[nextLocation] == 0 {
					toVisit = append(toVisit, nextLocation)
					distanceTravelled[nextLocation] = distanceTravelled[currentLocation] + 1
				}
			}
		}

		// Sort all locations to be visited from shortest distance to longest
		sort.Slice(toVisit, func(x, y int) bool {
			return distanceTravelled[toVisit[x]] < distanceTravelled[toVisit[y]]
		})
	}
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// heightmap represents and x y grid of elevation values
	heightmap := make([][]rune, 0)

	// start, end mark the start position and the end destination
	var end Pair

	// lowElectationsStarts holds all locations that have an elevation of a and are potential start locations for the trail
	var lowElevationStarts []Pair

	// Create the heightmap grid
	for fileScanner.Scan() {
		// currentRow is a row in the heightmap
		var currentRow []rune
		for index, elevation := range fileScanner.Text() {

			// When the start location is found set that pair. Adjust the elevation defined in the problem
			// len(heightmap) is also the current row number
			if elevation == 'S' || elevation == 'a' {
				elevation = 'a'
				lowElevationStarts = append(lowElevationStarts, Pair{index, len(heightmap)})

			}

			// When the end location is found set that pair. Adjust the elevation defined in the problem
			if elevation == 'E' {
				end = Pair{index, len(heightmap)}
				elevation = 'z'
			}

			currentRow = append(currentRow, elevation)
		}

		heightmap = append(heightmap, currentRow)
	}

	var shortestPath int
	heightmapHeight := len(heightmap)
	heightmapWidth := len(heightmap[0])

	for _, startLocation := range lowElevationStarts {
		visited := make(map[Pair]bool)
		toVisit := []Pair{startLocation}
		distanceTravelled := map[Pair]int{startLocation: 0}

		for {
			// If no unvisited nodes remain then this startLocation cannot reach the end location
			if len(toVisit) == 0 {
				break
			}

			currentLocation := toVisit[0]
			visited[currentLocation] = true
			toVisit = toVisit[1:]

			// If the end location is found check what the current shortest distance found is
			if currentLocation == end {
				// If the shortest distance is 0 this is the first complete path.
				// If this path is shorter than the previous shortest path update the length of the shortest path
				if distanceTravelled[end] < shortestPath || shortestPath == 0 {
					shortestPath = distanceTravelled[end]
				}
				// move to the next startLocation
				break
			}

			// Checking neighbours is the same as part 1
			neighbours := []Pair{{1, 0}, {0, -1}, {-1, 0}, {0, 1}}
			for _, nearby := range neighbours {
				moveX, moveY := nearby.x, nearby.y
				nextLocation := Pair{currentLocation.x + moveX, currentLocation.y + moveY}

				if !visited[nextLocation] && nextLocation.x >= 0 && nextLocation.y >= 0 && nextLocation.x < heightmapWidth && nextLocation.y < heightmapHeight && (heightmap[nextLocation.y][nextLocation.x]-heightmap[currentLocation.y][currentLocation.x] <= 1) {

					if distanceTravelled[nextLocation] == 0 {
						toVisit = append(toVisit, nextLocation)
						distanceTravelled[nextLocation] = distanceTravelled[currentLocation] + 1
					}
				}
			}

			// Sorting by shortest distanceTravelled is the same as part 1
			sort.Slice(toVisit, func(x, y int) bool {
				return distanceTravelled[toVisit[x]] < distanceTravelled[toVisit[y]]
			})
		}
	}

	return shortestPath, nil
}

func main() {
	distance, err := partOne()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(distance)

	distance, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(distance)

}
