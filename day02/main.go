package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/*
	ABC is opponent, XYZ is you.

	Rock  		A = X
	Paper 		B = Y
	Scissors 	C = Z
*/

/*

	Scoring:
		Rock = 1
		Paper = 2
		Scissors = 3

		Win = 6
		Draw = 3
		Lose = 0

*/

func partOne() (int, error) {

	scoring := map[string]int{"A X": 4, "A Y": 8, "A Z": 3, "B X": 1, "B Y": 5, "B Z": 9, "C X": 7, "C Y": 2, "C Z": 6}

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	totalScore := 0

	for fileScanner.Scan() {
		totalScore += scoring[fileScanner.Text()]
	}

	return totalScore, nil
}

/*
	Opponent plays:
		Rock  		A
		Paper 		B
		Scissors 	C

	Scoring:
		Rock = 1
		Paper = 2
		Scissors = 3

		Win = 6
		Draw = 3
		Lose = 0

	You need:
		X = need to lose
		Y = need to draw
		Z = need to win

*/

func partTwo() (int, error) {

	// Adjust scoring from part 1 to get the desired result based on what the opponent chooses
	scoring := map[string]int{"A X": 3, "A Y": 4, "A Z": 8, "B X": 1, "B Y": 5, "B Z": 9, "C X": 2, "C Y": 6, "C Z": 7}

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	totalScore := 0

	for fileScanner.Scan() {
		totalScore += scoring[fileScanner.Text()]
	}

	return totalScore, nil
}

func main() {
	score, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(score)

	score, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(score)
}
