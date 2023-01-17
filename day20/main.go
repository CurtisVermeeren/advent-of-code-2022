package main

import (
	"container/ring"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func partOne() (int, error) {
	input, err := os.ReadFile("input.txt")
	answer, err := mix(string(input), 1, 1)
	return answer, err
}

func partTwo() (int, error) {
	input, err := os.ReadFile("input.txt")
	answer, err := mix(string(input), 811589153, 10)
	return answer, err
}

// mixing is the process defined in the problem that is used to decrpted the file input
// A file is mixed by moving each number forward or backward in the file a number of positions equal to the value of the number being moved.
// The list is circular so moving off one end wraps around to the other side
// mix takes an input list of numbers in the encrypted file
// key is the decryption key needed for part two
// times is the number of times the encryped list must be mixed
func mix(input string, key, times int) (int, error) {
	split := strings.Fields(input)

	// Split the input string and create a new go ring and a map of idx
	// The go ring has functions to move items, traverse the ring with Prev and Next and link and unlink items
	r, idx, z := ring.New(len(split)), map[int]*ring.Ring{}, (*ring.Ring)(nil)
	for i, s := range split {
		v, err := strconv.Atoi(s)
		if err != nil {
			return -1, err
		}

		if v == 0 {
			z = r
		}
		// multiply each value by the decryption key
		r.Value, idx[i], r = v*key, r, r.Next()
	}

	// For each element in the index map remove the current ring element
	// Move the ring element to its new position with is mod (len(idx)-1)
	// Add the ring element back to the ring
	// It runs the number of times specified in the problem. once in part one. 10 times in part two
	for i := 0; i < times; i++ {
		for i := 0; i < len(idx); i++ {
			r = idx[i].Prev()
			c := r.Unlink(1)
			r.Move(c.Value.(int) % (len(idx) - 1)).Link(c)
		}
	}

	// sum the values of the 1000th 2000th and 3000th numbers
	return z.Move(1000).Value.(int) + z.Move(2000).Value.(int) + z.Move(3000).Value.(int), nil
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
