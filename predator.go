package main

import (
	"math/rand"
)

// UpdatePredator is a Predator method which will take a Predator input and update the position, initiate eating, reproduction, and age accordingly
func (shark *Predator) UpdatePredator(currEco *Ecosystem, i, j, curGen int) {
	// note we have moved the shark this timestep/generation
	shark.lastGenUpdated = curGen
	numRows := currEco.CountRows()
	numCols := currEco.CountCols()

	if shark.Organism.energy <= 0 {
		(*currEco)[i][j].predator = nil

	} else {
		//4. Reproduction
		if shark.CheckAge(ageThresholdPredator) && shark.CheckEnergy(energyThresholdPredator) {

			freeUnits := GetAvailableUnits(currEco, i, j)

			if len(freeUnits) != 0 {
				var babyShark Predator
				deltaX, deltaY := pickUnit(&freeUnits)
				newI := GetIndex(i, deltaX, numRows)
				newJ := GetIndex(j, deltaY, numCols)
				(*currEco)[newI][newJ].predator = &babyShark
				shark.Reproduce(&babyShark)

			}
		}

		//1. Update POSITION AND ENERGY first if energy is allowed
		// This UpdatePosition will scan through all 7 units, give a list of available units, and use GENOME to update
		//	We prioritize the GENOME instead of the fish
		// This function will UpdatePredatorPosition while returning the new index

		deltaRow, deltaCol, newDirection, geneIndex, newR, newC := shark.UseGenomeToMove(currEco, i, j)

		isMoving := deltaRow != 0 || deltaCol != 0
		shark.DecreaseEnergy(geneIndex, isMoving)

		if shark.energy > 0 {
			(*currEco)[newR][newC].predator = shark
			(*currEco)[newR][newC].predator.lastDirection = newDirection
			(*currEco)[i][j].predator = nil // remove the original pointer
		}

		// 2. FEEDING:
		// Check to eat fish or not
		shark.FeedShark(currEco, newR, newC)

		//3. AGE
		shark.UpdateAge() //Just add one

	}
}

// UseGenomeToMovePredator() uses the genome to decide the next location of the organism in a probabilistic manner
func (shark *Predator) UseGenomeToMove(currentEcosystem *Ecosystem, i, j int) (int, int, int, int, int, int) {
	var moveDeltas OrderedPair
	var geneIndex, newDirection, newI, newJ int
	currentPredator := (*currentEcosystem)[i][j].predator
	numTries := 0

	// 20 is the threshold for max number of tries we get to reselect a gene for movement
	// if numberTries >= 20 and isFreeUnitFlag is still false
	// the prey doesn't move
	for !shark.isFreeUnit(currentEcosystem, i, j) && numTries < 20 {
		r := rand.Float64()
		geneIndex = 0
		runningSum := 0.0
		for idx, gene := range shark.genome {
			runningSum += float64(gene)
			if runningSum >= r {
				geneIndex = idx
				break
			}
		}
		newDirection := (shark.lastDirection + geneIndex) % 8
		moveDeltas = deltas[newDirection]
		numRows := currentEcosystem.CountRows()
		numCols := currentEcosystem.CountCols()
		newI = GetIndex(i, moveDeltas.row, numRows)
		newJ = GetIndex(j, moveDeltas.col, numCols)

		// This check if the unit is free or not
		numTries += 1
	}
	// if numTries >= 20 and still haven't find a free unit, we don't move
	if !shark.isFreeUnit(currentEcosystem, i, j) {
		geneIndex = 0
		newDirection = currentPredator.lastDirection
		moveDeltas.row, moveDeltas.col = 0, 0
	}

	//lastDirection will be updated with my new direction
	return moveDeltas.row, moveDeltas.col, newDirection, geneIndex, newI, newJ
}

func (shark *Predator) isFreeUnit(currEco *Ecosystem, i, j int) bool {
	return (*currEco)[i][j].predator == nil
}

func (shark *Predator) FeedShark(currEco *Ecosystem, x, y int) {
	if (*currEco)[x][y].prey != nil {
		(*currEco)[x][y].prey = nil
		shark.IncreaseEngeryAfterMeal() //increase energy after eating a fish

	}

}

func (shark *Predator) IncreaseEngeryAfterMeal() {
	shark.Organism.energy += 1
}

func GetAvailableUnits(currEco *Ecosystem, r, c int) []int {
	var units []int
	var n int
	for i := r - 1; i <= r+1; i++ {
		i_updated := -1

		if i < 0 {
			i_updated = len(*currEco) - 1
		}
		if i == len(*currEco) {
			i_updated = 0
		}

		for j := c - 1; j <= c+1; j++ {
			j_updated := -1
			if j < 0 {
				j_updated = len((*currEco)[0]) - 1
			}
			if j == len((*currEco)[0]) {
				j_updated = 0
			}

			// j_updated didn't change, it should be j still
			if j_updated == -1 {
				j_updated = j
			}

			if i_updated == -1 {
				i_updated = i
			}

			if IsItAvailable((*currEco)[i_updated][j_updated], true) {
				n = GetUnit(r, c, i_updated, j_updated, len(*currEco))
				units = append(units, n)
			}

		}
	}
	return units
}

func (shark *Predator) Reproduce(babyShark *Predator) {
	//Already check age and energy!!!!
	shark.Organism.age = 0
	babyShark.Organism.energy = shark.Organism.energy / 2
	shark.Organism.energy /= 2
	babyShark.Organism.genome = shark.Organism.genome // Check if the array needs to be copied manually.
	UpdateDirection(&shark.Organism, &babyShark.Organism)
	UpdateGenome(&babyShark.Organism)
}

func (shark *Predator) CheckAge(threshold int) bool {
	return shark.Organism.age >= threshold
}

func (shark *Predator) CheckEnergy(threshold int) bool {
	return shark.Organism.energy >= threshold
}

func (shark *Predator) UpdateAge() {
	shark.Organism.age += 1
}

func (shark *Predator) DecreaseEnergy(geneIndex int, isMoving bool) {
	shark.energy -= costOfLivingPredator

	if isMoving {
		shark.energy -= energyCosts[geneIndex]
	}
}

func IsItAvailable(unit *Unit, IsThisAPredator bool) bool {
	//Check if there is any predator
	return unit.predator == nil
}

func GetUnit(r, c, i, j, n int) int {
	var unit int
	rowDelta := r - i
	colDelta := c - j

	//edge case
	if rowDelta < -1 { //first row
		rowDelta = 1
	}
	if rowDelta > 1 { //last row
		rowDelta = -1
	}
	if colDelta < -1 { //first col
		colDelta = 1
	}
	if colDelta > 1 { //last col
		colDelta = -1
	}

	if rowDelta == -1 && colDelta == -1 {
		unit = 0
	} else if rowDelta == -1 && colDelta == 0 {
		unit = 1
	} else if rowDelta == -1 && colDelta == 1 {
		unit = 2
	} else if rowDelta == 0 && colDelta == 1 {
		unit = 3
	} else if rowDelta == 1 && colDelta == 1 {
		unit = 4
	} else if rowDelta == 1 && colDelta == 0 {
		unit = 5
	} else if rowDelta == 1 && colDelta == -1 {
		unit = 6
	} else if rowDelta == 0 && colDelta == -1 {
		unit = 7
	}

	return unit
}
