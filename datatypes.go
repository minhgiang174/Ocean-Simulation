package main

// 2D array of Unit objects
type Ecosystem [][]*Unit

// OrderedPair
type OrderedPair struct {
	row, col int
}

// Each unit has food, predator, and/or prey. Note: food is static, and will never move so it's not a pointer, but predator and prey move between Unit objects in the Ecosystem, so they are pointers.
type Unit struct {
	food     Food
	predator *Predator
	prey     *Prey
}

type Food struct {
	isPresent bool
	// lastGenUpdated int // if newly made, gets set to current generation. if eaten, gets set to -1. if this is not the current generation, then the Prey can eat it (because it wasn't made during the current generation).
}

type Organism struct {
	// we don't need location OrderedPair because we are using an [][]Unit
	energy         int
	age            int
	genome         [8]Gene
	lastGenUpdated int // gets updated to current generation after the organism has moved (so it doesn't move twice when updating for the next generation)
	lastDirection  int // a number between 0 and 7, corresponding to which gene was chosen for the last movement
}

type Gene float64 // with range 0 to 1. all the genes of a genome add up to 1

type Prey struct {
	Organism
}

type Predator struct {
	Organism
}
