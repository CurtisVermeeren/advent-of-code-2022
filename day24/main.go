package main

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
)

type State struct {
	P image.Point
	T int
}

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	vall := map[image.Point]rune{}
	y := 0
	for fileScanner.Scan() {
		for x, r := range fileScanner.Text() {
			vall[image.Point{x, y}] = r
		}
		y++
	}

	var bliz image.Rectangle
	for p := range vall {
		bliz = bliz.Union(image.Rectangle{p, p.Add(image.Point{1, 1})})
	}
	bliz.Min, bliz.Max = bliz.Min.Add(image.Point{1, 1}), bliz.Max.Sub(image.Point{1, 1})

	bfs := func(start image.Point, end image.Point, time int) int {
		delta := map[rune]image.Point{
			'#': {0, 0}, '^': {0, -1}, '>': {1, 0}, 'v': {0, 1}, '<': {-1, 0},
		}

		queue := []State{{start, time}}
		seen := map[State]struct{}{queue[0]: {}}

		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]

		loop:
			for _, d := range delta {
				next := State{cur.P.Add(d), cur.T + 1}
				if next.P == end {
					return next.T
				}

				if _, ok := seen[next]; ok {
					continue
				}
				if r, ok := vall[next.P]; !ok || r == '#' {
					continue
				}

				if next.P.In(bliz) {
					for r, d := range delta {
						if vall[next.P.Sub(d.Mul(next.T)).Mod(bliz)] == r {
							continue loop
						}
					}
				}

				seen[next] = struct{}{}
				queue = append(queue, next)
			}
		}
		return -1
	}

	start, end := bliz.Min.Sub(image.Point{0, 1}), bliz.Max.Sub(image.Point{1, 0})
	fmt.Println(bfs(start, end, 0))
	fmt.Println(bfs(start, end, bfs(end, start, bfs(start, end, 0))))

	return -1, nil
}

func partTwo() (int, error) {
	return -1, nil
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
