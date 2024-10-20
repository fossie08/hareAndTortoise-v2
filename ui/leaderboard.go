package ui

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Player struct {
	Name     string
	Score    int
	MinSpeed float64
	MaxSpeed float64
	UUID     string
}

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

// UpdateCSV updates the CSV file with the modified player data
func UpdateCSV(filename string, players []Player) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Name", "Score", "Min Speed", "Max Speed", "UUID"})
	for _, player := range players {
		record := []string{
			player.Name,
			strconv.Itoa(player.Score),
			strconv.FormatFloat(player.MinSpeed, 'f', -1, 64),
			strconv.FormatFloat(player.MaxSpeed, 'f', -1, 64),
			player.UUID,
		}
		writer.Write(record)
	}
	return nil
}

func ShowEmbeddedEditForm(list *widget.Table, player *Player, players []Player, filename string, playerData *[][]string) {
	formWindow := fyne.CurrentApp().NewWindow("Edit Player")
	// Create entries for editing player data
	nameEntry := widget.NewEntry()
	nameEntry.SetText(player.Name)

	minSpeedEntry := widget.NewEntry()
	minSpeedEntry.SetText(strconv.FormatFloat(player.MinSpeed, 'f', -1, 64))

	maxSpeedEntry := widget.NewEntry()
	maxSpeedEntry.SetText(strconv.FormatFloat(player.MaxSpeed, 'f', -1, 64))

	saveButton := widget.NewButton("Save", func() {
		player.Name = nameEntry.Text
		player.MinSpeed, _ = strconv.ParseFloat(minSpeedEntry.Text, 64)
		player.MaxSpeed, _ = strconv.ParseFloat(maxSpeedEntry.Text, 64)

		// Save the changes back to the CSV file
		if err := SavePlayersToCSV(filename, players); err != nil {
			fmt.Println("Error saving to CSV:", err)
		}

		// Refresh the playerData after saving
		*playerData = [][]string{{"Name", "Score", "Min Speed", "Max Speed", "UUID"}} // Header row
		for _, p := range players {
			*playerData = append(*playerData, []string{p.Name, strconv.Itoa(p.Score), strconv.FormatFloat(p.MinSpeed, 'g', -1, 64), strconv.FormatFloat(p.MaxSpeed, 'g', -1, 64), p.UUID})
		}

		list.Refresh() // Refresh the list with the updated playerData
		formWindow.Close()
	})

	// Create the edit form
	editForm := container.NewVBox(
		widget.NewLabel("Edit Player Details"),
		widget.NewForm(
			widget.NewFormItem("Name", nameEntry),
			widget.NewFormItem("Min Speed", minSpeedEntry),
			widget.NewFormItem("Max Speed", maxSpeedEntry),
		),
		saveButton,
	)

	// Create a pop-up window or panel in your main UI to display the form
	// Since you cannot replace the table content, you should add this form to a new container in your UI
	formWindow.SetContent(editForm)
	formWindow.Resize(fyne.NewSize(300, 200))
	formWindow.Show()
}

// UpdateLeaderboardContent dynamically updates the table with new data
func UpdateLeaderboardContent(list *widget.Table, playerData [][]string) {
	list.Refresh()
}

func DisplayLeaderboard() *fyne.Container {
	playerData := [][]string{{"Name", "Score", "Min Speed", "Max Speed", "UUID"}} // Header row
	players, err := ReadCSV("data/animal.simulation")
	if err != nil {
		fmt.Println("Error loading leaderboard:", err)
		return nil
	}

	for _, player := range players {
		playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score), strconv.FormatFloat(player.MinSpeed, 'f', -1, 64), strconv.FormatFloat(player.MaxSpeed, 'f', -1, 64), player.UUID})
	}

	// Create a widget to show leaderboard data
	list := widget.NewTable(
		func() (int, int) { return len(playerData), 5 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(playerData[id.Row][id.Col])
			if id.Row == 0 { // Bold headers
				o.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
			}
		},
	)

	editButton := widget.NewButton("Edit Player", func() {
		playerNames := make([]string, len(players))
		for i, player := range players {
			playerNames[i] = player.Name
		}
	
		// Create a dropdown to select a player
		playerSelect := widget.NewSelect(playerNames, func(selected string) {
			for i, player := range players {
				if player.Name == selected {
					ShowEmbeddedEditForm(list, &players[i], players, "data/animal.simulation", &playerData) // Pass the list and playerData by reference
					break
				}
			}
		})
	
		// Create a container for the dropdown
		form := container.NewVBox(
			widget.NewLabel("Select a Player to Edit:"),
			playerSelect,
		)
	
		// Create and show the edit window
		editWindow := fyne.CurrentApp().NewWindow("Select Player")
		editWindow.SetContent(form)
		editWindow.Resize(fyne.NewSize(300, 200)) // Resize the window if needed
		editWindow.Show()
	})
	
	// Sorting buttons
	sortByScoreButton := widget.NewButton("Sort by Score", func() {
		sort.Slice(players, func(i, j int) bool {
			return players[i].Score > players[j].Score // Sort by score (descending)
		})
		// Update the playerData slice after sorting
		playerData = [][]string{{"Name", "Score", "Min Speed", "Max Speed", "UUID"}} // Header row
		for _, player := range players {
			playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score), strconv.FormatFloat(player.MinSpeed, 'g', -1, 64), strconv.FormatFloat(player.MaxSpeed, 'g', -1, 64), player.UUID})
		}
		UpdateLeaderboardContent(list, playerData) // Refresh the list with the updated playerData
	})
	
	sortByNameButton := widget.NewButton("Sort by Name", func() {
		sort.Slice(players, func(i, j int) bool {
			return players[i].Name < players[j].Name // Sort by name (alphabetical)
		})
		// Update the playerData slice after sorting
		playerData = [][]string{{"Name", "Score", "Min Speed", "Max Speed", "UUID"}} // Header row
		for _, player := range players {
			playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score), strconv.FormatFloat(player.MinSpeed, 'g', -1, 64), strconv.FormatFloat(player.MaxSpeed, 'g', -1, 64), player.UUID})
		}
		UpdateLeaderboardContent(list, playerData) // Refresh the list with the updated playerData
	})
	
	sortByUUIDButton := widget.NewButton("Sort by UUID", func() {
		sort.Slice(players, func(i, j int) bool {
			return players[i].UUID < players[j].UUID // Sort by UUID (alphabetical)
		})
		// Update the playerData slice after sorting
		playerData = [][]string{{"Name", "Score", "Min Speed", "Max Speed", "UUID"}} // Header row
		for _, player := range players {
			playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score), strconv.FormatFloat(player.MinSpeed, 'g', -1, 64), strconv.FormatFloat(player.MaxSpeed, 'g', -1, 64), player.UUID})
		}
		UpdateLeaderboardContent(list, playerData) // Refresh the list with the updated playerData
	})

	// Refresh button
	refreshButton := widget.NewButton("Refresh", func() {
		players, err = ReadCSV("data/animal.simulation")
		if err != nil {
			fmt.Println("Error refreshing leaderboard:", err)
			return
		}
		playerData = [][]string{{"Name", "Score", "Min Speed", "Max Speed", "UUID"}} // Header row
		for _, player := range players {
			playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score), strconv.FormatFloat(player.MinSpeed, 'g', -1, 64), strconv.FormatFloat(player.MaxSpeed, 'g', -1, 64), player.UUID})
		}
		list.Refresh() // Refresh the list with the updated playerData
	})

	// Setting column widths
	list.SetColumnWidth(0, 140)
	list.SetColumnWidth(1, 140)
	list.SetColumnWidth(2, 140)
	list.SetColumnWidth(3, 140)
	list.SetColumnWidth(4, 280)

	// Display list and sorting buttons
	content := container.NewBorder(
		container.NewHBox(refreshButton, editButton, sortByNameButton, sortByScoreButton, sortByUUIDButton),
		nil, nil, nil,
		list,
	)

	return content
}




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
