package main

import (
	"context"
	"log"
	"mcpbinance/internal/application"
	"mcpbinance/internal/config"
)

func main() {
	ctx := context.TODO()
	cfg, err := config.MustConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := application.NewApplication(cfg)
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
