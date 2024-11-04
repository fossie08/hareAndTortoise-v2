package ui
// import some stuff
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
// race data structure
type Race struct {
    UUID               string
    Place              int
    DistanceTravelled  float64
    Score              int
    TotalDistance      float64
    Rounds             int
    Date               string
    Time               string
    Name               string
}
// animal data strucutre
type Animal struct {
    Name string
    UUID string
}
// animal insights data structure
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
        fmt.Printf("Discovered file: %s\n", path)
    
        if strings.HasSuffix(info.Name(), ".png") || info.Name() == "animal.simulation" {
            return nil
        }
    
        if strings.HasSuffix(info.Name(), ".simulation") {
            races, err := parseRaceFile(path)
            if err != nil {
                fmt.Printf("Error parsing file %s: %v\n", path, err)
                return err
            }
            raceUUID := strings.TrimSuffix(info.Name(), ".simulation")
            fmt.Printf("Race UUID: %s, Races: %+v\n", raceUUID, races)
            raceMap[raceUUID] = races
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return raceMap, nil
}

// parseRaceFile parses a .simulation file and returns race data
func parseRaceFile(filename string) ([]Race, error) {
    fmt.Printf("Opening file: %s\n", filename)
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error opening file %s: %v\n", filename, err)
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(bufio.NewReader(file))
    records, err := reader.ReadAll()
    if err != nil {
        fmt.Printf("Error reading CSV data from %s: %v\n", filename, err)
        return nil, err
    }

    fmt.Printf("Records in %s: %+v\n", filename, records)
    if len(records) < 2 {
        fmt.Printf("No race data found in %s\n", filename)
        return nil, errors.New("no race data found")
    }

    var races []Race
    for _, record := range records[1:] {
        if len(record) != 9 { // Adjusted to expect 9 fields including Name
            fmt.Printf("Skipping malformed record in %s: %+v\n", filename, record)
            continue
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
            Date:              record[6],
            Time:              record[7],
            Name:              record[8], // Add Name field here if needed in Race struct
        })
    }

    fmt.Printf("Parsed races from %s: %+v\n", filename, races)
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
            fmt.Printf("Checking race UUID: %s against animal ID: %s\n", race.UUID, animalID)
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
                fmt.Printf("Found animal name %s with UUID %s\n", animalID, animal.UUID)
                return SearchAnimalInsights(raceData, animal.UUID, animalMap)
            }
        }
        return insights, errors.New("animal not found 2")
    }
    return insights, nil
}

// GetAnimalUUID returns the UUID for a given animal name from the animal map.
func GetAnimalUUID(animalMap map[string]Animal, animalName string) (string, error) {
    for _, animal := range animalMap {
        fmt.Printf("Comparing '%s' with '%s'\n", animal.Name, animalName)
        if strings.TrimSpace(animal.Name) == strings.TrimSpace(animalName) {
            return animal.UUID, nil
        }
    }
    
    return "", errors.New("animal not found")
}

// SearchAnimals sets up the search UI with the chart displayed below results.
func SearchAnimals(myWindow fyne.Window) *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Enter Animal Name...")
	resultsLabel := widget.NewLabel("Results will appear here.")
	searchContainer := container.NewVBox()

	searchButton := widget.NewButton("Search", func() {
		animalID := searchEntry.Text
		animalMap, _ := ReadAnimalData()
		uuid, err := GetAnimalUUID(animalMap, animalID)
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}

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

		results := fmt.Sprintf("Total Score: %d\nRaces Participated: %d\nBest Place: %d\nLast 10 Positions: %v", 
			insights.TotalScore, insights.RacesParticipated, insights.BestPlace, insights.Last10Positions)
		resultsLabel.SetText(results)
	})
    
    searchContainer.Add(searchEntry)
	searchContainer.Add(searchButton)
    searchContainer.Add(resultsLabel)
	return searchContainer
}