package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ShowSettingsWindow creates and shows the settings window for theme selection
func ShowSettingsWindow(app fyne.App) {
	settingsWindow := app.NewWindow("Settings")

	// Store the current theme setting
	currentTheme := "Dark"
	if app.Settings().Theme() == theme.LightTheme() {
		currentTheme = "Light"
	}

	// Theme selection (doesn't apply until Apply is clicked)
	themeLabel := widget.NewLabel("Select Theme:")
	themeSelect := widget.NewSelect([]string{"Light", "Dark"}, nil)
	themeSelect.SetSelected(currentTheme)

	// Apply button to apply the selected theme
	applyButton := widget.NewButton("Apply", func() {
		selectedTheme := themeSelect.Selected
		if selectedTheme == "Light" {
			app.Settings().SetTheme(theme.LightTheme())
		} else {
			app.Settings().SetTheme(theme.DarkTheme())
		}
		settingsWindow.Close()
	})

	// Layout the UI components
	content := container.NewVBox(
		themeLabel,
		themeSelect,
		applyButton,
	)

	// Show the window
	settingsWindow.SetContent(content)
	settingsWindow.Resize(fyne.NewSize(300, 150))
	settingsWindow.CenterOnScreen()
	settingsWindow.Show()
}
