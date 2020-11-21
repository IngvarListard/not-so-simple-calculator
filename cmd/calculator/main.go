package main

import (
	"github.com/IngvarListard/not-so-simple-calculator/pkg/server"
	"log"
)

func main() {
	srv := server.NewServer()
	err := srv.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
