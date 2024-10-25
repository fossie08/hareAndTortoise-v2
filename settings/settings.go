package settings

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"os"
)

// Settings holds the remote configuration
type Settings struct {
	RemoteURL      string `json:"remote_url"`
	RemoteUsername string `json:"remote_username"`
	RemotePassword string `json:"remote_password"`
}

// settingsFilePath defines where the settings will be saved
const settingsFilePath = "data/settings.json"

// loadSettings loads the settings from a JSON file
func loadSettings() (Settings, error) {
	var settings Settings
	file, err := os.Open(settingsFilePath)
	if err != nil {
		return settings, err // Return zero value if file doesn't exist or is inaccessible
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	return settings, err
}

// saveSettings saves the settings to a JSON file
func saveSettings(settings Settings) error {
	file, err := os.Create(settingsFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(settings)
}

// ShowSettingsWindow creates and shows the settings window for theme selection and remote settings
func ShowSettingsWindow(app fyne.App) {
	settingsWindow := app.NewWindow("Settings")

	// Store the current theme setting
	currentTheme := "Dark"
	if app.Settings().Theme() == theme.LightTheme() {
		currentTheme = "Light"
	}

	// Load existing settings
	existingSettings, err := loadSettings()
	if err != nil {
		existingSettings = Settings{} // Use zero values if loading fails
	}

	// Theme selection (doesn't apply until Apply is clicked)
	themeLabel := widget.NewLabel("Select Theme:")
	themeSelect := widget.NewSelect([]string{"Light", "Dark"}, nil)
	themeSelect.SetSelected(currentTheme)

	// Remote settings
	remoteURLLabel := widget.NewLabel("Remote URL:")
	remoteURLEntry := widget.NewEntry()
	remoteURLEntry.SetText(existingSettings.RemoteURL)

	remoteUsernameLabel := widget.NewLabel("Remote Username:")
	remoteUsernameEntry := widget.NewEntry()
	remoteUsernameEntry.SetText(existingSettings.RemoteUsername)

	remotePasswordLabel := widget.NewLabel("Remote Password:")
	remotePasswordEntry := widget.NewPasswordEntry()
	remotePasswordEntry.SetText(existingSettings.RemotePassword)

	// Apply button to apply the selected theme
	applyButton := widget.NewButton("Apply Theme", func() {
		selectedTheme := themeSelect.Selected
		if selectedTheme == "Light" {
			app.Settings().SetTheme(theme.LightTheme())
		} else {
			app.Settings().SetTheme(theme.DarkTheme())
		}
	})

	// Save button to save remote settings
	saveButton := widget.NewButton("Save", func() {
		remoteURL := remoteURLEntry.Text
		remoteUsername := remoteUsernameEntry.Text
		remotePassword := remotePasswordEntry.Text

		settings := Settings{
			RemoteURL:      remoteURL,
			RemoteUsername: remoteUsername,
			RemotePassword: remotePassword,
		}

		err := saveSettings(settings)
		if err != nil {
			println("Error saving settings:", err.Error())
		}
	})

	// Layout the UI components
	content := container.NewVBox(
		themeLabel,
		themeSelect,
		applyButton,
		remoteURLLabel,
		remoteURLEntry,
		remoteUsernameLabel,
		remoteUsernameEntry,
		remotePasswordLabel,
		remotePasswordEntry,
		saveButton,
	)

	// Show the window
	settingsWindow.SetContent(content)
	settingsWindow.Resize(fyne.NewSize(300, 250))
	settingsWindow.CenterOnScreen()
	settingsWindow.Show()
}
