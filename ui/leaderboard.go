package ui

import (
	"encoding/csv"
//	"fmt"
//	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
//	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"os"
	"sort"
	"strconv"
)

type Player struct {
	Name  string
	Score int
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
		players = append(players, Player{Name: record[0], Score: score})
	}
	return players, nil
}

// UpdateLeaderboardContent dynamically updates the table with sorted data
func updateLeaderboardContent(list *widget.Table, playerData [][]string) {
	list.Refresh()
}

// DisplayLeaderboard creates a Fyne window displaying the players sorted by a selected order
func DisplayLeaderboard(myApp fyne.App, window fyne.Window, players []Player) {
	// Add header row to the playerData slice
	playerData := [][]string{{"Name", "Score"}} // Header row

	// Fill player data (skipping the header)
	for _, player := range players {
		playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score)})
	}

	// Create a widget to show leaderboard data
	list := widget.NewTable(
		func() (int, int) { return len(playerData), 2 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(playerData[i.Row][i.Col])
			if i.Row == 0 { // Bold headers
				o.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
			}
		},
	)

	// Sorting buttons
	sortByScoreButton := widget.NewButton("Sort by Score", func() {
		sort.Slice(players, func(i, j int) bool {
			return players[i].Score > players[j].Score // Sort by score (descending)
		})
		// Update the playerData slice after sorting
		playerData = [][]string{{"Name", "Score"}} // Recreate header
		for _, player := range players {
			playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score)})
		}
		list.Refresh()
	})

	sortByNameButton := widget.NewButton("Sort by Name", func() {
		sort.Slice(players, func(i, j int) bool {
			return players[i].Name < players[j].Name // Sort by name (alphabetical)
		})
		// Update the playerData slice after sorting
		playerData = [][]string{{"Name", "Score"}} // Recreate header
		for _, player := range players {
			playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score)})
		}
		list.Refresh()
	})

	list.SetColumnWidth(0, 140)
	list.SetColumnWidth(1, 140)

	// Layout: display list and sorting options
	content := container.NewBorder(
		container.NewHBox(sortByScoreButton, sortByNameButton),
		nil, nil, nil,
		list,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(400, 500)) // Adjust size as needed
	window.CenterOnScreen()
}