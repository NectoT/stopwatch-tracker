package main

import (
	"context"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func createApp(t *testing.T) App {
	tempDir := t.TempDir()
	a := App{configPath: tempDir + "/config-test.json", dataPath: tempDir + "/data-test.json"}
	a.OnStartup(context.TODO(), application.DefaultServiceOptions)
	return a
}

func TestBasicOperations(t *testing.T) {
	a := createApp(t)

	a.AddStopwatch("1", "stopwatch", 0)
	a.UpdateStopwatchTime("1", true)
	time.Sleep(time.Second)
	a.UpdateStopwatchTime("1", false)
	a.DeleteStopwatch("1")
	if len(a.GetStopwatches()) != 0 {
		t.Fatal("There must be zero stopwatches after creating and deleting a stopwatch")
	}
}
