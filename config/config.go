package config

import (
	"encoding/json"
	"fmt"
	"github.com/StratoAPI/Core/registry"
	"io/ioutil"
	"os"
)

var config Config

func InitializeConfig() {
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		panic(err)
	}

	if _, err := os.Stat(config.ConfigDirectory); os.IsNotExist(err) {
		panic(err)
	}

	fmt.Println("Configs initialized")
}

func InitializePluginConfigs() {
	configs := registry.GetRegistryInternal().GetConfigs()

	for name, conf := range configs {
		configFile := config.ConfigDirectory + "/" + name + ".json"
		structure := (*conf).CreateStructure()
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			bytes, err := json.MarshalIndent(structure, "", "  ")
			if err != nil {
				panic(err)
			}

			err = ioutil.WriteFile(configFile, bytes, 0644)
			if err != nil {
				panic(err)
			}
		} else {
			configFile, err := os.Open(configFile)
			if err != nil {
				panic(err)
			}

			jsonParser := json.NewDecoder(configFile)
			if err = jsonParser.Decode(&structure); err != nil {
				panic(err)
			}

			configFile.Close()
		}
		(*conf).Set(structure)
	}

	fmt.Println("Plugin Configs initialized")
}

func Get() *Config {
	return &config
}
