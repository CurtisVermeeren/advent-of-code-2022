package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/RyanCarrier/dijkstra"
	"github.com/ernestosuarez/itertools"
)

func partOne(reachability map[string][]string, pressure map[string]int, mappingInt map[string]int, tunnels []string, matrix map[string]map[string]int) (int, error) {
	maxPressure := 0
	maxPressure = solveRecur(matrix, pressure, 0, 0, 0, "AA", tunnels, 30)

	return maxPressure, nil
}

func partTwo(reachability map[string][]string, pressure map[string]int, mappingInt map[string]int, tunnels []string, matrix map[string]map[string]int) (int, error) {

	maxPressure := 0
	mu := sync.Mutex{}

	for i := 1; i < len(tunnels)/2; i++ {
		for v := range itertools.CombinationsStr(tunnels, i) {
			vcopy := []string{}
			vcopy = append(vcopy, v...)
			go func() {
				maxPressure2a := solveRecur(matrix, pressure, 0, 0, 0, "AA", vcopy, 26)
				maxPressure2b := solveRecur(matrix, pressure, 0, 0, 0, "AA", createOpposite(tunnels, vcopy), 26)
				mu.Lock()
				if maxPressure2a+maxPressure2b > maxPressure {
					maxPressure = maxPressure2a + maxPressure2b
				}
				mu.Unlock()
			}()
		}
	}
	return maxPressure, nil

}

func calculateReachabilityMatrix(reachability map[string][]string, mappingInt map[string]int) map[string]map[string]int {
	graph := dijkstra.NewGraph()

	for k := range reachability {
		graph.AddVertex(mappingInt[k])
	}

	for k, v := range reachability {
		for _, l := range v {
			graph.AddArc(mappingInt[k], mappingInt[l], 1)
		}
	}

	matrix := map[string]map[string]int{}
	for k1, v1 := range mappingInt {
		matrix[k1] = map[string]int{}
		for k2, v2 := range mappingInt {
			best, _ := graph.Shortest(v1, v2)
			matrix[k1][k2] = int(best.Distance)
		}
	}

	return matrix
}

func solveRecur(matrix map[string]map[string]int, pressures map[string]int, currentTime int, currentPressure int, currentFlow int, currentTunnel string, remaining []string, limit int) int {

	nScore := currentPressure + (limit-currentTime)*currentFlow

	max := nScore

	for _, v := range remaining {
		distanceAndOpen := matrix[currentTunnel][v] + 1
		if currentTime+distanceAndOpen < limit {
			newTime := currentTime + distanceAndOpen
			newPressure := currentPressure + distanceAndOpen*currentFlow
			newFlow := currentFlow + pressures[v]
			possibleScore := solveRecur(matrix, pressures, newTime, newPressure, newFlow, v, removeFromList(remaining, v), limit)
			if possibleScore > max {
				max = possibleScore
			}
		}
	}

	return max
}

func createOpposite(all []string, paratial []string) []string {
	new := []string{}

loop:
	for _, v := range all {
		for _, w := range paratial {
			if v == w {
				continue loop
			}
		}

		new = append(new, v)
	}
	return new
}

// removeFromList returns a slive of all values of in that are not equal to v
func removeFromList(in []string, v string) []string {
	new := []string{}
	for _, i := range in {
		if i != v {
			new = append(new, i)
		}
	}
	return new
}

func main() {

	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)

	reachability := map[string][]string{}
	pressure := map[string]int{}

	mappingInt := map[string]int{}
	count := 0

	for fileScanner.Scan() {
		parts := strings.Split(fileScanner.Text(), " ")
		valve := parts[1]
		var rate int
		_, err := fmt.Sscanf(parts[4], "rate=%d;", &rate)
		if err != nil {
			log.Fatal(err)
		}

		paths := parts[9:]

		reachability[valve] = []string{}
		for _, v := range paths {
			reachability[valve] = append(reachability[valve], strings.TrimSuffix(v, ","))
		}
		pressure[valve] = rate

		mappingInt[valve] = count
		count++
	}

	// Build tunnels of only valves with a flow rate more than 0
	tunnels := []string{}
	for k, v := range pressure {
		if v != 0 {
			tunnels = append(tunnels, k)
		}
	}

	// Use Dijkstra to build a reachability matrix
	matrix := calculateReachabilityMatrix(reachability, mappingInt)

	pressureReleased, err := partOne(reachability, pressure, mappingInt, tunnels, matrix)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pressureReleased)

	pressureReleased, err = partTwo(reachability, pressure, mappingInt, tunnels, matrix)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pressureReleased)
}
