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
		// Generate all combinations of tunnel nodes of length i
		for v := range itertools.CombinationsStr(tunnels, i) {
			vcopy := []string{}
			vcopy = append(vcopy, v...)

			// For each combination of tunnels run two recursions. One on all combinations of tunnels in V and another on all combinations of tunnels not in v
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

// calculateReachabilityMatrix returns a reachability matrix for the graph represented by reachability.
// reachability maps graph nodes as a key to a slice of strings representing which nodes the key node can reach.
// mappingInt represents a mode of nodes to weight values.
// Returns a map of graph nodes which maps to a map of graph nodes reachable by the first node and an integer representing the shortest distance between the two nodes.
func calculateReachabilityMatrix(reachability map[string][]string, mappingInt map[string]int) map[string]map[string]int {

	// Build a graph of all nodes using the dijkstra package
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
		// Find the shortest distance between each set of nodes and store them in the reachability matrix.
		for k2, v2 := range mappingInt {
			best, _ := graph.Shortest(v1, v2)
			matrix[k1][k2] = int(best.Distance)
		}
	}

	return matrix
}

// matrix is the reachabilty matrix of all nodes
// pressure is a map of nodes to their pressure values
// currentTime is an integer representing the amount of time passed
// currentPressure is an integer representing the current pressure being released
// currentFlow is the current flow rate of all visited nodes added together
// currentTunnel is the current node being visited in the recursion
// remaining is a slice of all remaining nodes that haven't been visited
// limit represents the time limit enforced by the challenge
func solveRecur(matrix map[string]map[string]int, pressures map[string]int, currentTime int, currentPressure int, currentFlow int, currentTunnel string, remaining []string, limit int) int {

	nScore := currentPressure + (limit-currentTime)*currentFlow
	max := nScore

	for _, v := range remaining {
		// distance and open represents the time passed to move from currentTunnel to v and adding 1 for the time it takes to open the valve
		distanceAndOpen := matrix[currentTunnel][v] + 1
		newTime := currentTime + distanceAndOpen
		// check that this value is reached and opened within the time limit
		if newTime < limit {
			// calculate the new pressure released for visiting this node. Take the current pressure and release the currentFlow rate for each minute is takes to reach node v
			newPressure := currentPressure + distanceAndOpen*currentFlow
			// newFlow is the sum of the current flow and the added from visitng node v
			newFlow := currentFlow + pressures[v]
			// recursively continue from this node and attempt to release more pressure
			possibleScore := solveRecur(matrix, pressures, newTime, newPressure, newFlow, v, removeFromList(remaining, v), limit)
			// Find the largest pressure that can be released
			if possibleScore > max {
				max = possibleScore
			}
		}
	}

	return max
}

// createOpposite returns a slice which contains all the alaments from slice "all" that are not prestent in "partial"
func createOpposite(all []string, partial []string) []string {
	new := []string{}

loop:
	for _, v := range all {
		for _, w := range partial {
			if v == w {
				continue loop
			}
		}

		new = append(new, v)
	}
	return new
}

// removeFromList returns a slice of all values of slice in that are not equal to v
// this removes v from the returned slice
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

	// Read the input file
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)

	// Create reachability, pressure, and weight maps
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
	// These are the only destination nodes worth visiting to decrease pressure
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
