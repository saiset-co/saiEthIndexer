package main

import (
	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/server"
)

func main() {
	// todo: cli app?
	cfg := config.Load()
	srv := server.NewServer(cfg, false)

	if cfg.SocketServer.Host != "" {
		go srv.SocketStart()
	}

	srv.Start()
}
