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

var priority []composite = []composite{
	geode,
	obsidian,
	clay,
	ore,
}

type state struct {
	building composite
	minute   int

	collected map[composite]int
	bots      map[composite]int
}

// blueprint represents a blueprint line from the input file
type blueprint struct {
	id   int               // id of the blueprint
	bots map[composite]bot // holds a map of bots based on the material they collect
}

type bot struct {
	costs map[composite]int // map the material name to the cost needed to build this bot
}

type move struct {
	bot composite
}

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

		bp.bots[ore] = bot{costs: map[composite]int{ore: orebotcost}}
		bp.bots[clay] = bot{costs: map[composite]int{ore: claybotcost}}
		bp.bots[obsidian] = bot{costs: map[composite]int{ore: obsorecost, clay: obsclaycost}}
		bp.bots[geode] = bot{costs: map[composite]int{ore: geodeorecost, obsidian: geodeobscost}}

		list = append(list, bp)

	}

	return list, nil
}

func getQualitySum(blueprints []blueprint, minutes int) int {
	sum := 0
	for _, bp := range blueprints {
		res := simulate(bp, minutes)
		sum += (res * bp.id)
	}

	return sum
}

func getGeodeProduct(blueprints []blueprint, minutes int) int {
	prod := 1
	for _, bp := range blueprints {
		res := simulate(bp, minutes)
		prod *= res
	}
	return prod
}

func simulate(bp blueprint, minutes int) int {
	queue := NewQueue[state, int]()

	queue.Enqueue(state{
		minute:    0,
		building:  nothing,
		bots:      map[composite]int{ore: 1, clay: 0, obsidian: 0, geode: 0},
		collected: map[composite]int{ore: 0, clay: 0, obsidian: 0, geode: 0},
	})

	visited := make(map[cachekey]bool)
	maxc := 0

	for queue.Any() {
		current := queue.Dequeue()
		if current.minute > minutes {
			continue
		}

		if current.collected[geode] > maxc {
			maxc = current.collected[geode]
		}

		key := cachekey{
			bc:       current.bots[clay],
			bo:       current.bots[ore],
			bg:       current.bots[geode],
			bob:      current.bots[obsidian],
			min:      current.minute,
			building: current.building,
		}

		if _, ok := visited[key]; ok {
			continue
		}

		visited[key] = true

		current = collect(current)
		if current.building != nothing {
			current.bots[current.building] = current.bots[current.building] + 1
			current.building = nothing
		}

		mvs := getMoves(current, bp)
		if len(mvs) == 0 {
			current.minute++
			queue.Enqueue(current)
			continue
		}

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

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	nm := make(map[K]V)
	for k, v := range m {
		nm[k] = v
	}
	return nm
}

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
