package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	CardStackPath  string
	HTTPServerPort int
}

var defaultConfig = Config{
	CardStackPath:  "assets/cardstack.json",
	HTTPServerPort: 8080,
}

func LoadConfig() *Config {
	cfgFile, err := os.Open("config.json")
	if err != nil {
		log.Println("Config file not found, using default config")
		return &defaultConfig
	}
	defer cfgFile.Close()

	var config Config
	err = json.NewDecoder(cfgFile).Decode(&config)
	if err != nil {
		log.Println("invalid config, using default config")
		return &defaultConfig
	}
	return &config
}
