package ui
//import some libraries
import (
	"hareandtortoise/v2/simulation"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fmt"
)

// creating a custom text box for number only entry
type numericalEntry struct {
	widget.Entry
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}
// number only and .
func (e *numericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' {
		e.Entry.TypedRune(r)
	}
}
// prevents copy and paste in
func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}
//add animal ui
func AddAnimal(hareandtortoise fyne.App, window fyne.Window) {
	animalName := widget.NewEntry()
	animalName.SetPlaceHolder("Animal name")
	animalMinSpeed := newNumericalEntry()
	animalMinSpeed.SetPlaceHolder("Minimum speed")
	animalMaxSpeed := newNumericalEntry()
	animalMaxSpeed.SetPlaceHolder("Maximum speed")
	
	content := container.NewVBox(animalName, animalMinSpeed, animalMaxSpeed, widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), func() {
		minSpeed, _ := strconv.ParseFloat(animalMinSpeed.Text, 64)
		maxSpeed, _ := strconv.ParseFloat(animalMaxSpeed.Text, 64)
		
		// Check if minSpeed is greater than maxSpeed and swap if necessary
		if minSpeed > maxSpeed {
			minSpeed, maxSpeed = maxSpeed, minSpeed
		}
		// check if min or max speed is less than or equal to 0 
		if minSpeed <= 0 {
			minSpeed = 1
			dialog.NewError(fmt.Errorf("minimum speed cannot be 0 or below"), window)
		}

		if maxSpeed <= 0 {
			maxSpeed = 1
			dialog.NewError(fmt.Errorf("maximum speed cannot be 0 or below"), window)
		}

		// Convert back to string for saving
		simulation.CreateAnimal(animalName.Text, strconv.FormatFloat(minSpeed, 'f', -1, 64), strconv.FormatFloat(maxSpeed, 'f', -1, 64))
		window.Hide()
	}))	
	window.SetContent(content)
	window.Resize(fyne.NewSize(300, 250))
	window.CenterOnScreen()
}
