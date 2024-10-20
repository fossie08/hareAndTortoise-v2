package simulation

import (
	"encoding/csv"
	"os"
	"fmt"
	"strconv"
	//"time"
	"math/rand"
	//"hareandtortoise/v2/ui"
	"fyne.io/fyne/v2"
)

// ReadCSV reads the CSV file and returns a slice of Players
func ReadCSV(filename string) ([]Player, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var players []Player
	for _, record := range records[1:] { // Skipping the header in the CSV file
		fmt.Println("Processing record:", record) // Debugging line to check each record
		score, _ := strconv.Atoi(record[1]) // Convert score from string to int
		minSpeed, _ := strconv.ParseFloat(record[2], 64)
		maxSpeed, _ := strconv.ParseFloat(record[3], 64)
		players = append(players, Player{Name: record[0], Score: score, MinSpeed: minSpeed, MaxSpeed: maxSpeed, UUID: record[4]})
	}
	
	return players, nil
}

type Player struct {
	Name     string
	Score    int
	MinSpeed float64
	MaxSpeed float64
	UUID     string
	Distance float64 // to track how far they've gone
	Finished bool    // to track if the player has finished the race
	Place    int     // to track the finishing position
}

// Convert playerData ([][]string) to []Player
func CreatePlayers(playerData [][]string) ([]Player, error) {
	var players []Player

	for _, data := range playerData {
		minSpeed, err := strconv.ParseFloat(data[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid min speed for player %s: %v", data[0], err)
		}

		maxSpeed, err := strconv.ParseFloat(data[3], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid max speed for player %s: %v", data[0], err)
		}

		player := Player{
			Name:     data[0],
			MinSpeed: minSpeed,
			MaxSpeed: maxSpeed,
			UUID:     data[4],
			Distance: 0,
			Finished: false,
			Place:    0,
		}
		players = append(players, player)
	}

	return players, nil
}

func RunSimulation(app fyne.App, numberOfPlayers int, laneHeight int, windowWidth int, playerData [][]string, raceLengthEntry string) {
	// Convert playerData to []Player
	players, err := CreatePlayers(playerData[1:])
	if err != nil {
		fmt.Println("Error creating players:", err)
		return
	}

	// Convert race length from string to int, and handle any potential error
	raceLength, err := strconv.Atoi(raceLengthEntry)
	if err != nil {
		fmt.Println("Invalid race length:", err)
		return
	}

	// Start the race with the created players and parsed race length
	DrawRaceTrack(app, numberOfPlayers, laneHeight, float32(windowWidth),players, raceLength)
}
/*
// Simulate the race with the given players
func StartRace(players []Player, totalDistance int) {
	rand.Seed(time.Now().UnixNano())

	// Reset each player's distance and status
	for i := range players {
		players[i].Distance = 0
		players[i].Finished = false
		players[i].Place = 0
	}

	round := 1
	finishedPlayers := 0
	currentPlace := 1 // Tracks the finishing position

	for finishedPlayers < len(players) {
		for i := range players {
			if !players[i].Finished {
				// Random movement for each player within their speed range
				players[i].Distance += RandomFloat(players[i].MinSpeed, players[i].MaxSpeed)

				// Check if the player has finished the race
				if players[i].Distance >= float64(totalDistance) {
					players[i].Finished = true
					players[i].Place = currentPlace
					currentPlace++
					finishedPlayers++
				}
			}
		}

		// Print the race progress to the console
		fmt.Printf("Round %d\n", round)
		for _, player := range players {
			status := ""
			if player.Finished {
				status = fmt.Sprintf(" (Finished in place %d)", player.Place)
			}
			fmt.Printf("%s: %.2f m%s\n", player.Name, player.Distance, status)
		}
		fmt.Println()

		time.Sleep(10 * time.Millisecond) // Pause to simulate time between rounds

		round++
	}

	// Print the final results in order of placement
	fmt.Println("Race Finished!")
	fmt.Println("Final Placements:")
	for _, player := range players {
		fmt.Printf("%d. %s\n", player.Place, player.Name)
	}
}
*/
func RandomFloat(lowerLimit, upperLimit float64) float64 { 
    return lowerLimit + rand.Float64()*(upperLimit-lowerLimit)
}