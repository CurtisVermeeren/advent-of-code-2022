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

func mix(input string, key, times int) (int, error) {
	split := strings.Fields(input)

	r, idx, z := ring.New(len(split)), map[int]*ring.Ring{}, (*ring.Ring)(nil)
	for i, s := range split {
		v, err := strconv.Atoi(s)
		if err != nil {
			return -1, err
		}

		if v == 0 {
			z = r
		}
		r.Value, idx[i], r = v*key, r, r.Next()
	}

	for i := 0; i < times; i++ {
		for i := 0; i < len(idx); i++ {
			r = idx[i].Prev()
			c := r.Unlink(1)
			r.Move(c.Value.(int) % (len(idx) - 1)).Link(c)
		}
	}

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
