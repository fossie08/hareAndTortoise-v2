package main

import (
//	"fmt"
	"hareandtortoise/v2/ui"
	"hareandtortoise/v2/settings"
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

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			log.Println("New race")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.AccountIcon(), func() {
			addAnimalWindow := hareandtortoise.NewWindow("Add animal")
			ui.AddAnimal(hareandtortoise, addAnimalWindow)
			addAnimalWindow.Show()
		}),
		widget.NewToolbarAction(theme.WarningIcon(), func() {
			newWindow := hareandtortoise.NewWindow("Running Track")
			ui.DrawRaceTrack(hareandtortoise, newWindow, 10, 70, 1200.0)
			newWindow.Show()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			settings.ShowSettingsWindow(hareandtortoise)
		}),
		widget.NewToolbarAction(theme.ContentClearIcon(), func() {
			mainWindow.Close()
		}),
	)

	main := container.NewBorder(toolbar, nil, nil, nil, ui.DisplayLeaderboard())

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