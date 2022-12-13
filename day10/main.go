package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func partOne() (int, error) {
	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// registerX starts at 1 and takes two commands
	// "addx V" takes two cycles to complete after which the X register is increased by value V.
	// "noop" takes one cycle to complete and has no other effect.
	registerX := 1
	currentCycle := 0

	var signalStrength int

	for fileScanner.Scan() {
		instructions := strings.Fields(fileScanner.Text())

		currentCycle++

		// Check the signal at 20, 60, 100, 140, 180, 220
		if (currentCycle-20)%40 == 0 {
			signalStrength += currentCycle * registerX
		}

		// If the instreuction is addx it takes an extra cycle
		// increment the currentCycle again and do potential signal strength checks for that cycle.
		if instructions[0] == "addx" {

			// increment again. addx takes two cycles
			currentCycle++

			// Check the signal at 20, 60, 100, 140, 180, 220
			if (currentCycle-20)%40 == 0 {
				signalStrength += currentCycle * registerX
			}

			// Add the value to the registerX
			value, err := strconv.Atoi(instructions[1])
			if err != nil {
				return -1, nil
			}
			registerX += value

		}
	}

	return signalStrength, err
}

func partTwo() error {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(file)

	// registerX starts at 1 and takes two commands
	// "addx V" takes two cycles to complete after which the X register is increased by value V.
	// "noop" takes one cycle to complete and has no other effect.
	registerX := 1
	currentCycle := 0

	for fileScanner.Scan() {
		instructions := strings.Fields(fileScanner.Text())

		// Every 40 "pixels" of the screen is a new line
		if (currentCycle)%40 == 0 {
			fmt.Println()
		}

		// registerX sets the horizontal position of the middle of the sprite "###"
		// The crt screens draws a single pixel each cycle. The screen is 40 x 6. Cycle 1 to 40 draws the top row, Cycles 41 to 80 draw the second row and so on
		// If the sprite is position such that one of its three pixels is the pixel currently being drawn, a lit pixel is drawn "#"
		if registerX-1 == currentCycle%40 || registerX == currentCycle%40 || registerX+1 == currentCycle%40 {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}

		currentCycle++

		// If the instreuction is addx it takes an extra cycle
		// increment the currentCycle again and do drawing necessary depending on registerX
		if instructions[0] == "addx" {

			// Every 40 "pixels" of the screen is a new line
			if (currentCycle)%40 == 0 {
				fmt.Println()
			}

			// Draw again to account for the second addx cycle
			if registerX-1 == currentCycle%40 || registerX == currentCycle%40 || registerX+1 == currentCycle%40 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}

			// increment again. addx takes two cycles
			currentCycle++

			// Add the value to the registerX
			value, err := strconv.Atoi(instructions[1])
			if err != nil {
				return err
			}
			registerX += value

		}
	}

	return err
}

func main() {
	signal, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signal)

	partTwo()
}
