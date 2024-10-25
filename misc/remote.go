package misc

import (
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "os"
	"encoding/json"
    "strings"
)

type Settings struct {
	RemoteURL      string `json:"remote_url"`
	RemoteUsername string `json:"remote_username"`
	RemotePassword string `json:"remote_password"`
}

const settingsFilePath = "data/settings.json"

// loadSettings loads the settings from a JSON file
func loadSettings() (Settings, error) {
	var settings Settings
	file, err := os.Open(settingsFilePath)
	if err != nil {
		return settings, err // Return zero value if file doesn't exist or is inaccessible
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	return settings, err
}

func loginAndGetSession(username, password string, remoteurl string) (*http.Client, error) {
    // Create an HTTP client with a cookie jar to store session cookies
    jar, _ := cookiejar.New(nil)
    client := &http.Client{
        Jar: jar,
    }

    data := url.Values{}
    data.Set("username", username)
    data.Set("password", password)
	loginurl := remoteurl + "/login"
    req, err := http.NewRequest("POST", loginurl, strings.NewReader(data.Encode()))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Check for a successful login (e.g., HTTP 200 or redirect status)
    if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("login failed with status code: %d", resp.StatusCode)
    }

    return client, nil
}

func uploadFile(client *http.Client, filePath string, username string, remoteurl string) error {
    file, err := os.Open("data/"+filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)
    writer.WriteField("username", username)
    fileWriter, _ := writer.CreateFormFile("file", filePath)
    io.Copy(fileWriter, file)
    writer.Close()
	uploadurl := remoteurl + "/upload"
    req, err := http.NewRequest("POST", uploadurl, &requestBody)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // Send the request using the client with the session cookie
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusOK {
        return fmt.Errorf("file upload failed with status code: %d", resp.StatusCode)
    }

    return nil
}

func UploadToRemote(simulationfile string) {
	// Load existing settings
	existingSettings, err := loadSettings()
	if err != nil {
		existingSettings = Settings{} // Use zero values if loading fails
	}

    client, err := loginAndGetSession(existingSettings.RemoteUsername, existingSettings.RemotePassword, existingSettings.RemoteURL)
    if err != nil {
        fmt.Println("Error logging in:", err)
        return
    }

    err = uploadFile(client, simulationfile, existingSettings.RemoteUsername, existingSettings.RemoteURL)
    if err != nil {
        fmt.Println("Error uploading file:", err)
    }
}
