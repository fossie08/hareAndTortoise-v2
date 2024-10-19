package main

import (
//	"fmt"
	"hareandtortoise/v2/ui"
//	"hareandtortoise/v2/simulation"
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
	/*
	newCharacterButton := widget.NewButton("Add user", func() {
		ui.AddAnimal(hareandtortoise).Show()
	})
*/
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.AccountIcon(), func() {
			addAnimalWindow := hareandtortoise.NewWindow("Add animal")
			ui.AddAnimal(hareandtortoise, addAnimalWindow)
			addAnimalWindow.Show()
		}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	buttons := container.NewHBox(raceTrackButton)
	content := container.NewVBox(toolbar, buttons)
	main := container.NewBorder(content, nil, ui.DisplayLeaderboard(), nil, nil)
	//content := container.NewBorder(toolbar, nil, nil, nil,)

	mainWindow.SetContent(main)
	mainWindow.Resize(fyne.NewSize(1000, 500))
	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()
}
/*
func leaderboardWidget(hareandtortoise fyne.App) {
	newWindow := hareandtortoise.NewWindow("Leaderboard")
	players, err := ui.ReadCSV("data/leaderboard.simulation")
	if err != nil {
		fmt.Println("Error loading leaderboard:", err)
		return
	}
	ui.DisplayLeaderboard(hareandtortoise, newWindow, players)
	newWindow.Show()
}
*/