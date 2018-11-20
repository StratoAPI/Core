package config

import (
	"encoding/json"
	"fmt"
	"github.com/StratoAPI/Interface/plugins"
	"io/ioutil"
	"os"
)

var config Config
var pluginConfigs = make(map[string]*plugins.Config)

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

func PushPluginConfigs(configs map[string]*plugins.Config) {
	pluginConfigs = configs
}

func InitializePluginConfigs() {
	for name, conf := range pluginConfigs {
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
