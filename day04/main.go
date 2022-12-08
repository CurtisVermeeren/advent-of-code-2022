package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	// Track the number of pairs where one range contains the other
	totalContains := 0

	for fileScanner.Scan() {
		var start1, end1, start2, end2 int
		fmt.Sscanf(fileScanner.Text(), "%d-%d,%d-%d", &start1, &end1, &start2, &end2)
		/*
			Starts and ends listed in order
			condition Elf 1 contains elf 2 completely:   s1 s2 e2 e1
			condition Elf 2 contains elf 1 completely: 	 s2 s1 e1 e2
		*/
		if start2 >= start1 && end2 <= end1 || start1 >= start2 && end2 >= end1 {
			totalContains++
		}
	}

	return totalContains, nil
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	// Track the number of pairs where one range partially or fully contains the other
	totalContains := 0

	for fileScanner.Scan() {
		var start1, end1, start2, end2 int
		fmt.Sscanf(fileScanner.Text(), "%d-%d,%d-%d", &start1, &end1, &start2, &end2)
		/*
			If two ranges overlap there is some number N that exists in both ranges

			start1 <= N <= end1 and start2 <= N <= end2

			This can be reduced to start1 <= end2 AND start2 <= end1

			Each range must start before the other ends for the N value to overlap.
		*/
		if start1 <= end2 && start2 <= end1 {
			totalContains++
		}
	}

	return totalContains, nil
}

func main() {
	totalContains, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(totalContains)

	totalContains, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(totalContains)
}
