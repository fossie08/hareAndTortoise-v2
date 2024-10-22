package simulation

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// SavePlayersToCSV updates the relevant player data without overwriting the whole file
func SavePlayersToCSV(filename string, players []Player) error {
	// Open the existing file for reading
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll() // Read all records at once
	if err != nil {
		return err
	}

	// Create a map of players by UUID for quick lookup
	playerMap := make(map[string]Player)
	for _, player := range players {
		playerMap[player.UUID] = player
	}

	// Modify only the relevant player records
	for i, record := range records {
		if i == 0 {
			// Skip header row
			continue
		}

		uuid := record[4] // Assuming UUID is the 5th column
		if updatedPlayer, ok := playerMap[uuid]; ok {
			// Replace the player's data with updated values
			records[i] = []string{
				updatedPlayer.Name,
				strconv.Itoa(updatedPlayer.Score),
				strconv.FormatFloat(updatedPlayer.MinSpeed, 'f', -1, 64),
				strconv.FormatFloat(updatedPlayer.MaxSpeed, 'f', -1, 64),
				updatedPlayer.UUID,
			}
		}
	}

	// Open the file for writing (overwrite)
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Write all records (including updated ones) back to the file
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
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
	}

}

