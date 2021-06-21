package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

const (
	defaultServerPort         = 8080
	defaultJWTExpirationHours = 72
	defaultLogFile            = "./logs/app.log"
)

// Config represents an application configuration.
type Server struct {
	Port int `yaml:"port"`
}

type Database struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Name   string `yaml:"name"`
}
type Config struct {
	Version       string `yaml:"version"`
	Server        Server
	DBSource      string `yaml:"db_source"`
	Database      Database
	JWTSigningKey string `yaml:"jwt_signing_key"`
	JWTExpiration int    `yaml:"jwt_expiration"`
	LogFile       string `yaml:"log_file"`
}

// Validate validates the application configuration.
func (c Config) Validate() error {
	return nil
}

// Load returns an application configuration which is populated from the given configuration file and environment variables.
func Load(file string) (*Config, error) {
	// default config
	conf := Config{
		Server: Server{
			Port: defaultServerPort,
		},
		JWTExpiration: defaultJWTExpirationHours,
		LogFile:       defaultLogFile,
	}

	// load from YAML config file
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
