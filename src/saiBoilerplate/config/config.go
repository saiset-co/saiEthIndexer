package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	configinternal "github.com/webmakom-com/saiBoilerplate/internal/config-internal"
)

type Configuration struct {
	Common   configinternal.Common `yaml:"common"` // built-in framework config
	Specific `yaml:"specific"`
}

// Specific - specific for current microservice settings
type Specific struct {
	Mongo `yaml:"mongo"`
	Token string `yaml:"token"`
}

type Mongo struct {
	Atlas      bool   `yaml:"atlas"`
	User       string `yaml:"user"`
	Pass       string `yaml:"pass"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

func Load() (Configuration, error) {
	cfg := Configuration{}

	err := cleanenv.ReadConfig("config/config.yaml", &cfg)
	if err != nil {
		return Configuration{}, fmt.Errorf("config error: %w", err)
	}

	// fmt.Printf("loaded configuration:%+v\n", cfg)

	return cfg, nil
}
