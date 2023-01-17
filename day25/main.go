package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	// Read the file and split the contents into an array of strings.
	// Iterate over each line and initialize it as a value of n to 0
	input, _ := os.ReadFile("input.txt")
	// For part 1 find the sum of all fuel requirements (each line)
	sum := 0
	for _, s := range strings.Fields(string(input)) {
		n := 0
		// Map each rune to an integer value using the
		// SNAFU uses a power of 5 rather than 10
		for _, r := range s {
			n = 5*n + map[rune]int{'=': -2, '-': -1, '0': 0, '1': 1, '2': 2}[r]
		}
		sum += n
	}

	// Convert the decimal sum back into snafu to
	snafu := ""
	for sum > 0 {
		snafu = string("=-012"[(sum+2)%5]) + snafu
		sum = (sum + 2) / 5
	}
	fmt.Println(snafu)
}
