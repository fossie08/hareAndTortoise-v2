package simulation

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"time"
	"math/rand"
)

func DrawRaceTrack(myApp fyne.App, numLanes int, laneHeight int, windowWidth float32, players []Player, totalDistance int) {
	// Create main window
	mainWindow := myApp.NewWindow("Race Simulation")
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
	trackLine.Move(fyne.NewPos(0, float32(windowHeight)/2-2))

	// Create Start and Finish labels
	startText := canvas.NewText("Start", theme.ForegroundColor())
	startText.TextSize = 24
	startText.Move(fyne.NewPos(10, float32(windowHeight)/2-30))

	finishText := canvas.NewText("Finish", theme.ForegroundColor())
	finishText.TextSize = 24
	finishText.Move(fyne.NewPos(float32(windowWidth)-100, float32(windowHeight)/2-30))

	// Add elements to container
	trackContainer.Add(trackLine)
	trackContainer.Add(startText)
	trackContainer.Add(finishText)

	// Player icons representing players
	playerImages := make([]*canvas.Image, len(players))
	for i := 0; i < numLanes; i++ {
		animal := canvas.NewImageFromFile("data/image.png")
		animal.Resize(fyne.NewSize(50, 50))
		initialPos := fyne.NewPos(0, float32(laneHeight*i+laneHeight/2)-25)
		animal.Move(initialPos)
		playerImages[i] = animal
		trackContainer.Add(animal)
	}

	rand.Seed(time.Now().UnixNano())

	// Reset each player's distance and status
	for i := range players {
		players[i].Distance = 0
		players[i].Finished = false
		players[i].Place = 0
	}

	finishedPlayers := 0
	currentPlace := 1 // Tracks the finishing position

	// Race simulation loop using a goroutine
	go func() {
		for finishedPlayers < len(players) {
			for i := range players {
				if !players[i].Finished {
					// Random movement for each player within their speed range
					players[i].Distance += RandomFloat(players[i].MinSpeed, players[i].MaxSpeed)

					// Move the player image in the UI
					playerProgress := (players[i].Distance / float64(totalDistance)) * float64(windowWidth-50)
					newPos := fyne.NewPos(float32(playerProgress), float32(laneHeight*i+laneHeight/2)-25)
					playerImages[i].Move(newPos)
					canvas.Refresh(playerImages[i])

					// Check if the player has finished the race
					if players[i].Distance >= float64(totalDistance) {
						players[i].Finished = true
						players[i].Place = currentPlace
						currentPlace++
						finishedPlayers++
						playerImages[i].Move(fyne.NewPos(float32(totalDistance), float32(laneHeight*i+laneHeight/2)-25))
						canvas.Refresh(playerImages[i])
					}
				}
			}
			time.Sleep(100 * time.Millisecond) // Control speed of simulation
		}
	}()

	// Set the content of the window
	mainWindow.SetContent(trackContainer)
	mainWindow.Resize(fyne.NewSize(float32(windowWidth), float32(windowHeight)))
	mainWindow.CenterOnScreen()
	mainWindow.Show()
}
