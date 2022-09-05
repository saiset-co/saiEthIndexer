package config

import (
	configinternal "github.com/webmakom-com/saiP2P/internal/config-internal"
)

type Configuration struct {
	Common   configinternal.Common `yaml:"common"` // built-in framework config
	Specific `yaml:"specific"`
}

// Specific - specific for current microservice settings
type Specific struct {
	HttpCallback      `yaml:"http_callback"`
	SocketCallback    `yaml:"socket_callback"`
	WebsocketCallback `yaml:"websocket_callback"`
}

type HttpCallback struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
}

type SocketCallback struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
}

type WebsocketCallback struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
}
