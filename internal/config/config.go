package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	ConfigFilePath string
)

func Initialize() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting the home directory of the user")
		return err
	}
	configFileFolder := filepath.Join(homedir, ".fabriquik")
	err = os.MkdirAll(configFileFolder, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return err
	}
	ConfigFilePath = filepath.Join(homedir, ".fabriquik", "config.json")
	return nil
}
