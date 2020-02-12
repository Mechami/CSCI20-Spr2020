package main

import (
	"os"
	"fmt"
	"math"
	"math/rand"
	"strings"
)

// Golden Ratio
// Computed at runtime to ensure accuracy up to ~15 decimal places
var goldenRatio float64 = (1.0 + math.Sqrt(5.0)) / 2.0

// strScan
// Scans in a string and lowers the case of all characters
func strScan(response *string) {
	// Scan the string into the specified pointer
	fmt.Scanf("%s\n", response)

	// Then overwrite the pointed string with the corrected string
	*response = strings.ToLower(*response)
}

// strPrompt
// Reads in a boolean answer from the user with a specified prompt string
// Accepts any string that starts with [YyNn]
func strPrompt(prompt string) bool {
	// Allocate a variable to hold the entered string
	response := ""

	// For each iteration check the response string
	// If it begins with lowercase 'y' or lowercase 'n' do nothing and move on (strScan will sanitize the input for us)
	// Otherwise keep scanning string inputs from the user until they give a y/n response
	for ; !(strings.HasPrefix(response, "y") || strings.HasPrefix(response, "n")); strScan(&response) {
		// Prompt the user with the specified prompt string and a y/n to hint what we want from them
		fmt.Printf("\r%s (y/n)\t", prompt)
	}

	// Check if the response starts with 'y' and return that as the prompt value
	return strings.HasPrefix(response, "y")
}

// readDice
// Reads in an integer from the user within the range [4, 100]
func readDice() int {
	// Declare an integer and initialize it not zero
	var size int = ^0

	// While the size is outside the allowed range
	// Keep scanning inputs from the user and attempt to store them in size
	for ; (size < 4) || (size > 100); fmt.Scanf("%d\n", &size) {
		// Display a hint to the user for the accepted values
		fmt.Printf("\rSize of die? [4, 100]\t")
	}

	// Finally return the die size to the caller
	return size
}

// procTurn
// Processes a player's turn using the specified die size and prompt function
func procTurn(dieSize int, promptFunc func() bool) int {
	// Default to a net score of zero
	turnScore := 0

	// Assume the player wants to keep rolling
	keepRolling := true

	// While the player wants to keep rolling
	for keepRolling {
		// Roll a random number between [1, dieSize]
		roll := 1 + rand.Intn(dieSize)

		// If it rolled anything other than a one
		if roll > 1 {
			// Add it to the player's score
			turnScore += roll
		// Otherwise if did roll a one
		} else {
			// The player gets nothing
			turnScore = 0
		}

		// Then print out the roll statistics to the user
		fmt.Printf("Roll: %d\tGain: %d\n", roll, turnScore)

		// And roll again if the player did not roll a one or decide to hold
		keepRolling = (roll > 1) && promptFunc()
	}

	// Finally return the players accumulated score
	return turnScore
}

func main() {
	// Tell the player what game we are playing
	fmt.Println("Game of Pig - Singleplayer")

	// Assume the player wants to keep playing and ask every iteration
	// For every iteration the player continues to play
	for playing := true; playing; playing = strPrompt("Play again?") {
		// Create and initialize some score variables
		scorePC, scoreCPU := 0, 0

		// Then read in a dice size from the user
		diceSize := readDice()

		// CPU goes first for arbitrary reasons
		playersTurn := false

		// While both players have a score less than 100
		for (scorePC < 100) && (scoreCPU < 100) {
			// Switch the player turn based on boolean
			switch playersTurn {

			// If it's the [human] player's turn
			case true:
				{
					// Tell them so
					fmt.Println("# Player's turn:")

					// Then process their turn with the chosen dice size
					// A lambda that calls the prompt function will provide the continue boolean
					scorePC += procTurn(diceSize, func() bool { return strPrompt("Roll again?") })

					// Then tell the player how much they scored
					fmt.Printf("You scored: %d\n\n", scorePC)
				}

			// If it's the computer's turn
			case false:
				{
					// Tell the [human] player
					fmt.Println("# CPU's turn:")

					// Then process their turn as before
					// But instead of reading in from the keyboard generate a random float from [0, math.MaxFloat64)
					// On an exponentially distributed range, then compare that against the golden ratio
					// If less than 1/Gr then continue, otherwise hold
					scoreCPU += procTurn(diceSize, func() bool { return rand.ExpFloat64() < (1.0 / goldenRatio) })

					// Finally tell the [human] player what the CPU scored
					fmt.Printf("CPU scored: %d\n\n", scoreCPU)
				}
			// Easter egg: If somehow neither the true nor false cases are matched
			default:
				{
					// Tell the user they have done the unexpected
					fmt.Printf("You have accomplished the implausible.\n")

					// Then proceed to bail out back to the OS
					os.Exit(1)
				}
			}

			// Finally switch which players turn is next
			playersTurn = !playersTurn
		}

		// Define an empty win message string
		winMsg := ""

		// If the [human] player scored more
		if scorePC > scoreCPU {
			// Tell them they win
			winMsg = "You win!"
		
		// Otherwise the computer wins
		} else {

			// And inform the user
			winMsg = "The computer wins!"
		}

		// Finally print out a final score value of both players
		fmt.Printf("Final score: %d\nCPU's Score: %d\n%s\n", scorePC, scoreCPU, winMsg)
	}
}
