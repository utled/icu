package config

import (
	"encoding/json"
	"fmt"
	"igloo/data"
	"os"
	"path/filepath"
)

func InitializeConfig(homePath string, servicePath string) error {
	defaultConfig := data.Config{
		LargeSyncPath:      "/",
		QuickSyncPath:      homePath,
		LargeSyncFrequenzy: 5,
	}

	defaultConfigJSON, _ := json.MarshalIndent(defaultConfig, "", "  ")
	configFilepath := filepath.Join(servicePath, "igloo.conf")

	file, err := os.Create(configFilepath)
	if err != nil {
		return fmt.Errorf("failed to create config file:\n%v", err)
	}
	defer file.Close()

	_, err = file.WriteString(string(defaultConfigJSON))
	if err != nil {
		return fmt.Errorf("failed to write to config file\n%v", err)
	}
	file.Sync()

	return nil
}
