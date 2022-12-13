package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Pair represents a pair of coordinates on a grid of x and y
type Pair struct {
	x, y interface{}
}

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	var forest [][]rune

	// Create the forest reading one row (line) of trees at a time
	for fileScanner.Scan() {
		row := []rune{}
		for _, tree := range fileScanner.Text() {
			row = append(row, tree)
		}
		forest = append(forest, row)
	}

	// height and width of the forest
	forestHeight := len(forest)
	forestWidth := len(forest[0])

	// maxLeft and maxRight are slices that track the max value from the left and right for each row of trees
	// They are indexed using i so that each index represents the max value of row i
	maxLeft := make([]rune, len(forest))
	maxRight := make([]rune, len(forest))

	// Check the horizontal view by creating a map of all visible trees
	isVisible := make(map[Pair]bool)

	// go line by line
	for i := 0; i < forestHeight; i++ {
		// go across each line
		for j := 0; j < forestWidth; j++ {
			if j == 0 {
				// All trees on the outmost left and right edges of the forest are visible
				maxLeft[i] = forest[i][j]
				maxRight[i] = forest[i][forestWidth-1]

				// Set the outer edges as visible
				isVisible[Pair{i, j}] = true
				isVisible[Pair{i, forestWidth - 1}] = true

				continue
			}

			if forest[i][j] > maxLeft[i] {
				// Current tree is taller than the previous max on the left making it visible
				isVisible[Pair{i, j}] = true
				maxLeft[i] = forest[i][j]
			}

			if forest[i][forestWidth-1-j] > maxRight[i] {
				// Current tree is taller than the previous max on the right making it visible
				isVisible[Pair{i, forestWidth - 1 - j}] = true
				maxRight[i] = forest[i][forestWidth-1-j]
			}
		}
	}

	// maxTop and maxBottom are slices that track the max value from the top and bottom for each column of trees
	// They are indexed using j so that each index represents the max value of column j
	maxTop := make([]rune, len(forest))
	maxBottom := make([]rune, len(forest))

	for i := 0; i < forestHeight; i++ {
		for j := 0; j < forestWidth; j++ {
			if j == 0 {
				// All tree on the top and bottom edges of the forest are visible
				maxTop[j] = forest[i][j]
				maxBottom[j] = forest[forestHeight-1][j]

				isVisible[Pair{i, j}] = true
				isVisible[Pair{forestHeight - 1, j}] = true

				continue
			}

			if forest[i][j] > maxTop[j] {
				// Current tree is taller than the previous max above it making it visible
				isVisible[Pair{i, j}] = true
				maxTop[j] = forest[i][j]
			}

			if forest[forestHeight-1-i][j] > maxBottom[j] {
				// Current tree is taller than the previous max below it makint it visible
				isVisible[Pair{forestHeight - 1 - i, j}] = true
				maxBottom[j] = forest[forestHeight-1-i][j]
			}
		}
	}

	// Return the count of the number of trees in the isVisible map
	return len(isVisible), nil
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	var forest [][]rune

	// Create the forest reading one row (line) of trees at a time
	for fileScanner.Scan() {
		row := []rune{}
		for _, tree := range fileScanner.Text() {
			row = append(row, tree)
		}
		forest = append(forest, row)
	}

	var highestScenicScore int

	// height and width of the forest
	forestHeight := len(forest)
	forestWidth := len(forest[0])

	// Check the scencic score of each tree
	for i := 0; i < forestHeight; i++ {
		for j := 0; j < forestWidth; j++ {
			scenicScore := calcDown(forest, i, j, forest[i][j], true) *
				calcLeft(forest, i, j, forest[i][j], true) *
				calcRight(forest, i, j, forest[i][j], true) *
				calcUp(forest, i, j, forest[i][j], true)
			if scenicScore > highestScenicScore {
				highestScenicScore = scenicScore
			}
		}
	}

	return highestScenicScore, nil
}

/*
	calcUp, calcDown, calcLeft, calcRight count and return the number of trees shorter than location in each direction

	forest is the forest structure built from the input file.
	i and j are the current coordinates of the location being checked.
	location is the height of the tree at the potential location for the tree house
	firtRec is a bool that checks if the current recursion is the first to prevent the starting point forest[i][j] = location from being true.
*/

func calcUp(forest [][]rune, i, j int, location rune, firstRec bool) int {
	if i == 0 || !firstRec && forest[i][j] >= location {
		return 0
	}
	return 1 + calcUp(forest, i-1, j, location, false)
}

func calcDown(forest [][]rune, i, j int, location rune, firstRec bool) int {
	if i == len(forest)-1 || !firstRec && forest[i][j] >= location {
		return 0
	}
	return 1 + calcDown(forest, i+1, j, location, false)
}

func calcLeft(forest [][]rune, i, j int, location rune, firstRec bool) int {
	if j == 0 || !firstRec && forest[i][j] >= location {
		return 0
	}
	return 1 + calcLeft(forest, i, j-1, location, false)
}

func calcRight(forest [][]rune, i, j int, location rune, firstRec bool) int {
	if j == len(forest)-1 || !firstRec && forest[i][j] >= location {
		return 0
	}
	return 1 + calcRight(forest, i, j+1, location, false)
}

func main() {
	numVisible, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(numVisible)

	numVisible, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(numVisible)
}
