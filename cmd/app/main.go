package main

import (
	"log"
	app "docs_storage/internal/app"
)

func main() {
	config, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application := app.NewApp(config)
	if err := application.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}