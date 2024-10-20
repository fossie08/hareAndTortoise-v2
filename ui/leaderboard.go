package ui

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Player struct {
	Name  string
	Score int
	MinSpeed float64
	MaxSpeed float64
	UUID string
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

// UpdateLeaderboardContent dynamically updates the table with sorted data
func UpdateLeaderboardContent(list *widget.Table, playerData [][]string) {
	list.Refresh()
}

// DisplayLeaderboard creates a Fyne window displaying the players sorted by a selected order
func DisplayLeaderboard() (content *fyne.Container) {
	// Add header row to the playerData slice
	playerData := [][]string{{"Name", "Score", "Min Speed", "Max Speed", "UUID"}} // Header row
	players, err := ReadCSV("data/animal.simulation")
	if err != nil {
		fmt.Println("Error loading leaderboard:", err)
		return
	}
	// Fill player data (skipping the header)
	for _, player := range players {
		playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score), strconv.FormatFloat(player.MinSpeed, 'g', -1, 64), strconv.FormatFloat(player.MaxSpeed, 'g', -1, 64), player.UUID})
	}

	// Create a widget to show leaderboard data
	list := widget.NewTable(
		func() (int, int) { return len(playerData), 5 },
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

	refreshLeaderboard := func() {
		players, err = ReadCSV("data/animal.simulation")
		if err != nil {
			fmt.Println("Error refreshing leaderboard:", err)
			return
		}
		playerData = [][]string{{"Name", "Score", "Min Speed", "Max Speed", "UUID"}} // Header row
		for _, player := range players {
			playerData = append(playerData, []string{player.Name, strconv.Itoa(player.Score), strconv.FormatFloat(player.MinSpeed, 'g', -1, 64), strconv.FormatFloat(player.MaxSpeed, 'g', -1, 64), player.UUID})
		}
		UpdateLeaderboardContent(list, playerData)
	}

	// Refresh button
	refreshButton := widget.NewButton("Refresh", func() {
		refreshLeaderboard()
	})
	

	//setting column widths
	list.SetColumnWidth(0, 140)
	list.SetColumnWidth(1, 140)
	list.SetColumnWidth(2, 140)
	list.SetColumnWidth(3, 140)
	list.SetColumnWidth(4, 280)

	//display list and sorting buttons
	content = container.NewBorder(
		container.NewHBox(refreshButton, sortByScoreButton, sortByNameButton, sortByUUIDButton),
		nil, nil, nil,
		list,
	)
	return
	//window.SetContent(content)
	//window.Resize(fyne.NewSize(400, 500)) // Adjust size as needed
	//window.CenterOnScreen()
}