package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Configuration struct {
	Common   `yaml:"common"`
	Specific `yaml:"specific"`
}

// Common - common settings for microservice (server options, socket port and etc)
type Common struct {
	HttpServer   `yaml:"http_server"`
	SocketServer `yaml:"socket_server"`
	WebSocket    `yaml:"web_socket"`
}

// Specific - specific for current microservice settings
type Specific struct {
	Mongo `yaml:"mongo"`
	Token string `yaml:"token"`
}

type HttpServer struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

type SocketServer struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
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

type WebSocket struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	Url     string `yaml:"url"`
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
