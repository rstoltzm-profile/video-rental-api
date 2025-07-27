package main

import (
	"log"

	_ "github.com/rstoltzm-profile/video-rental-api/docs/swagger"
	"github.com/rstoltzm-profile/video-rental-api/internal/app"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @Security ApiKeyAuth
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
