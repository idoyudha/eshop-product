package main

import (
	"log"

	"github.com/idoyudha/eshop-product/config"
	"github.com/idoyudha/eshop-product/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Config error: ", err)
	}

	app.Run(cfg)
}
