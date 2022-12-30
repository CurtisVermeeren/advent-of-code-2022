package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type Pair struct {
	x, y int
}

func partOne() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// line tracks all positions x = key where a beacon cannot exist
	line := make(map[int]bool)

	for fileScanner.Scan() {
		// Read each sensor and its closest beacon
		var sensor, beacon Pair
		fmt.Sscanf(fileScanner.Text(), "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sensor.x, &sensor.y, &beacon.x, &beacon.y)

		// Calculate the manhattan distance between the sensor and beacon
		distanceFromBeacon := int(math.Abs(float64(sensor.x)-float64(beacon.x))) + int(math.Abs(float64(sensor.y)-float64(beacon.y)))

		// Calculate how far the sensor is from the puzzle line 2000000.  | Y - Sy|
		distanceFromLine := int(math.Abs(float64(2000000 - sensor.y)))

		// i represents all x values alone the horizontal line at Y = 2000000 where the diamond around sensor intersects
		for i := 0; i <= distanceFromBeacon-distanceFromLine; i++ {
			line[sensor.x+i] = true
			line[sensor.x-i] = true
		}

		// remove any beacons on the line from the count
		if beacon.y == 2000000 {
			delete(line, beacon.x)
		}
	}

	return len(line), nil
}

func partTwo() (int, error) {

	file, err := os.Open("./input.txt")
	defer file.Close()
	if err != nil {
		return -1, err
	}
	fileScanner := bufio.NewScanner(file)

	// sensors hold coordinates for all sensors are the key and the manhattan distance as the value
	sensors := map[Pair]int{}

	for fileScanner.Scan() {
		// Read each sensor and its closest beacon
		var sensor, beacon Pair
		fmt.Sscanf(fileScanner.Text(), "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sensor.x, &sensor.y, &beacon.x, &beacon.y)

		// Calculate the manhattan distance between the sensor and beacon
		sensors[sensor] = int(math.Abs(float64(sensor.x)-float64(beacon.x))) + int(math.Abs(float64(sensor.y)-float64(beacon.y)))
	}

	// The distress beacon must have x and y coordinates that are between 0 and 4000000
	for y := 0; y <= 4000000; y++ {
	loop:
		for x := 0; x <= 4000000; x++ {
			// Check the current x,y position against each found sensor
			for sensor, distance := range sensors {
				// dx and yx represents the difference from the sensor x y and the current point x y
				// If the distance from sensor to the current point is less than the distance from the sensor to the beacon then the current point is within the area covered by the the sensors diamond. No closer beacon is here
				if dx, dy := sensor.x-x, sensor.y-y; math.Abs(float64(dx))+math.Abs(float64(dy)) <= float64(distance) {
					// Increment x outside the current sensors diamond
					x += int(float64(distance)-math.Abs(float64(dy))) + dx
					continue loop
				}
			}
			return x*4000000 + y, nil
		}
	}

	return -1, nil
}

func main() {
	noBeacon, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(noBeacon)

	noBeacon, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(noBeacon)
}
