package simulation

import (
	"github.com/google/uuid"
	"encoding/csv"
	"fmt"
	"os"
)

func WriteCSV(filename string, data [][]string, appendMode bool) error {
	var file *os.File
	var err error

	if appendMode {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	} else {
		file, err = os.Create(filename) // Overwrite if not appendMode
	}

	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func CreateAnimal () {
	id := uuid.New().String()
	data := [][]string{{"Name","0",id}}
	err := WriteCSV("data/leaderboard.simulation", data, true)
	if err != nil {
		fmt.Println("Error:", err)
	}
}