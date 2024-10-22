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
		score, _ := strconv.Atoi(record[1]) // Convert score from string to int
		minSpeed, _ := strconv.ParseFloat(record[2], 64)
		maxSpeed, _ := strconv.ParseFloat(record[3], 64)
		players = append(players, Player{Name: record[0], Score: score, MinSpeed: minSpeed, MaxSpeed: maxSpeed, UUID: record[4]})
	}
	
	return players, nil
}


type Player struct {
    Name        string
    Distance    float64
    Finished    bool
    Place       int
    MinSpeed    float64
    MaxSpeed    float64
    Score       int
    UUID        string
    Endurance   float64
    Resting     bool
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
			Score:    0,
		}
		players = append(players, player)
	}

	return players, nil
}

func RunSimulation(app fyne.App, numberOfPlayers int, laneHeight int, windowWidth int, playerData [][]string, raceLengthEntry string) {
	// Convert playerData to []Player
	players, err := CreatePlayers(playerData[1:])
	if err != nil {
		return
	}

	// Convert race length from string to int, and handle any potential error
	raceLength, err := strconv.Atoi(raceLengthEntry)
	if err != nil {
		return
	}

	// Start the race with the created players and parsed race length
	DrawRaceTrack(app, numberOfPlayers, laneHeight, float32(windowWidth),players, raceLength)
}

func RandomFloat(lowerLimit, upperLimit float64) float64 { 
    return lowerLimit + rand.Float64()*(upperLimit-lowerLimit)
}