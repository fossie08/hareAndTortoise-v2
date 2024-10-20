package settings

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"hareandtortoise/v2/simulation"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ImageSelection(app fyne.App) {
	w := app.NewWindow("Image Selector")

	// Create variables to hold the selected image and animal name
	var previewImage *canvas.Image
	var selectedImagePath string
	var selectedAnimal string
	var selectedImage image.Image

	// Preview container
	previewLabel := widget.NewLabel("No image selected")
	previewImage = canvas.NewImageFromResource(nil)
	previewImage.FillMode = canvas.ImageFillContain
	previewImage.SetMinSize(fyne.NewSize(200, 200)) // Set a minimum size for the preview

	// Dropdown for selecting an animal
	players, err := simulation.ReadCSV("data/animal.simulation")
	if err != nil {
		dialog.ShowError(err, w)
	}
	var animalEntryOptions []string
	// map to store player UUIDs
	playerUUIDs := make(map[string]string)
	for _, player := range players {
		animalEntryOptions = append(animalEntryOptions, player.Name)
		playerUUIDs[player.Name] = player.UUID
	}

	animalEntry := widget.NewSelect(animalEntryOptions, func(value string) {
		selectedAnimal = value
	})
	animalEntry.PlaceHolder = "Select an animal"

	// File selection button
	fileBtn := widget.NewButton("Select Image", func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				// Load and display the image
				selectedImagePath = reader.URI().Path()
				file, err := os.Open(selectedImagePath)
				if err != nil {
					previewLabel.SetText("Error opening file")
					return
				}
				defer file.Close()

				// Check file extension
				fileExtension := strings.ToLower(filepath.Ext(selectedImagePath))
				if fileExtension == ".png" {
					// Handle PNG files directly (no need to decode and re-encode)
					img, err := png.Decode(file)
					if err != nil {
						previewLabel.SetText("Error decoding PNG image")
						return
					}
					selectedImage = img
				} else {
					// Try decoding as JPEG if it's not PNG
					img, err := jpeg.Decode(file)
					if err != nil {
						previewLabel.SetText("Error decoding image as JPEG")
						return
					}
					selectedImage = img
				}

				previewLabel.SetText(reader.URI().Name())
				previewImage = canvas.NewImageFromImage(selectedImage)
				previewImage.FillMode = canvas.ImageFillContain
				previewImage.SetMinSize(fyne.NewSize(200, 200))
				previewImage.Refresh()
			}
		}, w).Show()
	})

	// Import button action
	importBtn := widget.NewButton("Import", func() {
		if selectedImagePath != "" && selectedAnimal != "" {
			// Get the UUID of the selected animal
			playerUUID, ok := playerUUIDs[selectedAnimal]
			if !ok {
				dialog.ShowError(fmt.Errorf("Animal not found"), w)
				return
			}

			// Check the file extension
			fileExtension := strings.ToLower(filepath.Ext(selectedImagePath))
			newFileName := fmt.Sprintf("data/%s.png", playerUUID) // Always use .png for new files

			if fileExtension == ".png" {
				// Simply copy the PNG file without converting
				destinationFile, err := os.Create(newFileName)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				defer destinationFile.Close()

				// Copy the PNG directly
				sourceFile, err := os.Open(selectedImagePath)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				defer sourceFile.Close()

				_, err = destinationFile.ReadFrom(sourceFile)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

			} else {
				// Convert to PNG if it's not already PNG
				destinationFile, err := os.Create(newFileName)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				defer destinationFile.Close()

				// Encode and save the image as PNG
				err = png.Encode(destinationFile, selectedImage)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
			}

			// Success notification
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Import Successful",
				Content: fmt.Sprintf("Image assigned to %s (UUID: %s)", selectedAnimal, playerUUID),
			})
			w.Close()
		} else {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Import Failed",
				Content: "Please select an image and assign an animal.",
			})
		}
	})

	// Layout
	content := container.NewVBox(
		fileBtn,
		previewLabel,
		previewImage,
		animalEntry,
		importBtn,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 600))
	w.Show()
}
