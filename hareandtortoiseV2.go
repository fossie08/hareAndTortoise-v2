package main

import (
	"hareandtortoise/v2/ui"
	"hareandtortoise/v2/settings"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"fyne.io/fyne/v2/theme"
)


func main() {
	hareandtortoise := app.New()
	mainWindow := hareandtortoise.NewWindow("Animal Simulation")
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			ui.ShowSetupRaceMenu(hareandtortoise)
		}),
		widget.NewToolbarAction(theme.AccountIcon(), func() {
			addAnimalWindow := hareandtortoise.NewWindow("Add animal")
			ui.AddAnimal(hareandtortoise, addAnimalWindow)
			addAnimalWindow.Show()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.FileImageIcon(), func() {
			settings.ImageSelection(hareandtortoise)
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

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Leaderboard", theme.MenuIcon(), ui.DisplayLeaderboard()),
		container.NewTabItemWithIcon("Races", theme.HistoryIcon(), widget.NewLabel("test")),	
	)
	tabs.SetTabLocation(container.TabLocationTop)

	main := container.NewBorder(toolbar, nil, nil, nil, tabs)

	mainWindow.SetContent(main)
	mainWindow.Resize(fyne.NewSize(1000, 500))
	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()
}