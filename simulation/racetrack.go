package simulation

import (
	"fmt"
	"os"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/dialog"
	"image/color"
	"time"
	"math/rand"
	"sort"
	"github.com/google/uuid"

)



var roundNumber int = 1
var raceRunning bool = true

// Calculate scores based on race position, race distance, and number of players
func CalculateScores(players []Player, totalDistance int) {
	numPlayers := len(players)
	basePoints := func(place int) int {
		// Base points are higher for smaller races and decrease by position
		return (numPlayers - place + 1) * 10 // Example: 1st place gets 10 * numPlayers points
	}
	for i, player := range players {
		if player.Finished {
			distanceBonus := float64(player.Distance) / float64(totalDistance) // Bonus based on distance
			speedBonus := (player.MaxSpeed - player.MinSpeed) / 2               // Speed difference bonus
			players[i].Score = basePoints(player.Place) + int(distanceBonus*100) + int(speedBonus*10)
		}
	}
}


// Modify ShowRaceResultsWindow to include a "Save Race" button
func ShowRaceResultsWindow(app fyne.App, players []Player, mainWindow fyne.Window, totalDistance, numRounds int) {
	resultsWindow := app.NewWindow("Race Results")
	resultsContainer := container.NewVBox()

	sort.Slice(players, func(i, j int) bool {
		return players[i].Place < players[j].Place
	})

	for i, player := range players {
		if player.Finished {
			result := fmt.Sprintf("Place: %d - %s - Score: %d", player.Place, player.Name, players[i].Score)
			resultLabel := canvas.NewText(result, theme.ForegroundColor())
			resultsContainer.Add(resultLabel)
		}
	}

	// Add "Save Race" button
	saveButton := widget.NewButton("Save Race", func() {
		raceUUID := uuid.New().String()
		SaveRaceResults(players, totalDistance, numRounds, raceUUID)
		dialog.ShowInformation("Race Saved", "Race results have been saved successfully.", resultsWindow)
	})
	resultsContainer.Add(saveButton)

	resultsWindow.SetContent(resultsContainer)
	resultsWindow.Resize(fyne.NewSize(300, 400))
	resultsWindow.CenterOnScreen()
	resultsWindow.Show()
}

// Main race function
func DrawRaceTrack(myApp fyne.App, numLanes int, laneHeight int, windowWidth float32, players []Player, totalDistance int) {
	mainWindow := myApp.NewWindow("Race Simulation")
	trackContainer := container.NewWithoutLayout()
	windowHeight := float32(numLanes) * float32(laneHeight)

	// Display round number
	roundText := canvas.NewText(fmt.Sprintf("Round: %d", roundNumber), theme.ForegroundColor())
	roundText.TextSize = 24
	roundText.Move(fyne.NewPos(windowWidth/2-50, 10))

	lightGreen := color.RGBA{34, 139, 34, 255}
	darkGreen := color.RGBA{0, 100, 0, 255}

	playerProgressTexts := make([]*canvas.Text, len(players)) // Create array for progress text

	for i := 0; i < numLanes; i++ {
		laneColor := lightGreen
		if i%2 == 1 {
			laneColor = darkGreen
		}
		lane := canvas.NewRectangle(laneColor)
		lane.Resize(fyne.NewSize(float32(windowWidth), float32(laneHeight)))
		lane.Move(fyne.NewPos(0, float32(laneHeight)*float32(i)))
		trackContainer.Add(lane)

		// Display player names and distance travelled at the beginning of lanes
		playerNameText := canvas.NewText(players[i].Name, theme.ForegroundColor())
		playerNameText.TextSize = 18
		playerNameText.Move(fyne.NewPos(10, float32(laneHeight*i)+5))
		trackContainer.Add(playerNameText)

		// Distance text
		progressText := canvas.NewText(fmt.Sprintf("0.0/%d", totalDistance), theme.ForegroundColor())
		progressText.TextSize = 18
		progressText.Move(fyne.NewPos(150, float32(laneHeight*i)+5))
		playerProgressTexts[i] = progressText
		trackContainer.Add(progressText)
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

	// Add start, stop, and end buttons
	startButton := widget.NewButton("Start Race", func() {
		raceRunning = true
	})
	stopButton := widget.NewButton("Pause Race", func() {
		raceRunning = false
	})
	endButton := widget.NewButton("End Race", func() {
		dialog.NewConfirm("Are you sure?", "Are you sure you want to end the race?", 
		func(confirmed bool) {
			if confirmed {
				finishedPlayers = len(players)
				mainWindow.Close()
			}
		}, mainWindow).Show()
	})

	buttonContainer := container.NewHBox(startButton, stopButton, endButton, roundText)
	layout := container.NewVBox(buttonContainer, trackContainer)

	go func() {
		for finishedPlayers < len(players) {
			if raceRunning {
				roundNumber++ // Increment round number
				roundText.Text = fmt.Sprintf("Round: %d", roundNumber) // Update round number display
				canvas.Refresh(roundText)

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

						// Update distance travelled text
						playerProgressTexts[i].Text = fmt.Sprintf("%.1f/%d", players[i].Distance, totalDistance)
						canvas.Refresh(playerProgressTexts[i])

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
			}
			time.Sleep(100 * time.Millisecond)
		}

		CalculateScores(players, totalDistance)
		ShowRaceResultsWindow(myApp, players, mainWindow, totalDistance, roundNumber)
		mainWindow.Close()
	}()

	mainWindow.SetContent(layout)
	mainWindow.Resize(fyne.NewSize(float32(windowWidth), float32(windowHeight+100))) // Adjust window size
	mainWindow.CenterOnScreen()
	mainWindow.Show()
}
