package main

import (
	"log"

	"github.com/Arash-Afshar/gohtmx/pkg/endpoint"
)

func main() {
	if err := endpoint.Run(); err != nil {
		log.Fatal(err)
	}
}
