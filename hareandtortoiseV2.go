package main

import (
	"hareandtortoise/v2/settings"
	"hareandtortoise/v2/ui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"os/exec"
	"runtime"
	
)

func openbrowser(url string, myWindow fyne.Window) {
	var err error
	switch runtime.GOOS {
	case "linux":
	  err = exec.Command("xdg-open", url).Start()
	case "windows":
	  err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
	  err = exec.Command("open", url).Start()
	default:
	  
	}
	if err != nil {
	  dialog.NewError(err, myWindow)
	}
  }

func main() {
	// Create the app
	hareandtortoise := app.New()
	mainWindow := hareandtortoise.NewWindow("Animal Simulation")

	// Run the filesystem check before loading anything else
	settings.CheckAndCreateFolderAndFile(mainWindow)

	// Toolbar setup
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
			const url string = "github.com/fossie08/hareAndTortoise-v2"
			openbrowser(url, mainWindow)
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			settings.ShowSettingsWindow(hareandtortoise)
		}),
		widget.NewToolbarAction(theme.ContentClearIcon(), func() {
			mainWindow.Close()
		}),
	)

	// Tab setup
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Leaderboard", theme.MenuIcon(), ui.DisplayLeaderboard()),
		container.NewTabItemWithIcon("Races", theme.HistoryIcon(), ui.SearchAnimals(mainWindow)),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	// Main container
	main := container.NewBorder(toolbar, nil, nil, nil, tabs)

	// Set up the window and display
	mainWindow.SetContent(main)
	mainWindow.Resize(fyne.NewSize(1000, 500))
	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()
}
