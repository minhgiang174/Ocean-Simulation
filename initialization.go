package main

import (
	"math/rand"
)

// InitializePreyAndPredator
// Randomly generate numPrey and numPred predators in the initialEcosystem.
// Functions written by Akshat
func InitializePreyAndPredator(numRows, numCols, numPrey, numPred int, newEco *Ecosystem) {
	// Akshat wrote these: Randomly initialize the prey and predators
	count_Prey := 0
	count_Pred := 0

	for count_Pred < numPred {
		i := rand.Intn(numRows)
		j := rand.Intn(numCols)
		if (*newEco)[i][j].predator == nil {
			(*newEco)[i][j].predator = CreatePredator()
			count_Pred += 1
		}

	}

	for count_Prey < numPrey {
		i := rand.Intn(numRows)
		j := rand.Intn(numCols)
		if (*newEco)[i][j].prey == nil && (*newEco)[i][j].predator == nil {
			(*newEco)[i][j].prey = CreatePrey()
			count_Prey += 1
		}
	}
}

// CreatePrey initializes the Prey object
func CreatePrey() *Prey {
	var newPrey Prey
	newPrey.Organism.age = 0
	newPrey.Organism.energy = 50
	newPrey.Organism.age = 0
	newPrey.Organism.genome = CreateGenome()
	newPrey.Organism.lastGenUpdated = 0
	newPrey.Organism.lastDirection = 0
	return &newPrey
}

// CreatePrey initializes the Predator object
func CreatePredator() *Predator {
	var newPredator Predator
	newPredator.Organism.age = 0
	newPredator.Organism.energy = 50
	newPredator.Organism.age = 0
	newPredator.Organism.genome = CreateGenome()
	newPredator.Organism.lastGenUpdated = 0
	newPredator.Organism.lastDirection = 0
	return &newPredator

}

// CreateGenome creates the first version of the genome and returns it.
func CreateGenome() [8]Gene {
	var newGenome [8]Gene
	for i := range newGenome {
		newGenome[i] = Gene(0.125)
	}
	return newGenome
}

func InitializeEcosystem(numRows, numCols, numPrey, numPred int) Ecosystem {

	// initialize newEco, which has numRows rows. the outer dimension
	newEco := make(Ecosystem, numRows)
	for i := 0; i < numRows; i++ {
		newEco[i] = make([]*Unit, numCols)
		for j := 0; j < numCols; j++ {

			// initialize the pointer
			newEco[i][j] = new(Unit)

			// generate food randomly. 50% chance of generating food at every location in initial system
			randomFood := rand.Float64()
			if randomFood > 0.90 {
				newEco[i][j].food.isPresent = true
			}
		}
	}

	InitializePreyAndPredator(numRows, numCols, numPrey, numPred, &newEco)

	return newEco
}
