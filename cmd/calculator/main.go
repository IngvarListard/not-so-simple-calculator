package main

import (
	"github.com/IngvarListard/not-so-simple-calculator/internal/server"
	"github.com/caarlos0/env/v6"
	"log"
)

func main() {
	cfg := new(server.Config)
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err = srv.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
