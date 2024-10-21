package ui

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
)

// Struct to hold participant details
type ParticipantDetails struct {
	UUID     string
	Place    string
	Distance string
	Score    string
}

// Struct to hold overall race details
type RaceDetails struct {
	UUID          string
	TotalDistance string
	Rounds        string
	Date          string
	Time          string
	Participants  []ParticipantDetails
}


// Function to load all .simulation files from data/ and display in a window
func ShowPreviousRacesWindow(app fyne.App, mainWindow fyne.Window) {
	// Create a new window for showing previous races
	previousRacesWindow := app.NewWindow("Previous Races")

	// Container for the list of races
	raceContainer := container.NewVBox()

	// Read all .simulation files from the data folder
	files, err := ioutil.ReadDir("data/")
	if err != nil {
		dialog.ShowError(err, mainWindow)
		return
	}

	// Loop through files and load race details
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".simulation") {
			raceDetails, err := LoadRaceDetails("data/" + file.Name())
			if err != nil {
				fmt.Println("Failed to load race details:", err)
				continue
			}

			// Display the race details in the UI
			raceLabel := fmt.Sprintf("UUID: %s | Date: %s | Rounds: %s | Total Distance: %s", raceDetails.UUID, raceDetails.Date, raceDetails.Rounds, raceDetails.TotalDistance)
			raceContainer.Add(widget.NewLabel(raceLabel))

			// Display the participants and their positions
			for _, participant := range raceDetails.Participants {
				participantLabel := fmt.Sprintf("Place: %s | Name: %s", participant.Place, participant.UUID)
				raceContainer.Add(widget.NewLabel(participantLabel))
			}

			// Add a separator for each race
			raceContainer.Add(widget.NewSeparator())
		}
	}

	// Set the content of the window and show it
	previousRacesWindow.SetContent(raceContainer)
	previousRacesWindow.Resize(fyne.NewSize(500, 400)) // Adjust window size as needed
	previousRacesWindow.CenterOnScreen()
	previousRacesWindow.Show()
}

func LoadRaceDetails(filePath string) (RaceDetails, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return RaceDetails{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read() // Read the first line as headers
	if err != nil {
		return RaceDetails{}, err
	}

	// Ensure the headers have at least 8 fields
	if len(headers) < 8 {
		return RaceDetails{}, fmt.Errorf("invalid header structure")
	}

	// Collect race details
	raceDetails := RaceDetails{
		UUID:          headers[0], // Assuming the file name contains UUID, but using the first column as UUID
		TotalDistance: headers[4],
		Rounds:        headers[5],
		Date:          headers[6],
		Time:          headers[7],
	}

	var participants []ParticipantDetails
	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}

		// Ensure the record has at least 4 fields for participants (UUID, Place, Distance, Score)
		if len(record) < 4 {
			continue
		}

		// Collect participant data
		participant := ParticipantDetails{
			UUID:     record[0],
			Place:    record[1],
			Distance: record[2],
			Score:    record[3],
		}
		participants = append(participants, participant)
	}

	raceDetails.Participants = participants
	return raceDetails, nil
}

