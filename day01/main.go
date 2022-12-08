package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func partOne() (int, error) {
	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	maxCalories := 0
	currentElfCalories := 0

	for fileScanner.Scan() {
		if fileScanner.Text() == "" {
			if currentElfCalories > maxCalories {
				maxCalories = currentElfCalories
			}
			currentElfCalories = 0
		} else {
			cal, err := strconv.Atoi(fileScanner.Text())
			if err != nil {
				return -1, err
			}

			currentElfCalories += cal
		}

	}
	err = file.Close()
	if err != nil {
		return -1, err
	}

	return maxCalories, nil
}

func partTwo() (int, error) {
	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	// maxCalories1 is the most
	maxCalories1 := 0
	maxCalories2 := 0
	maxCalories3 := 0

	currentElfCalories := 0

	for fileScanner.Scan() {
		if fileScanner.Text() == "" {

			// Check if the current elf carries more calories than the top 3.
			// Swap positions as needed
			if currentElfCalories > maxCalories3 {
				maxCalories3 = currentElfCalories
			}
			if maxCalories3 > maxCalories2 {
				maxCalories3, maxCalories2 = maxCalories2, maxCalories3
			}
			if maxCalories2 > maxCalories1 {
				maxCalories2, maxCalories1 = maxCalories1, maxCalories2
			}

			currentElfCalories = 0
		} else {
			cal, err := strconv.Atoi(fileScanner.Text())
			if err != nil {
				return -1, err
			}

			currentElfCalories += cal
		}

	}
	err = file.Close()
	if err != nil {
		return -1, err
	}

	return maxCalories1 + maxCalories2 + maxCalories3, nil
}

func main() {
	calories, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(calories)

	calories2, err := partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(calories2)
}
