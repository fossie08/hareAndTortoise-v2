package settings
// import some stuff
import (
	"os"
	"path/filepath"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	folderName         = "data"
	fileHeader         = "Name,Score,Min Speed,Max Speed,UUID\n"
	defaultPictureName = "default.png"
	defaultSoundName   = "cheering.mp3"
)

var fileName = "animal.simulation"
var filePath = filepath.Join(folderName, fileName)
var pictureFilepath = filepath.Join(folderName, defaultPictureName)
var soundFilepath = filepath.Join(folderName, defaultSoundName)

// custom error function as the old error function didn't scale right
func showCustomError(err error, mainWindow fyne.Window) {
    content := container.NewVBox(
        widget.NewLabel(err.Error()),
    )
    
    customDialog := dialog.NewCustom("Error", "OK", content, mainWindow)
    customDialog.Resize(fyne.NewSize(200, 200))
    customDialog.Show()
}

// CheckAndCreateFolderAndFile checks if the folder and file exist, and creates them if necessary
func CheckAndCreateFolderAndFile(mainWindow fyne.Window) {
	allChecksPassed := true // Track whether all checks pass

	// Check if folder exists
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		// Create the folder if it doesn't exist
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			showCustomError(err, mainWindow)
			allChecksPassed = false
		}
	}
	
	// Check if the sound file exists
	if _, err := os.Stat(soundFilepath); os.IsNotExist(err) {
		showCustomError(err, mainWindow)
		allChecksPassed = false
	}
	
	// Check if the picture file exists
	if _, err := os.Stat(pictureFilepath); os.IsNotExist(err) {
		showCustomError(err, mainWindow)
		allChecksPassed = false
	}
	
	// Check if the main file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, err := os.Create(filePath)
		if err != nil {
			showCustomError(err, mainWindow)
			allChecksPassed = false
		} else {
			defer file.Close()
	
			// Write the header to the file
			_, err = file.WriteString(fileHeader)
			if err != nil {
				showCustomError(err, mainWindow)
				allChecksPassed = false
			}
		}
	}
	
	// Show success message if all checks passed
	if allChecksPassed {
		dialog.NewInformation("Filesystem check", "The file and folder check has completed successfully", mainWindow).Show()
	}
	
}
