package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	index, totalSum := 1, 0

	for fileScanner.Scan() {

		// Use Go's builtin JSON unmarshal to read each line into any
		var lineOne, lineTwo any

		// Read the first line
		json.Unmarshal([]byte(fileScanner.Text()), &lineOne)

		// Read the second line
		fileScanner.Scan()
		json.Unmarshal(fileScanner.Bytes(), &lineTwo)

		// skip the blank line
		fileScanner.Scan()

		if compare(lineOne, lineTwo) <= 0 {
			totalSum += index
		}

		index++
	}

	return totalSum, nil
}

func partTwo() (int, error) {
	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// packets holds each line of input packet
	packets := []any{}

	for fileScanner.Scan() {

		// Use Go's builtin JSON unmarshal to read each line into any
		var lineOne, lineTwo any

		// Read the first line
		json.Unmarshal([]byte(fileScanner.Text()), &lineOne)

		// Read the second line
		fileScanner.Scan()
		json.Unmarshal(fileScanner.Bytes(), &lineTwo)

		// skip the blank line
		fileScanner.Scan()

		// Add both input lines to the packets
		packets = append(packets, lineOne, lineTwo)
	}

	// add the divider packets with floats to match jsonUnmarshal
	packets = append(packets, []any{[]any{2.0}}, []any{[]any{6.0}})

	// sort the slice of packets by comparing each index
	sort.Slice(packets, func(i, j int) bool {
		return compare(packets[i], packets[j]) < 0
	})

	// find the decoder packets
	decoderKey := 1
	for i, pkt := range packets {
		// When the packets are found multiply their index by the other
		if fmt.Sprint(pkt) == "[[2]]" || fmt.Sprint(pkt) == "[[6]]" {
			decoderKey *= i + 1
		}
	}

	return decoderKey, nil
}

// compare takes two lines of input
// If the int returned is 0 or lower then the values are in the correct order
// If the int returned is over 0 then the values are in the wrong order
// If the int returned is 0 then the values are the same
func compare(lineOne, lineTwo any) int {
	// Use a type assertion to check if the current section of lineOne and lineTwo are slices
	sliceOne, oneOk := lineOne.([]any)
	sliceTwo, twoOk := lineTwo.([]any)

	switch {
	// If neither are slices because the type assertion fails then they are both values
	// jsonUnmarshal from the input reading section uses float64 for JSON numbers
	case !oneOk && !twoOk:
		// If the first value is larger this will return a positive number meaning the values are NOT in the right order
		// The lower value should be on the left (lineOne)
		return int(lineOne.(float64) - lineTwo.(float64))
	// If just one type assertion failes then one line is a list of values and the other is a value
	// Create a list from the value to compare against the other
	case !oneOk:
		sliceOne = []any{lineOne}
	case !twoOk:
		sliceTwo = []any{lineTwo}
	}

	// While there are values to compare in both lists
	for i := 0; i < len(sliceOne) && i < len(sliceTwo); i++ {
		// Compare each potential list inside of the currentlists recursively
		// If 0 is returned then the values are the same so move to the next
		// If there was a difference in values then return it
		if c := compare(sliceOne[i], sliceTwo[i]); c != 0 {
			return c
		}
	}

	// sliceOne must have less values than sliceTwo for the lists to be in the right order. This would create a negative values
	// if sliceOne is longer it will not be in the correct order and return a positive
	return len(sliceOne) - len(sliceTwo)
}

func main() {
	orderedCorrect, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(orderedCorrect)

	orderedCorrect, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(orderedCorrect)
}
