package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func partOne() (int, error) {
	// All monkeys are represented in a map using their name as a key and their job as the value
	var monkeys = map[string]string{}
	input, err := os.ReadFile("input.txt")
	if err != nil {
		return -1, err
	}
	for _, s := range strings.Split(strings.TrimSpace(string(input)), "\n") {
		s := strings.Split(s, ": ")
		monkeys[s[0]] = s[1]
	}

	// Find the value the monkey named root will yell
	return solve("root", monkeys), nil
}

func partTwo() (int, error) {
	// All monkeys are represented in a map using their name as a key and their job as the value
	var monkeys = map[string]string{}
	input, err := os.ReadFile("input.txt")
	if err != nil {
		return -1, err
	}
	for _, s := range strings.Split(strings.TrimSpace(string(input)), "\n") {
		s := strings.Split(s, ": ")
		monkeys[s[0]] = s[1]
	}

	// humn is you as defined in the problem rather than a monkey
	// We must determine what number needs to be yelled so that the two numbers that root receives are equal
	monkeys["humn"] = "0"
	s := strings.Fields(monkeys["root"])
	// Ensure that s[0] is larger for calculating sort.Find
	if solve(s[0], monkeys) < solve(s[2], monkeys) {
		s[0], s[2] = s[2], s[0]
	}

	// Find the value of humn such that the two numbers passed to root s[0] and s[2] are equal
	// s[0] and s[2] are swapped before this such that sort.Find will be equal before it is negative
	part2, _ := sort.Find(1e16, func(v int) int {
		monkeys["humn"] = strconv.Itoa(v)
		return solve(s[0], monkeys) - solve(s[2], monkeys)
	})

	return part2, nil
}

// solve is a recursive function to execute the monkey jobs
// expr represents the name of a monkey and is used as a key in the monkeys map
func solve(expr string, monkeys map[string]string) int {

	// If the job of monkey expr is a value return it
	if v, err := strconv.Atoi(monkeys[expr]); err == nil {
		return v
	}

	// Get the operation for the current expression in the recursion
	s := strings.Fields(monkeys[expr])
	operations := map[string]func(int, int) int{"+": func(a, b int) int { return a + b }, "-": func(a, b int) int { return a - b }, "*": func(a, b int) int { return a * b }, "/": func(a, b int) int { return a / b }}
	// Perform the operation recursively using the monkeys listed in the operation
	return operations[s[1]](solve(s[0], monkeys), solve(s[2], monkeys))
}

func main() {
	answer, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(answer)

	answer, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(answer)
}
