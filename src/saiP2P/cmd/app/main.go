package main

import (
	"log"

	"github.com/webmakom-com/saiP2P/internal/app"
	"github.com/webmakom-com/saiP2P/tasks"
)

func main() {
	app := app.New()

	//register config with specific options
	err := app.RegisterConfig("./config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	taskManager := tasks.New(app.Cfg, app.Logger)

	app.RegisterTask(taskManager)

	app.RegisterHandlers()

	app.Run()

}
