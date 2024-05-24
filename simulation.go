package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// SimulateEcosystemEvolution() takes an Ecosystem pointer initialEcosystem, int totalTimesteps, and string foodRule. The function sequentially simulates initialEcosystem evolving over the course of totalTimesteps generations and saves each Ecosystem into a collection for the output, a slice of Ecosystem pointers called allEcosystems.
func SimulateEcosystemEvolution(initialEcosystem *Ecosystem, totalTimesteps int, foodRule string) []*Ecosystem {
	fmt.Println("SimulateEcosystemEvolution is running")

	// initialize the slice to contain all the Ecosystems
	allEcosystems := make([]*Ecosystem, totalTimesteps+1)
	allEcosystems[0] = initialEcosystem
	// keep track of start of simulation
	var start time.Time = time.Now()

	// sequentially update allEcosystems (serially) and save them
	for i := 1; i <= totalTimesteps; i++ {

		allEcosystems[i] = UpdateEcosystem(allEcosystems[i-1], foodRule, i)

		// print status of simulation
		if (totalTimesteps / 10) != 0 {
			if i%(totalTimesteps/10) == 0 || i == 1 {
				fmt.Println("Simulation is", 100*float64(i)/float64(totalTimesteps), "percent complete. Generation =", i)
				elapsed := time.Since(start)
				log.Printf("This took total %s\n\n", elapsed)
			}
		}
	}

	return allEcosystems
}

func UpdateEcosystem(prevEcosystem *Ecosystem, foodRule string, curGen int) *Ecosystem {

	// initialize the nextEcosystem
	var nextEcosystem *Ecosystem = DeepCopyEcosystem(prevEcosystem)

	// get the numRows and numCols of the ecosystem
	numRows := nextEcosystem.CountRows()
	numCols := nextEcosystem.CountCols()
	arrayOfIndices := MakeIndicesArray(nextEcosystem)

	// get ready to loop through all the possible indices of the ecosystem using a while loop
	totalUpdates := numCols * numRows
	k := 0

	for k < totalUpdates {
		// choose a random index from the arrayOfIndices, each index has a unique OrderedPair representing a Unit location in the Ecosystem
		chosenOrderedPair := ChooseRandomIndices(len(arrayOfIndices))
		i := arrayOfIndices[chosenOrderedPair].row
		j := arrayOfIndices[chosenOrderedPair].col

		// update the array of indices
		arrayOfIndices = UpdateIndices(arrayOfIndices, chosenOrderedPair)

		// get a pointer to the current Unit that we need new values for (nextEcosystem)
		currentUnit := (*nextEcosystem)[i][j]

		// Update the currentUnit based on the nextEcosystem! since we want the system to change as things are disappearing (so each prey/predator is competing to get to their respective food source first)
		// only perform this operation is Unit contains predator
		if (*currentUnit).predator != nil {

			// skip already updated predator.
			if (*currentUnit).predator.lastGenUpdated != curGen {
				currentUnit.predator.UpdatePredator(nextEcosystem, i, j, curGen)
			}

		}
		if (*currentUnit).prey != nil { // only perform this operation if Unit contains prey

			// skip already updated prey.
			if (*currentUnit).prey.lastGenUpdated != curGen {
				UpdatePrey(nextEcosystem, i, j, curGen)
			}

		}

		// we allow predator and prey stacking on top of food
		if !(*currentUnit).food.isPresent { // if we made it here that means there can only be food in the current Unit. skip if the food is already true. note: we don't need to check the lastGenUpdated because food will be false if this is ran.

			// determine whether food appears randomly for the prey. GeneratePreyFoodRandomly() will update both fields of the food, if food is generated. otherwise it will leave it false.
			currentUnit.GeneratePreyFoodProbabilistically(foodRule, i, j, nextEcosystem)
			// currentUnit.food.lastGenUpdated = curGen

		}

		k++
	}

	return nextEcosystem
}

// Input: the number of rows and number of cols to choose from, numRows and numCols
// Output: two integers, randomly choosen between for intervals [0,numRows) and [0,numCols)
func ChooseRandomIndices(numChoices int) int {
	chosenOrderedPair := rand.Intn(numChoices)

	return chosenOrderedPair
}

// Input: an Ecosystem pointer
// Output: a 1D array containing every combination of row and col (all the indices of that Ecosystem) as an OrderedPair
// https://go.dev/play/p/k3_wRjPRuzH
func MakeIndicesArray(someEcosystem *Ecosystem) []OrderedPair {

	// get the number of rows and cols in the ecosystem
	numRows := someEcosystem.CountRows()
	numCols := someEcosystem.CountCols()

	// initialize the arrays to contain all the rows and column indices in the ecosystem
	arrayOfIndices := make([]OrderedPair, numRows*numCols)

	index := 0
	// range over all the arrayOfIndices and set equal to the index
	for i := 0; i < numRows; i++ {
		// assign all the OrderedPair into the array
		for j := 0; j < numCols; j++ {
			arrayOfIndices[index].row = i
			arrayOfIndices[index].col = j
			index += 1
		}
	}

	return arrayOfIndices
}

// Input: indicesAvailable, a 2D array of OrderedPair containing the indices available for the ecosystem to choose. dont confuse the indices of this array, with the indices available!
// Output: none! operates on a pointer to makes updated indicesAvailable
func UpdateIndices(indicesAvailable []OrderedPair, chosenOrderedPair int) []OrderedPair {
	// delete the chosen location from the array of indices available.
	indicesAvailable = append((indicesAvailable)[:chosenOrderedPair], (indicesAvailable)[chosenOrderedPair+1:]...)

	return indicesAvailable
}

// assumes rectangular Ecosystem (rows all the same length)
func (someEcosystem *Ecosystem) CountRows() int {
	return len(*someEcosystem)
}

// assumes rectangular Ecosystem (cols all the same length)
func (someEcosystem *Ecosystem) CountCols() int {
	return len((*someEcosystem)[0])
}

func DeepCopyEcosystem(someEcosystem *Ecosystem) *Ecosystem {
	numCols := someEcosystem.CountCols()
	numRows := someEcosystem.CountRows()

	// make the rows of the copy (outermost dimension)
	var copyEcosystem Ecosystem = make([][]*Unit, numRows)

	// range over all the Units in someEcosystem's (every combo of row and col)
	for i := 0; i < numRows; i++ {
		// make the columns of the copy
		copyEcosystem[i] = make([]*Unit, numCols)
		for j := 0; j < numCols; j++ {
			// initialize the Unit that will go at i, j
			copyEcosystem[i][j] = new(Unit)

			// copy the corresponding fields of the Unit (deep copy)
			copyEcosystem[i][j].food = (*someEcosystem)[i][j].food

			// only attempt to copy if its there
			if (*someEcosystem)[i][j].prey != nil {
				copyEcosystem[i][j].prey = (*someEcosystem)[i][j].prey.DeepCopyOrganism()
			}

			// only attempt to copy if its there
			if (*someEcosystem)[i][j].predator != nil {
				copyEcosystem[i][j].predator = (*someEcosystem)[i][j].predator.DeepCopyOrganism()
			}
		}
	}

	return &copyEcosystem
}

func (somePrey *Prey) DeepCopyOrganism() *Prey {
	var preyCopy Prey
	preyCopy.age = somePrey.age
	preyCopy.energy = somePrey.energy

	// range over the genome and copy all its genes
	var copyGenome [8]Gene
	for i := range copyGenome {
		copyGenome[i] = (*somePrey).genome[i]
	}

	// assign it into preyCopy
	preyCopy.genome = copyGenome

	return &preyCopy // return a pointer to the new copy
}

func (somePred *Predator) DeepCopyOrganism() *Predator {
	var predCopy Predator

	predCopy.age = somePred.age
	predCopy.energy = somePred.energy

	// range over the genome and copy all its genes
	var copyGenome [8]Gene
	for i := range copyGenome {
		copyGenome[i] = (*somePred).genome[i]
	}

	// assign it into predCopy
	predCopy.genome = copyGenome

	return &predCopy // return a pointer to the new copy
}
