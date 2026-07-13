package main

import (
	"log"

	"github.com/kirillshkro/gmart-loyalty/internal/app"
	"github.com/kirillshkro/gmart-loyalty/internal/config"
)

func main() {
	cfg := config.GetAppConfig()
	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
