package ui

import (
    "bufio"
    "encoding/csv"
    "errors"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    "os"
    "path/filepath"
    "strconv"
	"fmt"
    "strings"
)

type Race struct {
    UUID               string
    Place              int
    DistanceTravelled  float64
    Score              int
    TotalDistance      float64
    Rounds             int
    Date               string
    Time               string
}

type Animal struct {
    Name string
    UUID string
}

type AnimalInsights struct {
    TotalScore         int
    RacesParticipated  int
    BestPlace          int
    RaceData           []Race
    Last10Positions    []int
}

// ReadRaceFiles reads all .simulation files in the data/ directory.
func ReadRaceFiles() (map[string][]Race, error) {
    raceMap := make(map[string][]Race)

    err := filepath.Walk("data/", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Ignore .png files and animal.simulation
        if strings.HasSuffix(info.Name(), ".png") || info.Name() == "animal.simulation" {
            return nil
        }

        if strings.HasSuffix(info.Name(), ".simulation") {
            races, err := parseRaceFile(path)
            if err != nil {
                return err
            }
            raceUUID := strings.TrimSuffix(info.Name(), ".simulation")
            raceMap[raceUUID] = races
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    return raceMap, nil
}

// parseRaceFile parses a .simulation file and returns race data.
func parseRaceFile(filename string) ([]Race, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(bufio.NewReader(file))
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    if len(records) < 2 {
        return nil, errors.New("no race data found")
    }

    var races []Race
    generalData := records[0] // General data is in the first row

    for _, record := range records[1:] {
        if len(record) != 8 { // Ensure the record has the correct number of fields
            continue // Skip malformed records
        }

        place, _ := strconv.Atoi(record[1])
        distanceTravelled, _ := strconv.ParseFloat(record[2], 64)
        score, _ := strconv.Atoi(record[3])
        totalDistance, _ := strconv.ParseFloat(record[4], 64)
        rounds, _ := strconv.Atoi(record[5])

        races = append(races, Race{
            UUID:              record[0],
            Place:             place,
            DistanceTravelled: distanceTravelled,
            Score:             score,
            TotalDistance:     totalDistance,
            Rounds:            rounds,
            Date:              generalData[6], // Extracting date from the first row
            Time:              generalData[7], // Extracting time from the first row
        })
    }

    return races, nil
}

// ReadAnimalData reads the animal data from animal.simulation.
func ReadAnimalData() (map[string]Animal, error) {
    animalMap := make(map[string]Animal)

    file, err := os.Open("data/animal.simulation")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(bufio.NewReader(file))
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    for _, record := range records[1:] {
        if len(record) < 2 {
            continue // Skip malformed records
        }
        animalUUID := record[4]
        animalName := record[0]
        animalMap[animalUUID] = Animal{
            Name: animalName,
            UUID: animalUUID,
        }
    }

    return animalMap, nil
}

// SearchAnimalInsights provides insights for a specific animal UUID or name.
func SearchAnimalInsights(raceData map[string][]Race, animalID string, animalMap map[string]Animal) (AnimalInsights, error) {
    var insights AnimalInsights
    var foundAnimal bool

    for _, races := range raceData {
        for _, race := range races {
            if race.UUID == animalID {
                foundAnimal = true
                insights.TotalScore += race.Score
                insights.RacesParticipated++
                insights.Last10Positions = append(insights.Last10Positions, race.Place)
                if insights.BestPlace == 0 || race.Place < insights.BestPlace {
                    insights.BestPlace = race.Place
                }
                insights.RaceData = append(insights.RaceData, race)
            }
        }
    }

    if !foundAnimal {
        // If searching by name, look up the UUID first
        for _, animal := range animalMap {
            if animal.Name == animalID {
                return SearchAnimalInsights(raceData, animal.UUID, animalMap)
            }
        }
        return insights, errors.New("animal not found")
    }

    return insights, nil
}

func SearchAnimals(myWindow fyne.Window) *fyne.Container {

    // Create a search entry
    searchEntry := widget.NewEntry()
    searchEntry.SetPlaceHolder("Enter Animal Name...")

    // Create a label for results
    resultsLabel := widget.NewLabel("Results will appear here.")

    // Create a button to trigger the search
    searchButton := widget.NewButton("Search", func() {
        animalID := searchEntry.Text
        animalMap, err := ReadAnimalData()
        if err != nil {
            dialog.ShowError(err, myWindow)
            return
        }
		uuid, err := GetAnimalUUID(animalMap, animalID)
		if err != nil {
			dialog.ShowError(err, myWindow)
		} else {
			raceData, err := ReadRaceFiles()
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
	
			insights, err := SearchAnimalInsights(raceData, uuid, animalMap)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
	
			// Show results in the same window
			results := fmt.Sprintf("Total Score: %d\nRaces Participated: %d\nBest Place: %d\nLast 10 Positions: %v", 
				insights.TotalScore, insights.RacesParticipated, insights.BestPlace, insights.Last10Positions)
			resultsLabel.SetText(results)
		}
    })

    // Layout the search bar, button, and results label
    content := container.NewVBox(searchEntry, searchButton, resultsLabel)
	return content
}

// GetAnimalUUID returns the UUID for a given animal name from the animal map.
func GetAnimalUUID(animalMap map[string]Animal, animalName string) (string, error) {
    for _, animal := range animalMap {
        if animal.Name == animalName {
            return animal.UUID, nil
        }
    }
    return "", errors.New("animal not found")
}
