package main

import (
	"backend/internal/app"
	"log"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	if err := application.Run(); err != nil {
		log.Fatal("Application error:", err)
	}
}
