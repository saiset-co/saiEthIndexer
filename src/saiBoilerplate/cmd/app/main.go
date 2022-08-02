package main

import (
	"log"

	"github.com/webmakom-com/saiBoilerplate/internal/app"
)

func main() {
	app := app.New()

	err := app.RegisterConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = app.RegisterStorage()
	if err != nil {
		log.Fatal(err)
	}

	app.RegisterUsecase()

	app.RegisterHandlers()

	app.Run()

}
