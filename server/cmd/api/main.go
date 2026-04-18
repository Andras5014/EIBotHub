package main

import (
	"log"

	"github.com/Andras5014/EIBotHub/server/internal/app"
)

func main() {
	server, err := app.New()
	if err != nil {
		log.Fatalf("start app: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
