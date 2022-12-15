package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Monkey struct {
	items        []int         // The items the monkey holds
	operation    func(int) int // This function defines how the worry level changes as the monkey inspects an item
	testAndThrow func(int) int // This functions tests which monkey the item will be thrown to
}

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// monkeys holds all monkeys from the input file
	var monkeys []Monkey

	// Read the input and parse each monkey's information
	for fileScanner.Scan() {

		var newMonkey Monkey

		// Read the starting items
		fileScanner.Scan()
		// Split the "Starting items: 76, 88, 96, 97, 58, 61, 67" strings by comma without the " Starting items: " text
		// Then parse each worry value and add it to the monkey's items
		// all items are in the order they will be inspected by the monkey
		for _, item := range strings.Split(fileScanner.Text()[len("  Starting items: "):], ", ") {
			worry, err := strconv.Atoi(item)
			if err != nil {
				return -1, err
			}
			newMonkey.items = append(newMonkey.items, worry)
		}

		// Read the Operation:
		fileScanner.Scan()
		var operator rune
		var value int = 0
		// Check if the operation for the monkey inpuit is "old * old" or an "old * value"
		if fileScanner.Text() == "  Operation: new = old * old" {
			operator = '*'
		} else {
			fmt.Sscanf(fileScanner.Text(), " Operation: new = old %c %d", &operator, &value)
		}
		// Create the Operation function for the monkey Operation input
		newMonkey.operation = createOperationOne(operator, value)

		// Read the Test:
		fileScanner.Scan()
		var testValue int
		// Read the test case
		fmt.Sscanf(fileScanner.Text(), "  Test: divisible by %d", &testValue)

		// Read the true expression to find the monkey the item is thrown to if the test case is true
		fileScanner.Scan()
		var monkeyThrownIfTrue int
		fmt.Sscanf(fileScanner.Text(), "    If true: throw to monkey %d", &monkeyThrownIfTrue)

		// Read the false expression to find the monkey the item is thrown to if the test case is false
		fileScanner.Scan()
		var monkeyThrownIfFalse int
		fmt.Sscanf(fileScanner.Text(), "    If false: throw to monkey %d", &monkeyThrownIfFalse)

		// Create a function that runs the test case
		newMonkey.testAndThrow = createTest(testValue, monkeyThrownIfTrue, monkeyThrownIfFalse)

		// Add the monkey to the slice of all monkeys
		monkeys = append(monkeys, newMonkey)

		// Skip the blank line between monkeys in the input file
		fileScanner.Scan()

	}

	// track how many items each monkey has inspected
	inspectedItems := make(map[int]int)

	// Problem asks for answer after 20 rounds of monkey buisness
	for round := 0; round < 20; round++ {
		// iterate over each monkey and their items
		for monkeyId, currentMonkey := range monkeys {
			for _, item := range currentMonkey.items {
				// calculate the new worry value for each item
				newValue := currentMonkey.operation(item)
				// determine who the item is thrown to
				throwTo := currentMonkey.testAndThrow(newValue)
				// append the item to the monkey who receives the item
				monkeys[throwTo].items = append(monkeys[throwTo].items, newValue)
			}

			// Count the number of items the currentMonkey has inspected and thrown
			inspectedItems[monkeyId] += len(monkeys[monkeyId].items)
			// Once a monkey has inspected and thrown all their items they will have none for this round
			monkeys[monkeyId].items = []int{}
		}
	}

	// Get the level of monkey buisness by multiply the number of items inspected of the top two monkeys
	var first, second int
	for _, count := range inspectedItems {
		if count > second {
			second = count
		}
		if second > first {
			first, second = second, first
		}
	}

	return first * second, nil
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// monkeys holds all monkeys from the input file
	var monkeys []Monkey

	// limit large numbers
	// (a mod kn) = a mod n for any integer k
	// To keep numbers small we want to keep only numbers around that are modulo the multiple of all divisors to get the k value.
	// All monkey Test divisors are prime so their least common multiple is just the product of them.
	// When newValue of worry is created we take the mod of the result of the operation.
	// This allows checking the Test for all monkeys while also preventing int from overflowing.
	// The actual worry value doesn't matter to the answer. Only the solution to Test which monkey the item is thrown to must remain the same.
	var limit int = 1

	// Read the input and parse each monkey's information
	for fileScanner.Scan() {

		var newMonkey Monkey

		// Read the starting items
		fileScanner.Scan()
		// Split the "Starting items: 76, 88, 96, 97, 58, 61, 67" strings by comma without the " Starting items: " text
		// Then parse each worry value and add it to the monkey's items
		// all items are in the order they will be inspected by the monkey
		for _, item := range strings.Split(fileScanner.Text()[len("  Starting items: "):], ", ") {
			worry, err := strconv.Atoi(item)
			if err != nil {
				return -1, err
			}
			newMonkey.items = append(newMonkey.items, worry)
		}

		// Read the Operation:
		fileScanner.Scan()
		var operator rune
		var value int = 0
		// Check if the operation for the monkey inpuit is "old * old" or an "old * value"
		if fileScanner.Text() == "  Operation: new = old * old" {
			operator = '*'
		} else {
			fmt.Sscanf(fileScanner.Text(), " Operation: new = old %c %d", &operator, &value)
		}
		// Create the Operation function for the monkey Operation input
		newMonkey.operation = createOperationTwo(operator, value)

		// Read the Test:
		fileScanner.Scan()
		var testValue int
		// Read the test case
		fmt.Sscanf(fileScanner.Text(), "  Test: divisible by %d", &testValue)

		// set the limit for the modulo
		// we get a large limit by multiplying all test values for the monkeys together
		// This is the least common multiple as all monkey Test divisors are prime
		limit *= testValue

		// Read the true expression to find the monkey the item is thrown to if the test case is true
		fileScanner.Scan()
		var monkeyThrownIfTrue int
		fmt.Sscanf(fileScanner.Text(), "    If true: throw to monkey %d", &monkeyThrownIfTrue)

		// Read the false expression to find the monkey the item is thrown to if the test case is false
		fileScanner.Scan()
		var monkeyThrownIfFalse int
		fmt.Sscanf(fileScanner.Text(), "    If false: throw to monkey %d", &monkeyThrownIfFalse)

		// Create a function that runs the test case
		newMonkey.testAndThrow = createTest(testValue, monkeyThrownIfTrue, monkeyThrownIfFalse)

		// Add the monkey to the slice of all monkeys
		monkeys = append(monkeys, newMonkey)

		// Skip the blank line between monkeys in the input file
		fileScanner.Scan()

	}

	// track how many items each monkey has inspected
	inspectedItems := make(map[int]int)

	// Problem asks for answer after 10000 rounds of monkey buisness for part two
	for round := 0; round < 10000; round++ {
		// iterate over each monkey and their items
		for monkeyId, currentMonkey := range monkeys {
			for _, item := range currentMonkey.items {
				// calculate the new worry value for each item
				newValue := currentMonkey.operation(item) % limit
				// determine who the item is thrown to
				throwTo := currentMonkey.testAndThrow(newValue)
				// append the item to the monkey who receives the item
				monkeys[throwTo].items = append(monkeys[throwTo].items, newValue)
			}

			// Count the number of items the currentMonkey has inspected and thrown
			inspectedItems[monkeyId] += len(monkeys[monkeyId].items)
			// Once a monkey has inspected and thrown all their items they will have none for this round
			monkeys[monkeyId].items = []int{}
		}
	}

	// Get the level of monkey buisness by multiply the number of items inspected of the top two monkeys
	var first, second int
	for _, count := range inspectedItems {
		if count > second {
			second = count
		}
		if second > first {
			first, second = second, first
		}
	}

	return first * second, nil
}

// createOperation is used to create a new function that represents a monkey's Operation from the input file for Part One
func createOperationOne(operator rune, value int) func(int) int {
	// func operation takes n which is the worry of an item and returns the new worry of an item
	operation := func(n int) int {
		var valueUsed int = value
		// if valueUsed is zero then it was the old * old operation so value is n
		if valueUsed == 0 {
			valueUsed = n
		}

		if operator == '+' {
			// Perform the worry addition operation then divide the worry by three in "relief"
			return (n + valueUsed) / 3
		}

		// Perform the worry multiplication operation then divide the worry by three in "relief"
		return (n * valueUsed) / 3
	}

	// Return the new function of the Monkey "Operation" from the input
	return operation
}

// createOperation is used to create a new function that represents a monkey's Operation from the input file for Part Two
// worry is no longer divided by 3
func createOperationTwo(operator rune, value int) func(int) int {
	// func operation takes n which is the worry of an item and returns the new worry of an item
	operation := func(n int) int {
		var valueUsed int = value
		// if valueUsed is zero then it was the old * old operation so value is n
		if valueUsed == 0 {
			valueUsed = n
		}

		if operator == '+' {
			// Perform the worry addition operation. Worry is no longer divided by 3
			return (n + valueUsed)
		}

		// Perform the worry multiplication operation. Worry is no longer divided by 3
		return (n * valueUsed)
	}

	// Return the new function of the Monkey "Operation" from the input
	return operation
}

// createTest is used to create a new function that represtents a monkey's Test from the input file
func createTest(testValue, monkeyThrownIfTrue, monkeyThrownIfFalse int) func(int) int {
	// func test takes n which is the worry level of an item and returns the number of the monkey the item is thrown to
	test := func(n int) int {
		if n%testValue == 0 {
			return monkeyThrownIfTrue
		}
		return monkeyThrownIfFalse
	}
	return test
}

func main() {
	monkeyBuisness, err := partOne()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(monkeyBuisness)

	monkeyBuisness, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(monkeyBuisness)
}
