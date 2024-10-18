package main

import (
	"fmt"
	"hareandtortoise/v2/ui"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
//	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	hareandtortoise := app.New()
	mainWindow := hareandtortoise.NewWindow("Animal Simulation")

	raceTrackButton := widget.NewButton("Open racetrack", func() {
		newWindow := hareandtortoise.NewWindow("Running Track")
		ui.DrawRaceTrack(hareandtortoise, newWindow, 5, 70, 1000.0)
		newWindow.Show() // Not ShowAndRun, because ShowAndRun is only for the main window
	})
	
	leaderboardButton := widget.NewButton("Launch leaderboard", func() {
		newWindow := hareandtortoise.NewWindow("Leaderboard")
		players, err := ui.ReadCSV("data/leaderboard.simulation")
		if err != nil {
			fmt.Println("Error loading leaderboard:", err)
			return
		}
		ui.DisplayLeaderboard(hareandtortoise, newWindow, players)
		newWindow.Show() // Again, use Show() instead of ShowAndRun
	})
	

	content := container.NewBorder(
		container.NewHBox(raceTrackButton, leaderboardButton),
		nil, nil, nil,
	)

	mainWindow.SetContent(content)
	mainWindow.Resize(fyne.NewSize(1000, 500))
	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()
}
