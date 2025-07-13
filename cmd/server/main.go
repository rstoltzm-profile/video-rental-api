package main

import (
	"log"

	"github.com/rstoltzm-profile/video-rental-api/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
