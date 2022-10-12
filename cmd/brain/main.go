package main

import (
	"log"

	"github.com/sno6/brain"

	"github.com/sno6/brain/tui"
)

func main() {
	b, err := brain.New()
	if err != nil {
		log.Fatal(err)
	}

	app := tui.NewApp(b)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
