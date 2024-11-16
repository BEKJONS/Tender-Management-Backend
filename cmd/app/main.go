package main

import (
	"tender_management/config"
	"tender_management/internal/app"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
