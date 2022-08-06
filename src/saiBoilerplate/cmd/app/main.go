package main

import (
	"context"
	"log"

	"github.com/webmakom-com/saiBoilerplate/internal/app"
	"github.com/webmakom-com/saiBoilerplate/storage"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"github.com/webmakom-com/saiBoilerplate/tasks/repo"
)

func main() {
	app := app.New()

	//register config with specific options
	err := app.RegisterConfig("../../config.config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// get storage instance (mongodb collection here)
	storage, err := storage.GetStorageInstance(context.Background(), app.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	// register storage in app
	err = app.RegisterStorage(storage)
	if err != nil {
		log.Fatal(err)
	}

	task := tasks.New(&repo.SomeRepo{
		Collection: storage.Collection,
	})

	app.RegisterTask(task)

	app.RegisterHandlers()

	app.Run()

}
