package simulation

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
	"strconv"
)

func SavePlayersToCSV(filename string, players []Player) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Name", "Score", "Min Speed", "Max Speed", "UUID"})

	// Write player data
	for _, player := range players {
		writer.Write([]string{
			player.Name,
			strconv.Itoa(player.Score),
			strconv.FormatFloat(player.MinSpeed, 'f', -1, 64),
			strconv.FormatFloat(player.MaxSpeed, 'f', -1, 64),
			player.UUID,
		})
	}

	return nil
}

// Save the race results to a CSV file, updating the existing score
func SaveRaceResults(players []Player, totalDistance, numRounds int, uuid string) {
	filePath := fmt.Sprintf("data/%s.simulation", uuid)

	// Get the current date and time
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Open or create the CSV file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers, including Date and Time
	writer.Write([]string{"UUID", "Place", "Distance Travelled", "Score", "Total Distance", "Rounds", "Date", "Time"})

	// Write player data
	for _, player := range players {
		record := []string{
			player.UUID,
			fmt.Sprintf("%d", player.Place),
			fmt.Sprintf("%.1f", player.Distance),
			fmt.Sprintf("%d", player.Score),
			fmt.Sprintf("%d", totalDistance),
			fmt.Sprintf("%d", numRounds),
			currentTime[:10], // Date
			currentTime[11:], // Time
		}
		writer.Write(record)
	}

	// Update the player scores in "data/animal.simulation"
	if err := SavePlayersToCSV("data/animal.simulation", players); err != nil {
		fmt.Println("Failed to update scores in animal.simulation:", err)
	}

	fmt.Println("Race results saved to:", filePath)
}

