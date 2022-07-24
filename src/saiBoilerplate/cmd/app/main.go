package main

import (
	"log"

	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/internal/app"
)

func main() {
	// todo: cli app?
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// srv := server.NewServer(cfg, false)

	// if cfg.SocketServer.Host != "" {
	// 	go srv.SocketStart()
	// }

	app.Run(&cfg)

	//srv.Start()
}
