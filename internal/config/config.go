package config

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v3"
)

const (
	defaultServerPort    = 8080
	defaultTokenDuration = time.Hour
	defaultLogFile       = "./logs/app.log"
	defaultServerHost    = "0.0.0.0"
)

// Config represents an application configuration.
type Server struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type Database struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Name   string `yaml:"name"`
}
type Config struct {
	Version         string `yaml:"server_port"`
	Server          Server
	DBSource        string `yaml:"db_source"`
	Database        Database
	TokenSecretKey  string        `yaml:"token_secret_key"`
	TokenDuration   time.Duration `yaml:"token_duration"`
	LogFile         string        `yaml:"log_file"`
	FileStoragePath string        `yaml:"file_storage_path"`
	NewrelicKey     string        `yaml:"newrelic_key"`
	AssetsURL       string        `yaml:"assets_url"`
}

// Validate validates the application configuration.
func (c Config) Validate() error {
	return nil
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(file string) (*Config, error) {
	// Default config
	conf := Config{
		Server: Server{
			Port: defaultServerPort,
			Host: defaultServerHost,
		},
		TokenDuration: defaultTokenDuration,
		LogFile:       defaultLogFile,
	}

	// Load from YAML config file
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &conf); err != nil {
		return nil, err
	}

	// validation
	if err = conf.Validate(); err != nil {
		return nil, err
	}

	return &conf, err
}
