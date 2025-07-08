package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Returns the config from config file stored on the file system.
//
// Contains the database URL and the current users name.
func Read() (Config, error) {
	config := Config{}

	// Get the config file path.
	// Read the JSON and decode it to Config struct.

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("unable to get config file path: %w", err)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, fmt.Errorf("unable to read config file: %w", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("unable to parse the json: %w", err)
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get home directory: %w", err)
	}
	return fmt.Sprintf("%s/%s", homeDir, configFileName), nil
}

// Updates the current users name to the config file stored on the file system.
func (cfg Config) SetUser(name string) error {

	cfg.CurrentUserName = name

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	configFileName, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(configFileName, data, 0666)
	if err != nil {
		return err
	}
	return nil
}
