package main

import (
	"log"
	"os"

	"github.com/sno6/brain"
	"github.com/sno6/brain/tui"
)

func main() {
	b, err := brain.New()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) == 1 {
		if err := b.Write(); err != nil {
			log.Fatal(err)
		}

		return
	}

	if os.Args[1] == "read" {
		app := tui.NewApp(b)
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	}
}
