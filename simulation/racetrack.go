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



var raceRunning bool = true

// score calculation - revered positions last gets 1
func CalculateScores(players []Player, totalDistance int) {
	numPlayers := len(players)
	for i, player := range players {
		if player.Finished {
			players[i].Score = (numPlayers - player.Place)+1
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

func DrawRaceTrack(myApp fyne.App, numLanes int, laneHeight int, windowWidth float32, players []Player, totalDistance int) {
    mainWindow := myApp.NewWindow("Race Simulation")
	var roundNumber int = 1
    trackContainer := container.NewWithoutLayout()
    windowHeight := float32(numLanes) * float32(laneHeight)

    // Initialize endurance for each player (endurance starts full)
    for i := range players {
        players[i].Endurance = 100 // Example starting endurance, could be different depending on player
        players[i].Resting = false // Not resting at start
    }

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
				raceRunning = false
                mainWindow.Close()
            }
        }, mainWindow).Show()
    })

    buttonContainer := container.NewHBox(startButton, stopButton, endButton, roundText)
    layout := container.NewVBox(buttonContainer, trackContainer)

    go func() {
		raceRunning = true
        for finishedPlayers < len(players) {
            if raceRunning {
                roundNumber++ // Increment round number
                roundText.Text = fmt.Sprintf("Round: %d", roundNumber) // Update round number display
                canvas.Refresh(roundText)

                for i := range players {
                    if players[i].Finished {
                        continue // Skip finished players
                    }

                    if players[i].Resting {
                        // Recover endurance and skip this round
                        players[i].Endurance += 3 * players[i].MinSpeed
                        players[i].Resting = false
                        continue
                    }

                    // Deduct endurance based on the distance run this round
                    distanceRun := RandomFloat(players[i].MinSpeed, players[i].MaxSpeed)
                    players[i].Endurance -= distanceRun

                    if players[i].Endurance <= 0 {
                        players[i].Endurance = 0
                        players[i].Resting = true
                        continue
                    }

                    // Move player if not resting
                    players[i].Distance += distanceRun
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
