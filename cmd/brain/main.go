package main

import (
	"log"
	"os"

	"github.com/sno6/brain"
	"github.com/sno6/brain/tui"
)

func main() {
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	page := tui.PageIndex
	switch arg {
	case "read":
		page = tui.PageSearch
	case "write":
		page = tui.PageWrite
	}

	b, err := brain.New()
	if err != nil {
		log.Fatal(err)
	}

	app := tui.NewApp(b, page)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
