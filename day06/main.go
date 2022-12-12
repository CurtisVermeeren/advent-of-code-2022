package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// parameter distinct is the number of distinct runes that must be part of the sequence
func partOne(distinct int) (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)

	fileScanner.Scan()
	characters := []rune(fileScanner.Text())

	// Check sequences along the text 4 runes at a time
	for i, _ := range fileScanner.Text() {

		// map is used to track characters found in the current sequence
		check := make(map[rune]bool)

		// check the 4 character sequence from i
		for x := 0; x < distinct; x++ {

			// Check for rune already found in this sequence
			if !check[characters[i+x]] {
				check[characters[i+x]] = true

				// If 4 runes have been found they are unique
				if x == distinct-1 {
					return i + distinct, nil
				}
			} else {
				// Break if a repeat rune was found in the map
				break
			}
		}
	}

	return -1, nil
}

func main() {
	charactersProcessed, err := partOne(4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(charactersProcessed)

	// reuse partOne with a sequence of 14 distinct runes
	charactersProcessed, err = partOne(14)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(charactersProcessed)

}
