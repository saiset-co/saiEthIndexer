package main

import (
	"log"

	"github.com/webmakom-com/saiBoilerplate/internal/app"
	"github.com/webmakom-com/saiBoilerplate/tasks"
)

func main() {
	app := app.New()

	//register config with specific options
	err := app.RegisterConfig("./config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	t, err := tasks.NewManager(app.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer t.Logger.Sync()

	app.RegisterTask(t)

	app.RegisterHandlers()

	app.Run()

}
