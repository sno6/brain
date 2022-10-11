package main

import (
	"log"

	"github.com/sno6/brain"
)

func main() {
	b, err := brain.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := b.Write(); err != nil {
		log.Fatal(err)
	}
}
