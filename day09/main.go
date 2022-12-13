package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Pair represents a pair of coordinates on a grid of x and y
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

	// Visited tracks the pairs of xy coordinates the tail has reached
	visited := make(map[Pair]bool)

	// head and tail start at 0,0
	head := Pair{0, 0}
	tail := Pair{0, 0}
	visited[tail] = true

	// read head instructions from the input file
	for fileScanner.Scan() {
		instruct := rune(fileScanner.Text()[0])
		numMoved, err := strconv.Atoi(fileScanner.Text()[2:])
		if err != nil {
			return -1, err
		}

		// move the head up, down, left or right one step at a time
		for numMoved > 0 {
			switch instruct {
			case 'U':
				head.y++
			case 'D':
				head.y--
			case 'L':
				head.x--
			case 'R':
				head.x++
			}

			tail = moveTail(tail, head)
			visited[tail] = true
			numMoved--
		}
	}

	return len(visited), nil
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// Visited tracks the pairs of xy coordinates the tail has reached
	visited := make(map[Pair]bool)

	// rope contains all 10 knots index 0 being the head, index 9 being the tail
	rope := make([]Pair, 10)
	visited[rope[9]] = true

	for fileScanner.Scan() {
		instruct := rune(fileScanner.Text()[0])
		numMoved, err := strconv.Atoi(fileScanner.Text()[2:])
		if err != nil {
			return -1, err
		}

		// move the head up, down, left or right one step at a time
		for numMoved > 0 {
			switch instruct {
			case 'U':
				rope[0].y++
			case 'D':
				rope[0].y--
			case 'L':
				rope[0].x--
			case 'R':
				rope[0].x++
			}

			// each knot of rope acts as the head for the previous knot
			// move each knot after the head (index 1 to 9) one at a time according to the previous knots position
			for i := range rope[:len(rope)-1] {
				rope[i+1] = moveTail(rope[i+1], rope[i])
			}
			visited[rope[9]] = true
			numMoved--
		}
	}

	return len(visited), nil
}

// moveTail takes the tail and the just moved head and adjusted the tails position according to the head postion
func moveTail(tail, head Pair) (movedTail Pair) {
	movedTail = tail
	// Check the difference in the head and tail position
	switch (Pair{head.x - tail.x, head.y - tail.y}) {
	case Pair{0, 2}:
		movedTail.y++ // Tail too far below
	case Pair{2, 0}:
		movedTail.x++ // Tail too far left
	case Pair{0, -2}:
		movedTail.y-- // Tail too far above
	case Pair{-2, 0}:
		movedTail.x-- // Tail too far right
	case Pair{-2, 1}, Pair{-1, 2}, Pair{-2, 2}:
		movedTail.y++ // Move tail diagonally NW
		movedTail.x--
	case Pair{1, 2}, Pair{2, 1}, Pair{2, 2}:
		movedTail.y++ // Move tail diagonally NE
		movedTail.x++
	case Pair{2, -1}, Pair{1, -2}, Pair{2, -2}:
		movedTail.x++ // Move tail diagonally SE
		movedTail.y--
	case Pair{-1, -2}, Pair{-2, -1}, Pair{-2, -2}:
		movedTail.y-- // Move tail diagonally SW
		movedTail.x--
	}
	return
}

func main() {
	visited, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(visited)

	visited, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(visited)
}
