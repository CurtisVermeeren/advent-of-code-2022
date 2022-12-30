package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

	// cave represents the 2D grid of objects in the cave
	// rock is "#" , air is "." , and sand will be "o"
	cave := make(map[Pair]rune)

	// maxY tracks the furhest rock from the top. Anything below this point is the "void"
	var maxY int

	// Generate the cave structure
	for fileScanner.Scan() {
		location := strings.Split(fileScanner.Text(), " -> ")
		for i := range location[:len(location)-1] {
			// Get each Pair of x,y for the origin and end location of the rock
			locationFrom := strings.Split(location[i], ",")
			locationTo := strings.Split(location[i+1], ",")

			// Get the x and y values for the start and end of the rock
			fromX, _ := strconv.Atoi(locationFrom[0])
			fromY, _ := strconv.Atoi(locationFrom[1])

			toX, _ := strconv.Atoi(locationTo[0])
			toY, _ := strconv.Atoi(locationTo[1])

			// Add the rock sections in the cave
			cave[Pair{toX, toY}] = '#'
			cave[Pair{fromX, fromY}] = '#'

			// Add all sections of rock between the start and end either in either a horizontal or vertical
			for fromX != toX || fromY != toY {
				cave[Pair{fromX, fromY}] = '#'
				switch {
				case fromX < toX:
					// The rock formation goes from left to right
					fromX++
				case fromY < toY:
					// The rock formation goes from top to bottom
					fromY++
				case fromX > toX:
					// The rock formation goes from right to left
					fromX--
				case fromY > toY:
					// The rock formation goes from bottom to top
					fromY--
				}
			}

			// Update the furthest furthest from the top to define the new edge of the "void"
			if fromY >= maxY {
				maxY = toY
			}

			if toY >= maxY {
				maxY = toY
			}
		}
	}

	voidReached := false // track if a unit of sand has fallen to the abyss
	sand := 0            // count the units of sand

	for !voidReached {
		newSand := Pair{500, 0}
		for {

			// If the sand has fallen past the last rock it will fall into the void
			// all other sand will continue into the void
			if newSand.y+1 > maxY {
				voidReached = true
				break
			}
			// directions is a slice of Pair that correspond to the 3 directions sand can fall. [0] = down, [1] = diagonal left, [2] = diagonal right
			directions := []Pair{{newSand.x, newSand.y + 1}, {newSand.x - 1, newSand.y + 1}, {newSand.x + 1, newSand.y + 1}}

			if cave[directions[0]] < '#' { // attempt to move straight down
				newSand.y++
			} else if cave[directions[1]] < '#' { // attempt to move SW
				newSand.x--
				newSand.y++
			} else if cave[directions[2]] < '#' { // attempt to move SE
				newSand.y++
				newSand.x++
			} else { // sand cannot move to falls into its final place
				cave[newSand] = 'o'
				sand++
				break
			}
		}
	}

	return sand, nil
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// cave represents the 2D grid of objects in the cave
	// rock is "#" , air is "." , and sand will be "o"
	cave := make(map[Pair]rune)

	// maxY tracks the furhest rock from the top. Anything below this point is the "void"
	var maxY int

	// track the maximum and minimum values of the floor so it can be filled with rocks
	var floorStart, floorEnd int

	// Generate the cave structure
	for fileScanner.Scan() {
		location := strings.Split(fileScanner.Text(), " -> ")
		for i := range location[:len(location)-1] {
			// Get each Pair of x,y for the origin and end location of the rock
			locationFrom := strings.Split(location[i], ",")
			locationTo := strings.Split(location[i+1], ",")

			// Get the x and y values for the start and end of the rock
			fromX, _ := strconv.Atoi(locationFrom[0])
			fromY, _ := strconv.Atoi(locationFrom[1])

			toX, _ := strconv.Atoi(locationTo[0])
			toY, _ := strconv.Atoi(locationTo[1])

			// Add the rock sections in the cave
			cave[Pair{toX, toY}] = '#'
			cave[Pair{fromX, fromY}] = '#'

			// Add all sections of rock between the start and end either in either a horizontal or vertical
			for fromX != toX || fromY != toY {
				cave[Pair{fromX, fromY}] = '#'
				switch {
				case fromX < toX:
					// The rock formation goes from left to right
					fromX++
				case fromY < toY:
					// The rock formation goes from top to bottom
					fromY++
				case fromX > toX:
					// The rock formation goes from right to left
					fromX--
				case fromY > toY:
					// The rock formation goes from bottom to top
					fromY--
				}
			}

			// Update the furthest furthest from the top to define the new edge of the "void"
			if fromY >= maxY {
				maxY = toY
			}
			if toY >= maxY {
				maxY = toY
			}

			// Update the width of the floor
			if toX > floorEnd {
				floorEnd = toX
			}
			if toX < floorStart {
				floorStart = toX
			}
			if fromX > floorEnd {
				floorEnd = fromX
			}
			if fromX < floorStart {
				floorStart = fromX
			}
		}
	}

	// fill in the entire floor with rocks.
	// use the maxY value as extra padding to imitate the floor being "infinite"
	for i := floorStart - maxY; i < floorEnd+maxY; i++ {
		cave[Pair{i, maxY + 2}] = '#'
	}

	sand := 0 // count the units of sand

	for {
		newSand := Pair{500, 0}

		// If the start point has been marked as sand it is blocked
		if cave[newSand] == 'o' {
			break
		}

		for {
			// directions is a slice of Pair that correspond to the 3 directions sand can fall. [0] = down, [1] = diagonal left, [2] = diagonal right
			directions := []Pair{{newSand.x, newSand.y + 1}, {newSand.x - 1, newSand.y + 1}, {newSand.x + 1, newSand.y + 1}}

			if cave[directions[0]] < '#' { // attempt to move straight down
				newSand.y++
			} else if cave[directions[1]] < '#' { // attempt to move SW
				newSand.x--
				newSand.y++
			} else if cave[directions[2]] < '#' { // attempt to move SE
				newSand.y++
				newSand.x++
			} else { // sand cannot move to falls into its final place
				cave[newSand] = 'o'
				sand++
				break
			}
		}
	}

	return sand, nil
}

func main() {
	sand, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sand)

	sand, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sand)
}
