package ui

import (
	"hareandtortoise/v2/simulation"
	"strconv"

	"fyne.io/fyne/v2"

	//	"fyne.io/fyne/v2/app"
	//	"fyne.io/fyne/v2/container"

	//	"fyne.io/fyne/v2/layout"
	//	"log"

	//	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// extend the original widget
type numericalEntry struct {
	widget.Entry
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *numericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ',' {
		e.Entry.TypedRune(r)
	}
}

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

func AddAnimal(hareandtortoise fyne.App, window fyne.Window) {
	animalName := widget.NewEntry()
	animalName.SetPlaceHolder("Animal name")
	animalMinSpeed := newNumericalEntry()
	animalMinSpeed.SetPlaceHolder("Minimum speed")
	animalMaxSpeed := newNumericalEntry()
	animalMaxSpeed.SetPlaceHolder("Maximum speed")
	content := container.NewVBox(animalName, animalMinSpeed, animalMaxSpeed, widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), func() {
		simulation.CreateAnimal(animalName.Text, animalMinSpeed.Text, animalMaxSpeed.Text)
		window.Hide()
	}))	
	window.SetContent(content)
	window.Resize(fyne.NewSize(600,500))
	window.CenterOnScreen()
}