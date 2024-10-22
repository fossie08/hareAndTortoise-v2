package settings

import (
	"os"
	"path/filepath"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

const (
	folderName         = "data"
	fileHeader         = "Name,Score,Min Speed,Max Speed,UUID\n"
	defaultPictureName = "default.png"
)

var fileName = "animal.simulation"
var filePath = filepath.Join(folderName, fileName)
var pictureFilepath = filepath.Join(folderName, defaultPictureName)

// CheckAndCreateFolderAndFile checks if the folder and file exist, and creates them if necessary
func CheckAndCreateFolderAndFile(mainWindow fyne.Window) {
	// Check if folder exists
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		// Create the folder if it doesn't exist
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			dialog.NewError(err, mainWindow).Show()
		}
	}

	// Check if the picture file exists
	if _, err := os.Stat(pictureFilepath); os.IsNotExist(err) {
		dialog.NewError(err, mainWindow).Show()
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, err := os.Create(filePath)
		if err != nil {
			dialog.NewError(err, mainWindow).Show()
		}
		defer file.Close()

		// Write the header to the file
		_, err = file.WriteString(fileHeader)
		if err != nil {
			dialog.NewError(err, mainWindow).Show()
		}
	} else {
		dialog.NewInformation("Filesystem check", "The file and folder check has completed successfully", mainWindow).Show()
	}
}
