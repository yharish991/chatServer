package config

// Config struct contains the basic configuration settings
type Config struct {
	Host                 string      `json:"host"`
	Port                 string      `json:"port"`
	ConnectionType       string      `json:"connectionType"`
	LogFilePath          string      `json:"logFilePath"`
}