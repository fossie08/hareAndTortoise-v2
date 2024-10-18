package ui

import (
//	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

func DrawRaceTrack(myApp fyne.App, mainWindow fyne.Window, numLanes int, laneHeight int, windowWidth float32) {
	// Create custom track container
	trackContainer := container.NewWithoutLayout()

	windowHeight := float32(numLanes) * float32(laneHeight)

	// Alternate shades of green for lanes
	lightGreen := color.RGBA{34, 139, 34, 255} // Lighter green
	darkGreen := color.RGBA{0, 100, 0, 255}    // Darker green

	// Draw lanes
	for i := 0; i < numLanes; i++ {
		laneColor := lightGreen
		if i%2 == 1 {
			laneColor = darkGreen
		}
		lane := canvas.NewRectangle(laneColor)
		lane.Resize(fyne.NewSize(float32(windowWidth), float32(laneHeight)))
		lane.Move(fyne.NewPos(0, float32(laneHeight)*float32(i)))
		trackContainer.Add(lane)
	}

	// Create track line at the middle of the lanes
	trackLine := canvas.NewLine(theme.ForegroundColor())
	trackLine.StrokeWidth = 5
	trackLine.Resize(fyne.NewSize(float32(windowWidth), 5))
	trackLine.Move(fyne.NewPos(0, float32(windowHeight)/2-2)) // Cast to float32

	// Create Start and Finish labels
	startText := canvas.NewText("Start", theme.ForegroundColor())
	startText.TextSize = 24
	startText.Move(fyne.NewPos(10, float32(windowHeight)/2-30)) // Cast to float32

	finishText := canvas.NewText("Finish", theme.ForegroundColor())
	finishText.TextSize = 24
	finishText.Move(fyne.NewPos(float32(windowWidth)-100, float32(windowHeight)/2-30)) // Cast to float32

	// Add elements to container
	trackContainer.Add(trackLine)
	trackContainer.Add(startText)
	trackContainer.Add(finishText)

	// Set the content of the window
	mainWindow.SetContent(trackContainer)
	mainWindow.Resize(fyne.NewSize(float32(windowWidth), float32(windowHeight))) // Fixed window size
	mainWindow.CenterOnScreen()
}

