package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var config Config

func InitializeConfig() {
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		panic(err)
	}

	fmt.Println("Config initialized")
}

func Get() *Config {
	return &config
}
