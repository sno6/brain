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

	var app *tui.App
	if os.Args[1] == "write" {
		app = tui.NewApp(b, tui.PageWrite)
	} else if os.Args[1] == "read" {
		app = tui.NewApp(b, tui.PageSearch)
	} else {
		panic("undefined")
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
