package config

import (
	"fmt"

	"github.com/tkanos/gonfig"
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

func Load() Configuration {
	var config Configuration
	err := gonfig.GetConf("config.json", &config)

	if err != nil {
		fmt.Println("Configuration problem:", err)
		panic(err)
	}

	return config
}
