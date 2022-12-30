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

	var monkeys = map[string]string{}

	input, err := os.ReadFile("input.txt")
	if err != nil {
		return -1, err
	}

	for _, s := range strings.Split(strings.TrimSpace(string(input)), "\n") {
		s := strings.Split(s, ": ")
		monkeys[s[0]] = s[1]
	}

	return solve("root", monkeys), nil
}

func partTwo() (int, error) {
	var monkeys = map[string]string{}

	input, err := os.ReadFile("input.txt")
	if err != nil {
		return -1, err
	}

	for _, s := range strings.Split(strings.TrimSpace(string(input)), "\n") {
		s := strings.Split(s, ": ")
		monkeys[s[0]] = s[1]
	}

	monkeys["humn"] = "0"
	s := strings.Fields(monkeys["root"])
	if solve(s[0], monkeys) < solve(s[2], monkeys) {
		s[0], s[2] = s[2], s[0]
	}

	part2, _ := sort.Find(1e16, func(v int) int {
		monkeys["humn"] = strconv.Itoa(v)
		return solve(s[0], monkeys) - solve(s[2], monkeys)
	})

	return part2, nil
}

func solve(expr string, monkeys map[string]string) int {
	if v, err := strconv.Atoi(monkeys[expr]); err == nil {
		return v
	}

	s := strings.Fields(monkeys[expr])
	return map[string]func(int, int) int{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int { return a / b },
	}[s[1]](solve(s[0], monkeys), solve(s[2], monkeys))
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
