package main

import (
	"fmt"
	"hareandtortoise/v2/ui"
	"hareandtortoise/v2/simulation"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
//	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"fyne.io/fyne/v2/theme"
)

func main() {
	hareandtortoise := app.New()
	mainWindow := hareandtortoise.NewWindow("Animal Simulation")

	raceTrackButton := widget.NewButton("Open racetrack", func() {
		newWindow := hareandtortoise.NewWindow("Running Track")
		ui.DrawRaceTrack(hareandtortoise, newWindow, 10, 70, 1200.0)
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
	
	newCharacterButton := widget.NewButton("Add user", func() {
		simulation.CreateAnimal()
	})

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	buttons := container.NewHBox(raceTrackButton, leaderboardButton, newCharacterButton)
	content := container.NewVBox(toolbar, buttons)
	//content := container.NewBorder(toolbar, nil, nil, nil,)

	mainWindow.SetContent(content)
	mainWindow.Resize(fyne.NewSize(1000, 500))
	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()
}
