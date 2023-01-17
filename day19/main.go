package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Create a type defition for composite to string
type composite string

// Use constants for each of the building materials
const (
	nothing  composite = "nothing"
	ore      composite = "ore"
	clay     composite = "clay"
	obsidian composite = "obsidian"
	geode    composite = "geode"
)

// A slice containing all bots in order of prioirty
// A geode bot is the most important while an ore bot is the least
var priority []composite = []composite{
	geode,
	obsidian,
	clay,
	ore,
}

// state holds the current "gamestate"
// It tracks what material is currently being built, the minutes that have passed, how many of each material has been collected and the types of robots availabe.
type state struct {
	building composite
	minute   int

	collected map[composite]int
	bots      map[composite]int
}

// blueprint represents a blueprint line from the input file
// Each blueprint has a unique id to identify it
// bots holds the cost to build each type of robot in this blueprint
type blueprint struct {
	id   int
	bots map[composite]bot
}

// bot contains a map of materials to integers that represents the number of each material needed to build this bot.
type bot struct {
	costs map[composite]int
}

type move struct {
	bot composite
}

// cachekey is used as a unique key for state
// bc, bo, bg, bob represents the number of robots of each type
// min is the number of minutes that have passed
// building is the bot type currently being built in the state
type cachekey struct {
	bc, bo, bg, bob int
	min             int
	building        composite
}

func partOne() (int, error) {

	blueprints, err := parseInput()
	if err != nil {
		return -1, err
	}

	return getQualitySum(blueprints, 24), nil
}

func partTwo() (int, error) {

	blueprints, err := parseInput()
	if err != nil {
		return -1, err
	}

	return getGeodeProduct(blueprints[:3], 32), nil
}

// parseInput uses the problem input file to create a slice of all blueprints in it.
// Scan each line of the input file and create blueprint objects be parsing the cost of each robot from the text.
func parseInput() ([]blueprint, error) {
	list := []blueprint{}

	file, err := os.Open("./input.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {

		var id int
		var orebotcost, claybotcost, obsorecost, obsclaycost, geodeorecost, geodeobscost int

		fmt.Sscanf(fileScanner.Text(), "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.", &id, &orebotcost, &claybotcost, &obsorecost, &obsclaycost, &geodeorecost, &geodeobscost)

		bp := blueprint{id: id}
		bp.bots = make(map[composite]bot)

		// For each robot type set the cost needed to built that type of robot as defined in the input blueprint
		bp.bots[ore] = bot{costs: map[composite]int{ore: orebotcost}}
		bp.bots[clay] = bot{costs: map[composite]int{ore: claybotcost}}
		bp.bots[obsidian] = bot{costs: map[composite]int{ore: obsorecost, clay: obsclaycost}}
		bp.bots[geode] = bot{costs: map[composite]int{ore: geodeorecost, obsidian: geodeobscost}}

		list = append(list, bp)

	}
	return list, nil
}

// getQualitySum returns the quality level of all blueprints added together
// Quality level of a blueprint is determined by multiplying the blueprint ID with the largest number of geodes that can be opened in X minutes using that blueprint.
// minutes is the number of minutes the simulation of building runs. In part one this is set as 24.
func getQualitySum(blueprints []blueprint, minutes int) int {
	sum := 0
	for _, bp := range blueprints {
		res := simulate(bp, minutes)
		sum += (res * bp.id)
	}

	return sum
}

// getGeodeProduct calculates the number of geodes for each blueprint in blueprint and multiplies them together.
// For part two only a slice of the first 3 blueprints is used and the time limit for simulation is 32 minutes
func getGeodeProduct(blueprints []blueprint, minutes int) int {
	prod := 1
	for _, bp := range blueprints {
		res := simulate(bp, minutes)
		prod *= res
	}
	return prod
}

// simulate uses a breadth-first search to explore different system states to find the maximum number of geodes that can be collected.
func simulate(bp blueprint, minutes int) int {
	queue := NewQueue[state, int]()

	// Each simulation begins with 1 ore collecting bot as defined in the problem
	queue.Enqueue(state{
		minute:    0,
		building:  nothing,
		bots:      map[composite]int{ore: 1, clay: 0, obsidian: 0, geode: 0},
		collected: map[composite]int{ore: 0, clay: 0, obsidian: 0, geode: 0},
	})

	// visited is a map of all visited states
	// Each state gets a cachekey which is a unique identifier for the current state
	visited := make(map[cachekey]bool)
	// maximum number of geodes collected
	maxc := 0

	// Continue testing the next state while one exists in the queue
	for queue.Any() {
		current := queue.Dequeue()
		// Check that the number of minutes passed is within the limit of the problem
		if current.minute > minutes {
			continue
		}

		// Update that max collected geodes if this state has more
		if current.collected[geode] > maxc {
			maxc = current.collected[geode]
		}

		// Create a unique identifier for the current state based on the number of bots, minutes passed, and the current bot being built
		key := cachekey{
			bc:       current.bots[clay],
			bo:       current.bots[ore],
			bg:       current.bots[geode],
			bob:      current.bots[obsidian],
			min:      current.minute,
			building: current.building,
		}

		// Skip states that have alread been visited
		// If the state is the same it will have the same result
		if _, ok := visited[key]; ok {
			continue
		}
		visited[key] = true

		// collect all materials for the current state
		// If a bot was being built increment the number of that bot type
		current = collect(current)
		if current.building != nothing {
			current.bots[current.building] = current.bots[current.building] + 1
			current.building = nothing
		}

		// check for next moves that are possible with the new materials collected
		mvs := getMoves(current, bp)

		// If no bots can be built add the current updated state to the queue and continue to the next state while increment the time
		if len(mvs) == 0 {
			current.minute++
			queue.Enqueue(current)
			continue
		}

		// If there are bots that can be built iterate over each option as a next possible state
		// Update what is being built, the number of materials collected, and the time and add each new state to the queue
		for _, mv := range mvs {
			st := state{building: current.building, minute: current.minute + 1, collected: copyMap(current.collected), bots: copyMap(current.bots)}
			if mv.bot != nothing {
				st.building = mv.bot
				for c, cost := range bp.bots[st.building].costs {
					st.collected[c] = st.collected[c] - cost
				}
			}
			queue.Enqueue(st)
		}
	}

	return maxc
}

// copyMap is used to create a copy of map m
func copyMap[K comparable, V any](m map[K]V) map[K]V {
	nm := make(map[K]V)
	for k, v := range m {
		nm[k] = v
	}
	return nm
}

// getMoves returns a slice of all possible moves for the current state and blueprint
// Each bot type is checked in priority order from best bot to worst
// For each prioirty the materials needed to build that bot are checked against the materials available.
// If the bot is able to be built then it is a possible move for this state
func getMoves(st state, bp blueprint) []move {
	mvs := []move{}

	for _, c := range priority {
		bot := bp.bots[c]
		build := true
		for needed, cost := range bot.costs {
			if st.collected[needed] < cost {
				build = false
			}
		}
		if build {
			mvs = append(mvs, move{bot: c})
		}
	}
	// add a wait move
	mvs = append(mvs, move{bot: nothing})

	return mvs
}

// collect iterates over all material types and increases the amount of material based on the current robots available.
// The returned state has updated materials after collection
func collect(st state) state {
	for _, c := range priority {
		st.collected[c] = st.collected[c] + st.bots[c]
	}
	return st
}

func main() {
	qualityLevel, err := partOne()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(qualityLevel)

	qualityLevel, err = partTwo()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(qualityLevel)
}
