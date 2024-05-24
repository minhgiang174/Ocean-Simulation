package main

import (
	"fmt"
	"gifhelper"
	"math/rand"
	"time"
)

// GLOBAL VARIABLES

// Set a constant dictionary where keys are the directionIndex and the values are the orderedPair with corresponding deltaX and deltaY
var deltas map[int]OrderedPair
var energyCosts map[int]int
var maxEnergy int = 1500
var energyThresholdPrey int = 50
var ageThresholdPrey int = 21
var costOfLivingPrey int = 0
var energyGainedPerPlankton int = 50
var energyThresholdPredator int = 100 // 800
var ageThresholdPredator int = 42     // 50
var costOfLivingPredator int = 0

// DON'T MESS WITH THIS. SET THEM IN MAIN
// we will use these to track numPrey and numPred globally
var numPrey int = 0
var numPred int = 0

func main() {

	// assign deltas
	deltas = map[int]OrderedPair{
		0: {-1, 1},
		1: {0, 1},
		2: {1, 1},
		3: {-1, 0},
		4: {1, 0},
		5: {-1, -1},
		6: {0, -1},
		7: {1, -1},
	}

	energyCosts = map[int]int{
		0: 0,
		1: -1,
		2: -2,
		3: -4,
		4: -8,
		5: -4,
		6: -2,
		7: -1,
	}

	// var numRows int = 250
	// var numCols int = 250
	var numRows int = 50
	var numCols int = 50
	numPrey = 10
	numPred = 50

	if (numRows * numCols) < (numPrey + numPred) {
		panic("there's too many predator and prey in total")
	}
	var initialEcosystem Ecosystem = InitializeEcosystem(numRows, numCols, numPrey, numPred)
	var totalTimesteps int = 10
	var foodRule string = "gardenOfEden"

	// seed the PRNG approximately randomly
	rand.Seed(time.Now().UnixNano())
	allEcosystems := SimulateEcosystemEvolution(&initialEcosystem, totalTimesteps, foodRule)

	canvasWidth := 1000
	frequency := 1
	scalingFactor := 1.0
	imageList := AnimateSystem(allEcosystems, canvasWidth, frequency, scalingFactor)

	gifhelper.ImagesToGIF(imageList, "ecosystem")
	fmt.Println("GIF drawn.")

	// use this for debugging and seeing characteristics of specific ecosystem(s)
	// PrintEcosystem(allEcosystems[len(allEcosystems)-1])

}

func PrintEcosystem(someEcosystem *Ecosystem) {
	countPrey, countPred := 0, 0
	checkIfTwo := false
	for i := range *(someEcosystem) {
		for j := range (*someEcosystem)[i] {
			checkIfTwo = false
			curUnit := (*someEcosystem)[i][j]
			// fmt.Println(i, j, "=")
			// if curUnit.food.isPresent {
			// 	fmt.Println("i,j are ", i, j, " = food")
			// }

			if curUnit.predator != nil {
				// fmt.Println("i,j are ", i, j, " = predator")
				countPred += 1
				checkIfTwo = true
			}

			if curUnit.prey != nil {
				// fmt.Println("i,j are ", i, j, " = prey")
				countPrey += 1
				if checkIfTwo {
					panic("there are predator and prey in same cell")
				}
			}
		}
	}
	fmt.Println("Number prey =", countPrey)
	fmt.Println("Number pred =", countPred)
}
