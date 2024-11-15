package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// App struct
type App struct {
	ctx        context.Context
	userConfig UserConfig
}

type ColorTheme string

const (
	Light ColorTheme = "light"
	Dark  ColorTheme = "dark"
)

var ColorThemes = []struct {
	Value  ColorTheme
	TSName string
}{
	{Light, "Light"},
	{Dark, "Dark"},
}

type UserConfig struct {
	ColorTheme ColorTheme `json:"ColorTheme"`
}

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	return configDir + "/timer-tracker.json", err
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	configPath, err := getConfigPath()
	if err != nil {
		return
	}
	if configFile, err := os.Open(configPath); errors.Is(err, os.ErrNotExist) {
		configFile, err := os.Create(configPath)
		if err != nil {
			panic(err)
		}
		a.userConfig = UserConfig{
			ColorTheme: Light,
		}
		json.NewEncoder(configFile).Encode(a.userConfig)
	} else {
		json.NewDecoder(configFile).Decode(&a.userConfig)
	}
}

func (a *App) GetUserColorTheme() ColorTheme {
	return a.userConfig.ColorTheme
}

func (a *App) SetUserColorTheme(theme ColorTheme) {
	a.userConfig.ColorTheme = theme
	configPath, err := getConfigPath()
	if err != nil {
		fmt.Println("Could not get a config path")
	}

	configFile, err := os.Create(configPath)
	if err != nil {
		fmt.Println("Could not truncate and open a config file")
	}
	json.NewEncoder(configFile).Encode(a.userConfig)
}
