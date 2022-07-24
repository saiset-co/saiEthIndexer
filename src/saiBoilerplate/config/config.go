package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Configuration struct {
	HttpServer   `yaml:"http_server"`
	SocketServer `yaml:"socket_server"`
	Token        string `yaml:"token"`
	Mongo        `yaml:"mongo"`
	WebSocket    `yaml:"web_socket"`
}

type HttpServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type SocketServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
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
	Token string `yaml:"token"`
	Url   string `yaml:"url"`
}

func Load() (Configuration, error) {
	cfg := Configuration{}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	err = cleanenv.ReadConfig("config/config.yaml", &cfg)
	if err != nil {
		return Configuration{}, fmt.Errorf("config error: %w", err)
	}

	fmt.Printf("loaded configuration:%+v\n", cfg)

	return cfg, nil
}
