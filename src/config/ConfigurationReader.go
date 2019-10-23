package config

// ConfigurationReader interface to read the configuration
type ConfigurationReader interface {
	Read (file string) *Config
}