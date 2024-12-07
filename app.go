package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// App struct
type App struct {
	// Path to the file that stores config information
	configPath string

	// Path to the file that stores data information
	dataPath string

	userConfig UserConfig
	userData   UserData
}

type ColorTheme string

const (
	Light ColorTheme = "light"
	Dark  ColorTheme = "dark"
)

type UserConfig struct {
	ColorTheme ColorTheme `json:"ColorTheme"`
}

type StopwatchData struct {
	IsActive        bool
	Hue             float32
	Name            string
	LastUpdated     time.Time
	CreatedAt       time.Time
	TimeAccumulated int

}

type UserData map[string]StopwatchData

func createJson[K any](jsonPath string, defaultJson K) {
	configFile, err := os.Create(jsonPath)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	json.NewEncoder(configFile).Encode(defaultJson)
}

func getOrCreateJson[K any](jsonPath string, defaultJson K) K {
	if configFile, err := os.Open(jsonPath); errors.Is(err, os.ErrNotExist) {
		createJson(jsonPath, defaultJson)
		return defaultJson
	} else {
		var jsonStruct K
		err := json.NewDecoder(configFile).Decode(&jsonStruct)
		if err != nil {
			createJson(jsonPath, defaultJson)
			return defaultJson
		}
		return jsonStruct
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) OnStartup(ctx context.Context, options application.ServiceOptions) error {
	a.userConfig = getOrCreateJson(a.configPath, UserConfig{ColorTheme: Light})

	a.userData = getOrCreateJson(a.dataPath, UserData{})

	return nil
}

func (a *App) updateUserConfigFile() {
	configFile, err := os.Create(a.configPath)
	if err != nil {
		fmt.Println("Could not truncate and open a config file")
		return
	}
	defer configFile.Close()
	json.NewEncoder(configFile).Encode(a.userConfig)
}

func (a *App) updateUserDataFile() {
	userFile, err := os.Create(a.dataPath)
	if err != nil {
		fmt.Println("Could not truncate and open a data file")
		return
	}
	defer userFile.Close()
	json.NewEncoder(userFile).Encode(a.userData)
}

func (a *App) GetUserColorTheme() ColorTheme {
	return a.userConfig.ColorTheme
}

func (a *App) SetUserColorTheme(theme ColorTheme) {
	a.userConfig.ColorTheme = theme
	a.updateUserConfigFile()
}

func (a *App) GetStopwatchIds() []string {
	ids := make([]string, 0)
	datas := make([]StopwatchData, 0)

	for key, value := range a.userData {
		inserted := false
		for i := 0; i < len(ids); i++ {
			if value.CreatedAt.Before(datas[i].CreatedAt) {
				ids = append(ids[:i+1], ids[i:]...)
				ids[i] = key
				datas = append(datas[:i+1], datas[i:]...)
				datas[i] = value
				inserted = true
				break
			}
		}
		if !inserted {
			ids = append(ids, key)
			datas = append(datas, value)
		}
	}

	return ids
}

func (a *App) GetStopwatches() map[string]StopwatchData {
	return a.userData
}

func (a *App) HasStopwatch(id string) bool {
	_, has := a.userData[id]
	return has
}

func (a *App) GetStopwatch(id string) StopwatchData {
	return a.userData[id]
}

func (a *App) AddStopwatch(id string, name string, hue float32) (bool, error) {
	if _, exists := a.userData[id]; exists {
		return false, errors.New("stopwatch with such id already exists")
	}
	a.userData[id] = StopwatchData{
		Name:      name,
		Hue:       hue,
		CreatedAt: time.Now(),
	}
	a.updateUserConfigFile()
	return true, nil
}

func (a *App) DeleteStopwatch(id string) (bool, error) {
	if !a.HasStopwatch(id) {
		fmt.Printf("stopwatch with such id %s does not exist\n", id)
		return false, errors.New("stopwatch with such id does not exist")
	}
	delete(a.userData, id)
	a.updateUserDataFile()
	return true, nil
}

func (a *App) UpdateStopwatchTime(id string, timerActive bool, timeAccumulated int) {
	data := a.userData[id]
	data.TimeAccumulated = timeAccumulated
	data.IsActive = timerActive
	data.LastUpdated = time.Now()
	a.userData[id] = data
	a.updateUserDataFile()
}

func (a *App) UpdateStopwatchName(id string, name string) {
	data := a.userData[id]
	data.Name = name
	a.userData[id] = data
	a.updateUserDataFile()
}
