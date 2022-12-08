package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func partOne() (string, error) {

	stacks := make([][]rune, 9)
	for i := range stacks {
		stacks[i] = make([]rune, 0)
	}

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return "", err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	fileScanner.Scan()
	/*
		Read until the numbered stack labels.

		For each line read 1 rune at a time. If it is not a space or container bracket at the container letter to the coressponding stack

		Divide index by 4 as container labels appear every 4 runes.

		Finally append the container to the start of the slice as we work from top down.
	*/
	for fileScanner.Text() != " 1   2   3   4   5   6   7   8   9 " {
		for i, r := range fileScanner.Text() {
			if r != ' ' && r != '[' && r != ']' {
				stacks[i/4] = append([]rune{r}, stacks[i/4]...)
			}
		}
		fileScanner.Scan()
	}

	// Read blank line after number labels
	fileScanner.Scan()

	for fileScanner.Scan() {
		var numberMoved, from, to int
		fmt.Sscanf(fileScanner.Text(), "move %d from %d to %d", &numberMoved, &from, &to)

		// change from number labels to index values
		from--
		to--

		// move tracks the number of items moved so far
		for move := 0; move < numberMoved; move++ {
			// Get the container to be moved from end of "from" stack
			itemMoved := stacks[from][len(stacks[from])-1]
			// Remove the container from the end of "from" stack
			stacks[from] = stacks[from][:len(stacks[from])-1]
			// Append the container to the end of "to" stack
			stacks[to] = append(stacks[to], itemMoved)
		}
	}

	// Get the top container from each stack to build the solution
	solution := ""
	for _, v := range stacks {
		solution += string(v[len(v)-1])
		fmt.Println(v[len(v)-1])
	}

	return solution, nil

}

func partTwo() (string, error) {
	stacks := make([][]rune, 9)
	for i := range stacks {
		stacks[i] = make([]rune, 0)
	}

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return "", err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	fileScanner.Scan()
	/*
		Read until the numbered stack labels.

		For each line read 1 rune at a time. If it is not a space or container bracket at the container letter to the coressponding stack

		Divide index by 4 as container labels appear every 4 runes.

		Finally append the container to the start of the slice as we work from top down.
	*/
	for fileScanner.Text() != " 1   2   3   4   5   6   7   8   9 " {
		for i, r := range fileScanner.Text() {
			if r != ' ' && r != '[' && r != ']' {
				stacks[i/4] = append([]rune{r}, stacks[i/4]...)
			}
		}
		fileScanner.Scan()
	}

	// Read blank line after number labels
	fileScanner.Scan()

	for fileScanner.Scan() {
		var numberMoved, from, to int
		fmt.Sscanf(fileScanner.Text(), "move %d from %d to %d", &numberMoved, &from, &to)

		// change from number labels to index values
		from--
		to--

		// Select the "numberMoved" containers from the slice in order
		itemsMoved := stacks[from][len(stacks[from])-numberMoved:]
		// Remove the containers from the end of "from" stack
		stacks[from] = stacks[from][:len(stacks[from])-numberMoved]
		// Append the containers to the end of "to" stack
		stacks[to] = append(stacks[to], itemsMoved...)

	}

	// Get the top container from each stack to build the solution
	solution := ""
	for _, v := range stacks {
		solution += string(v[len(v)-1])
		fmt.Println(v[len(v)-1])
	}

	return solution, nil
}

func main() {
	solution, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(solution)

	solution, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(solution)
}
