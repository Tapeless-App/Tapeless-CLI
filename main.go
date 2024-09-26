package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/cmd"
	_ "tapeless.app/tapeless-cli/cmd/projects"
	_ "tapeless.app/tapeless-cli/cmd/repos"
	_ "tapeless.app/tapeless-cli/cmd/sync"
)

func main() {
	// 1. Determine the configuration path
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error finding home directory:", err)
		os.Exit(1)
	}

	configDir := filepath.Join(home, ".tapeless")
	configFile := filepath.Join(configDir, "config.json")

	// 2. Check if the config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// 3. Create the config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Println("Error creating config directory:", err)
			os.Exit(1)
		}

		// Create a default configuration
		defaultConfig := []byte("key: default_value\nanother_key: 123\n")
		if err := os.WriteFile(configFile, defaultConfig, 0644); err != nil {
			fmt.Println("Error writing default config file:", err)
			os.Exit(1)
		}

		fmt.Println("Created default configuration file at", configFile)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error since it's optional
			fmt.Println("Config file not found")
		} else {
			// Config file was found but another error was produced
		}
	}

	cmd.Execute()
}
