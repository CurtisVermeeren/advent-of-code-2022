package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"
)

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)

	// The part 1 solution is the total item priority
	var sumPriority int

	// Scan each "rucksack" line by line
	for fileScanner.Scan() {
		// items will set found items to trues
		items := make(map[rune]bool)

		// Check all items in the left half compartment of the rucksack
		for _, leftItem := range fileScanner.Text()[:len(fileScanner.Text())/2] {
			items[leftItem] = true
		}

		// Check for the "exactly one item" that was mistakenly packed by an elf
		for _, rightItem := range fileScanner.Text()[len(fileScanner.Text())/2:] {
			/*
				Lowercase item types a through z have priorities 1 through 26
				Uppercase item types A through Z have priorities 27 through 52
			*/
			if items[rightItem] {
				sumPriority += int(unicode.ToLower(rightItem) - 96)
				if unicode.IsUpper(rightItem) {
					sumPriority += 26
				}
				break
			}
		}
	}

	return sumPriority, nil
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)

	// The part 2 solution is the total item priority
	var sumPriority int

	// Scan each "rucksack" line by line
	for fileScanner.Scan() {
		// Get 3 lines at a time and map items that are found in each rucksack
		rucksack1 := checkItems(fileScanner.Text())
		fileScanner.Scan()
		rucksack2 := checkItems(fileScanner.Text())
		fileScanner.Scan()
		rucksack3 := checkItems(fileScanner.Text())

		// Iterate through all items of rucksack1 and check if that item exists in rucksack 2 and 3 to find the matching badges
		for itemFirst := range rucksack1 {
			if rucksack2[itemFirst] && rucksack3[itemFirst] {
				sumPriority += int(unicode.ToLower(itemFirst) - 96)
				if unicode.IsUpper(itemFirst) {
					sumPriority += 26
				}
				break
			}
		}
	}

	return sumPriority, nil
}

// checkItems takes a string of items and maps if an item was found
func checkItems(items string) (rucksack map[rune]bool) {
	rucksack = make(map[rune]bool)
	for _, item := range items {
		rucksack[item] = true
	}
	return
}

func main() {
	priority, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(priority)

	priority, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(priority)

}
