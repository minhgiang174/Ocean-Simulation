package main

import (
	"math"
	"math/rand"
)

//Set a constant dictionary where keys are the directionIndex and the values are the orderedPair with corresponding deltaX and deltaY
// movementDeltas := map[int]OrderedPair {
// 	0: OrderedPair{-1, -1},
// 	1: OrderedPair{0, -1},
// 	2: OrderedPair{1, -1},
// 	3: OrderedPair{1, 0},
// 	4: OrderedPair{1, 1},
// 	5: OrderedPair{0, 1},
// 	6: OrderedPair{-1, 1},
// 	7: OrderedPair{-1, 0},
// }

// Gene index to energy cost
// energyCosts := map[int]int {
// 	0: 0,
// 	1: -1,
// 	2: -2,
// 	3: -4,
// 	4: -8,
// 	5: -4,
// 	6: -2,
// 	7: -1,
// }

// Input: currentUnit is a pointer to a unit, currentEcosystem is a pointer to the ecosystem, i and j are the indices of the location of the unit we are about to move, curGen is the number of generations of the unit we are about to move.
// Output: none, operates on pointers
func MovePrey(currentEcosystem *Ecosystem, i, j int) {
	currentUnit := (*currentEcosystem)[i][j]
	currentPrey := currentUnit.prey

	deltaX, deltaY, newDirection, geneIndex, newI, newJ := UseGenomeToMovePrey(currentEcosystem, currentPrey, i, j)

	// energy decreases based on how drastic the change in direction is for the movement
	// if at least one of deltaX or deltaY is not equal to 0, we move the prey
	isMoving := deltaX != 0 || deltaY != 0
	currentPrey.DecreaseEnergy(geneIndex, isMoving)

	currentUnit.prey = nil

	// check if energy level > 0
	// if it is not, update direction
	if currentPrey.energy > 0 {

		// when deltaX and deltaY == 0, currentPrey stay at unit [i, j]
		(*currentEcosystem)[newI][newJ].prey = currentPrey

		// comes after moving the prey
		currentPrey.lastDirection = newDirection

	}

	if CheckIfEats((*currentEcosystem)[newI][newJ], currentPrey) {
		currentPrey.FeedOrganism((*currentEcosystem)[newI][newJ])
	}
}

func CheckIfEats(currentUnit *Unit, currentPrey *Prey) bool {
	return currentUnit.food.isPresent && (currentPrey.energy < maxEnergy)
}

func (currentPrey Prey) FeedOrganism(currentUnit *Unit) {
	currentUnit.food.isPresent = false
	currentPrey.energy += energyGainedPerPlankton
}

// cannot move to unit where there's shark (predator)
func UseGenomeToMovePrey(currentEcosystem *Ecosystem, currentPrey *Prey, i, j int) (int, int, int, int, int, int) {
	var moveDeltas OrderedPair
	var geneIndex, newDirection, newI, newJ int
	isFreeUnitFlag := false
	numTries := 0

	// 20 is the threshold for max number of tries we get to reselect a gene for movement
	// if numberTries >= 20 and isFreeUnitFlag is still false
	// the prey doesn't move
	for !isFreeUnitFlag && numTries < 20 {
		r := rand.Float64()
		geneIndex = 0
		runningSum := 0.0
		for idx, gene := range currentPrey.genome {
			runningSum += float64(gene)
			if runningSum >= r {
				geneIndex = idx
				break
			}
		}
		newDirection := (currentPrey.lastDirection + geneIndex) % 8
		moveDeltas = deltas[newDirection]
		numRows := currentEcosystem.CountRows()
		numCols := currentEcosystem.CountCols()
		newI = GetIndex(i, moveDeltas.row, numRows)
		newJ = GetIndex(j, moveDeltas.col, numCols)

		isFreeUnitFlag = isFreeUnit(currentEcosystem, i, j)
		numTries += 1
	}
	// if numTries >= 20 and still haven't find a free unit, we don't move
	if !isFreeUnitFlag {
		geneIndex = 0
		newDirection = currentPrey.lastDirection
		moveDeltas.row, moveDeltas.col = 0, 0
	}

	//lastDirection will be updated with my new direction
	return moveDeltas.row, moveDeltas.col, newDirection, geneIndex, newI, newJ
}

func (currentPrey *Prey) DecreaseEnergy(geneIndex int, isMoving bool) {
	currentPrey.energy -= costOfLivingPrey

	// if prey needs to be moved since either deltaX or deltaY or both are not equal to 0
	// we decrease the energy based on the geneIndex
	if isMoving {
		currentPrey.energy -= energyCosts[geneIndex]
	}

}

// check if unit (i, j) is unoccupied by a predator or another prey
func isFreeUnit(currentEcosystem *Ecosystem, i, j int) bool {
	if (*currentEcosystem)[i][j].prey == nil && (*currentEcosystem)[i][j].predator == nil {
		return true
	} else {
		return false
	}
}

// pass in the row, column indices and the delta for movement
// return new row, column indices within the boundary
// boundary is the numRow and numCol of the ecosystem board
func GetIndex(index, delta, boundary int) int {
	newIndex := index + delta
	if newIndex < 0 {
		newIndex = boundary + newIndex
	} else if newIndex >= boundary {
		newIndex = newIndex % boundary
	}
	return newIndex
}

// UpdateAge takes in a pointer to a prey object, updates it age by incrementing it by one and then returns it.
func UpdateAgePrey(p *Prey) {
	p.Organism.age += 1
}

// UpdateAge takes in a pointer to a prey object, updates it age by incrementing it by one and then returns it.
func UpdateAgePredator(p *Predator) {
	p.Organism.age += 1
}

func ReproducePrey(parent, child *Prey) {
	//This function will only be called if the age and energy and requirements are met. Check these requirements before calling this function.
	parent.Organism.age = 0
	child.Organism.energy = parent.Organism.energy / 2
	parent.Organism.energy /= 2
	child.Organism.genome = parent.Organism.genome // Check if the array needs to be copied manually.
	UpdateDirection(&parent.Organism, &child.Organism)
	UpdateGenome(&child.Organism)
}

// UpdateDirection updates the direction of that the child is moving in based on the parents genome and direction of movement
func UpdateDirection(parent, child *Organism) {
	r := rand.Float64()
	var sum Gene
	index := 0
	for i := range parent.genome {
		if Gene(r) < sum {
			index = i
			break
		} else {
			sum += parent.genome[i]
		}
	}
	child.lastDirection = (parent.lastDirection + index) % 8
}

// UpdateGenome updates the genome of the child based on the last known movement.
func UpdateGenome(currentOrganism *Organism) {
	currentDirection := currentOrganism.lastDirection
	delta := 0.8
	for i := range currentOrganism.genome {
		if i != currentDirection {
			if currentOrganism.genome[i]-Gene(delta)*currentOrganism.genome[currentDirection] > 0 {
				currentOrganism.genome[i] -= Gene(delta) * currentOrganism.genome[currentDirection] / 7.0
			}
		}
	}
	currentOrganism.genome[currentDirection] += Gene(delta) * currentOrganism.genome[currentDirection]
	CheckGenome(currentOrganism.genome)
}

// CheckGenome checks that we have not exceeded 1 by summing the genes for a given input genome
func CheckGenome(currentGenome [8]Gene) bool {
	sum := Gene(0.0)
	for i := range currentGenome {
		sum += currentGenome[i]
	}

	result := true

	// if the sum minus 1 rounded to nearest integer is less than 0, that would mean that the sum is under 0.5
	// if the sum plus 1 rounded to the nearest integer is greater than or equal to 1, than means sum is over 1.5
	if math.Round(float64(sum))-1 < 0 || math.Round(float64(sum))+1 >= 1 {
		result = false
	}

	return result
}
func ReproducePredator(p *Predator) *Predator {
	//This function will only be called if the age and energy and requirements are met. Check these requirements before calling this function.
	var child Predator
	p.Organism.age = 0
	child.Organism.energy = p.Organism.energy / 2
	p.Organism.energy /= 2
	child.Organism.genome = p.Organism.genome // Check if the array needs to be copied manually.
	UpdateDirection(&p.Organism, &child.Organism)
	UpdateGenome(&child.Organism)
	return &child
}

func UpdatePrey(currentEcosystem *Ecosystem, i, j, currGen int) {
	numRows := currentEcosystem.CountRows()
	numCols := currentEcosystem.CountCols()
	currentPrey := (*currentEcosystem)[i][j].prey
	// note we have moved the prey this timestep/generation
	currentPrey.lastGenUpdated = currGen

	if currentPrey.Organism.energy <= 0 {
		(*currentEcosystem)[i][j].prey = nil
		return
	}

	UpdateAgePrey(currentPrey)

	if (*currentEcosystem)[i][j].prey.energy >= energyThresholdPrey && (*currentEcosystem)[i][j].prey.age >= ageThresholdPrey {
		var babyPrey Prey

		freeUnits := GetAvailableUnits(currentEcosystem, i, j)

		if len(freeUnits) != 0 {
			deltaX, deltaY := pickUnit(&freeUnits)
			newI := GetIndex(i, deltaX, numRows)
			newJ := GetIndex(j, deltaY, numCols)
			(*currentEcosystem)[newI][newJ].prey = &babyPrey
			ReproducePrey(currentPrey, &babyPrey)
		}

	}
	MovePrey(currentEcosystem, i, j)

}

func pickUnit(freeUnits *[]int) (r, c int) {
	length := len(*freeUnits)
	random := rand.Intn(length)
	chosenUnit := (*freeUnits)[random]
	return GetIndices(&chosenUnit)
}

func GetIndices(chosenUnit *int) (r, c int) {
	if *chosenUnit == 0 {
		return -1, -1
	}
	if *chosenUnit == 1 {
		return -1, 0
	}
	if *chosenUnit == 2 {
		return -1, 1
	}
	if *chosenUnit == 3 {
		return 0, 1
	}
	if *chosenUnit == 4 {
		return 1, 1
	}
	if *chosenUnit == 5 {
		return 1, 0
	}
	if *chosenUnit == 6 {
		return 1, -1
	}
	if *chosenUnit == 7 {
		return 0, -1
	}
	return 0, 0
}
