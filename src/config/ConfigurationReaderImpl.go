package config

import (
	"encoding/json"
	"log"
	"os"
)
// ConfigurationReaderImpl struct for configuration reader
type ConfigurationReaderImpl struct {
}

// NewReaderImpl creates a new Reader
func NewReaderImpl() *ConfigurationReaderImpl {
	return &ConfigurationReaderImpl{}
}

// Read reads the configuration from the configuration file
func (configurationReader *ConfigurationReaderImpl) Read(file string) *Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config
}