package main

import (
	"log"

	"github.com/gobackend/tcp/api"
)

func main() {
	server := api.NewServer(":8080")
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
