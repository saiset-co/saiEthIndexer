package main

import (
	"context"
	"fmt"
	"log"

	"github.com/webmakom-com/saiBoilerplate/internal/app"
	"github.com/webmakom-com/saiBoilerplate/storage"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"github.com/webmakom-com/saiBoilerplate/tasks/repo"
)

func main() {
	app := app.New()

	//register config with specific options
	err := app.RegisterConfig("./config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// get storage instance (mongodb collection here)
	storage, client, err := storage.GetStorageInstance(context.Background(), app.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			fmt.Printf("error when disconnect to mongo instance : %s", err.Error())
		}
	}()

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
