package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"igloo/data"
	"igloo/db"
)

func initializeConfig(homePath string, servicePath string) error {
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

func Main() error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to identify user home directory:%v", err)
	}

	servicePath := filepath.Join(homePath, ".igloo")
	fmt.Println("servicePath:", servicePath)

	if info, err := os.Stat(servicePath); os.IsNotExist(err) {
		fmt.Println("starting setup process")

		os.MkdirAll(servicePath, os.ModePerm)
		db.InitializeDB(servicePath)
		initializeConfig(homePath, servicePath)

		fmt.Println("setup complete")
	} else if err != nil {
		return fmt.Errorf("conflicting service path was found. path already exist: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("conflicting service path was found. path exist but is not a directory%v", info.Name())
	}

	return nil
}
