package simulation

import (
	"fmt"
	"os"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"time"
	"math/rand"
	"sort"
)

// Function to display the race results in a list format
func ShowRaceResultsWindow(app fyne.App, players []Player, mainWindow fyne.Window) {
	// Create a new window for race results
	resultsWindow := app.NewWindow("Race Results")

	
	// Container for the list of results
	resultsContainer := container.NewVBox()

	// Sort players by their finishing place
	sort.Slice(players, func(i, j int) bool {
		return players[i].Place < players[j].Place
	})

	// Loop through sorted players and display their place and name
	for _, player := range players {
		if player.Finished {
			// Add player's finishing place and name to the results list
			result := fmt.Sprintf("Place: %d - %s", player.Place, player.Name)
			resultLabel := canvas.NewText(result, theme.ForegroundColor())
			resultsContainer.Add(resultLabel)
		}
	}

	// Set and display the results window content
	resultsWindow.SetContent(resultsContainer)
	resultsWindow.Resize(fyne.NewSize(300, 400)) // Adjust size as needed
	resultsWindow.CenterOnScreen()
	resultsWindow.Show()
}


// Main race function
func DrawRaceTrack(myApp fyne.App, numLanes int, laneHeight int, windowWidth float32, players []Player, totalDistance int) {
	mainWindow := myApp.NewWindow("Race Simulation")
	trackContainer := container.NewWithoutLayout()
	windowHeight := float32(numLanes) * float32(laneHeight)

	lightGreen := color.RGBA{34, 139, 34, 255}
	darkGreen := color.RGBA{0, 100, 0, 255}

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

	trackLine := canvas.NewLine(theme.ForegroundColor())
	trackLine.StrokeWidth = 5
	trackLine.Resize(fyne.NewSize(float32(windowWidth), 5))
	trackLine.Move(fyne.NewPos(0, float32(windowHeight)/2-2))

	startText := canvas.NewText("Start", theme.ForegroundColor())
	startText.TextSize = 24
	startText.Move(fyne.NewPos(10, float32(windowHeight)/2-30))

	finishText := canvas.NewText("Finish", theme.ForegroundColor())
	finishText.TextSize = 24
	finishText.Move(fyne.NewPos(float32(windowWidth)-100, float32(windowHeight)/2-30))

	trackContainer.Add(trackLine)
	trackContainer.Add(startText)
	trackContainer.Add(finishText)

	playerImages := make([]*canvas.Image, len(players))
	for i := 0; i < numLanes; i++ {
		imagePath := fmt.Sprintf("data/%s.png", players[i].UUID)

		// Check if the image exists, if not use default.png
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			fmt.Printf("Image for %s not found, using default.png\n", players[i].Name)
			imagePath = "data/default.png"
		} else {
			fmt.Printf("Using image for %s: %s\n", players[i].Name, imagePath)
		}

		animal := canvas.NewImageFromFile(imagePath)
		animal.Resize(fyne.NewSize(50, 50))
		initialPos := fyne.NewPos(0, float32(laneHeight*i+laneHeight/2)-25)
		animal.Move(initialPos)
		playerImages[i] = animal
		trackContainer.Add(animal)
	}

	rand.Seed(time.Now().UnixNano())

	for i := range players {
		players[i].Distance = 0
		players[i].Finished = false
		players[i].Place = 0
	}

	finishedPlayers := 0
	currentPlace := 1

	go func() {
		for finishedPlayers < len(players) {
			for i := range players {
				if !players[i].Finished {
					players[i].Distance += RandomFloat(players[i].MinSpeed, players[i].MaxSpeed)
					playerProgress := (players[i].Distance / float64(totalDistance)) * float64(windowWidth-50)
					if playerProgress > float64(windowWidth-50) {
						playerProgress = float64(windowWidth - 50)
					}
					newPos := fyne.NewPos(float32(playerProgress), float32(laneHeight*i+laneHeight/2)-25)
					playerImages[i].Move(newPos)
					canvas.Refresh(playerImages[i])

					if players[i].Distance >= float64(totalDistance) {
						players[i].Finished = true
						players[i].Place = currentPlace
						currentPlace++
						finishedPlayers++
						playerImages[i].Move(fyne.NewPos(float32(windowWidth-50), float32(laneHeight*i+laneHeight/2)-25))
						canvas.Refresh(playerImages[i])
					}
				}
			}
			time.Sleep(100 * time.Millisecond)
		}

		// Once race finishes, show podium window
		ShowRaceResultsWindow(myApp, players, mainWindow)
		mainWindow.Close()
	}()	

	mainWindow.SetContent(trackContainer)
	mainWindow.Resize(fyne.NewSize(float32(windowWidth), float32(windowHeight)))
	mainWindow.CenterOnScreen()
	mainWindow.Show()
}
