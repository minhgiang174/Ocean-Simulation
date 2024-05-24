package main

import (
	"math/rand"
	"time"
)

// GeneratePreyFoodProbabilistically() is a method operating on a Unit pointer someUnit. it uses some probabilistic function determined by foodRule, to determine whether food will be generated in this Unit or not. NOTE: this function shouldn't be called if there is something else in the Unit already
// Input: foodRule string, row and col indices for the Unit, and a Ecosystem pointer someEcosystem
// Output: none. operates on a pointer
func (someUnit *Unit) GeneratePreyFoodProbabilistically(foodRule string, row, col int, someEcosystem *Ecosystem) {
	//don't want a race condition, where all processes share a single PRNG object
	source := rand.NewSource(time.Now().UnixNano()) // create new PRNG object
	generator := rand.New(source)
	numRows := someEcosystem.CountRows()
	numCols := someEcosystem.CountCols()

	if foodRule == "gardenOfEden" {
		someUnit.GenerateEden(row, col, numRows, numCols, generator)
	} else if foodRule == "even" {
		someUnit.GenerateRandom(row, col, numRows, numCols, generator)
	} else if foodRule == "lineRunner" {
		someUnit.GenerateLineRunner(row, col, numRows, numCols, generator)
	} else {
		panic("invalid foodRule string inputted. should be eden, random, or lineRunner!")
	}
}

func (someUnit *Unit) GenerateRandom(row, col, numRows, numCols int, generator *rand.Rand) {
	// generate a random floating point value on half open interval [0,1), with the unique generator
	probability := generator.Float64() // this is instead of rand.Float64() accessing the global random PRNG object
	// if meets the probability, randomly make food.
	if probability >= 0.999 {
		someUnit.food.isPresent = true
	}
}

// GenerateEden() is a method operating on a Unit pointer someUnit. probabilistically determines how often food appears, with high weighting towards the inner rectangle of the board
// Input: foodRule string, row and col indices for the Unit
// Output: none. operates on a pointer
func (someUnit *Unit) GenerateEden(row, col, numRows, numCols int, generator *rand.Rand) {

	// initialize the center of the board
	centerRow := numRows / 2
	centerCol := numCols / 2
	/* if ecosystem has length l and width w then it will have
	area = l * w
	and the below parameters means the eden rectangle will have
	area = (l/fractionOfBoard)* (w/fractionOfBoard)*/
	fractionOfBoard := 10
	halfCenterRecLength := numCols / (fractionOfBoard)
	halfCenterRecWidth := numRows / (fractionOfBoard)

	// generate a random floating point value on half open interval [0,1), with the unique generator
	probability := generator.Float64() // this is instead of rand.Float64() accessing the global random PRNG object

	// check if the row and col of the current unit is within the center rectangle
	if CheckIsInCenter(row, col, centerRow, centerCol, halfCenterRecLength, halfCenterRecWidth) {
		// if within the center rectangle then much higher likelihood of generating food
		if probability >= 0.90 {
			someUnit.food.isPresent = true
		}
	} else {
		// if not within the center rectangle then much less likely to generate food
		if probability >= 0.99 {
			someUnit.food.isPresent = true
		}
	}
}

// checks if the row and col indices lie within the center of the ecosystem, return true if in central rectangle and false otherwise
func CheckIsInCenter(row, col, centerRow, centerCol, halfCenterRecLength, halfCenterRecWidth int) bool {
	if (centerRow-halfCenterRecWidth) <= row && row <= centerRow+halfCenterRecWidth {
		if (centerCol-halfCenterRecLength) <= col && col <= (centerCol+halfCenterRecLength) {
			// we are within the center, return true
			return true
		}
	}

	// we are not within the center, return false
	return false
}

// GenerateLineRunner: Method that sets grid for food eaters to run along
func (someUnit *Unit) GenerateLineRunner(row, col, numRows, numCols int, generator *rand.Rand) {
	// generate a random floating point value on half open interval [0,1), with the unique generator
	probability := generator.Float64() // this is instead of rand.Float64() accessing the global random PRNG object
	gridRow := numRows / 4
	gridCol := numCols / 4

	// check if the row and col of the current unit is within the center rectangle
	if CheckIsOnGridLine(&row, &col, &gridRow, &gridCol) {
		// if on the grid line then much higher likelihood of generating food
		if probability >= 0.95 {
			someUnit.food.isPresent = true
		}
	} else {
		// if not within the center rectangle then much less likely to generate food
		if probability >= 0.99999 {
			someUnit.food.isPresent = true
		}
	}
}

func CheckIsOnGridLine(row, col, gridRow, gridCol *int) bool {
	if *row == 0 || *col == 0 {
		return false
	}
	if *row%*gridRow == 0 {
		return true
	}
	if *col%*gridCol == 0 {
		return true
	}
	return false

}
