package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"fmt"
	"hareandtortoise/v2/simulation"
)

// ShowSetupRaceMenu shows the setup race menu, allowing users to select animals and race parameters
func ShowSetupRaceMenu(app fyne.App) [][]string {
	setupWindow := app.NewWindow("Setup Race")

	// Use the custom ReadCSV function to load animals
	players, err := ReadCSV("data/animal.simulation")
	if err != nil {
		dialog.ShowError(err, setupWindow)
		return nil
	}

	// Create animal selection checkboxes and convert them to fyne.CanvasObject
	var selectedAnimals []Player
	animalCheckboxes := make([]fyne.CanvasObject, len(players)) // This should be []fyne.CanvasObject
	for i, player := range players {
		checkbox := widget.NewCheck(player.Name, func(checked bool) {
			if checked {
				selectedAnimals = append(selectedAnimals, player)
			} else {
				// Remove player if unchecked
				for j, p := range selectedAnimals {
					if p.Name == player.Name {
						selectedAnimals = append(selectedAnimals[:j], selectedAnimals[j+1:]...)
						break
					}
				}
			}
		})
		animalCheckboxes[i] = checkbox // Assign as a fyne.CanvasObject
	}

	// Race length entry
	raceLengthLabel := widget.NewLabel("Race Length (meters):")
	raceLengthEntry := widget.NewEntry()
	raceLengthEntry.SetPlaceHolder("Enter race length")

	// Start Race button
	startRaceButton := widget.NewButton("Start Race", func() {
		if len(selectedAnimals) == 0 {
			dialog.ShowInformation("Error", "Please select at least one animal for the race.", setupWindow)
			return
		}
		if raceLengthEntry.Text == "" {
			dialog.ShowInformation("Error", "Please enter a valid race length.", setupWindow)
			return
		}

		// Create playerData in the specified format
		playerData := [][]string{{"Name", "Score", "Min Speed", "Max Speed", "UUID"}} // Header row
		for _, player := range selectedAnimals {
			playerData = append(playerData, []string{
				player.Name,
				strconv.Itoa(player.Score),
				strconv.FormatFloat(player.MinSpeed, 'g', -1, 64),
				strconv.FormatFloat(player.MaxSpeed, 'g', -1, 64),
				player.UUID,
			})
		}

		fmt.Println("Race Length:", raceLengthEntry.Text)
		fmt.Println("Selected Players:", playerData)
		simulation.RunSimulation(playerData, raceLengthEntry.Text)
		// Close the window
		setupWindow.Close()

		// Returning player data for further use
	})

	// Organize UI components
	content := container.NewVBox(
		widget.NewLabel("Select Animals:"),
		container.NewVBox(animalCheckboxes...), // Pass converted checkboxes
		raceLengthLabel,
		raceLengthEntry,
		startRaceButton,
	)

	setupWindow.SetContent(content)
	setupWindow.Resize(fyne.NewSize(400, 400))
	setupWindow.CenterOnScreen()
	setupWindow.Show()
	// Return empty playerData initially, will be updated when the race starts
	return nil
}
