package main

import "github.com/charmbracelet/log"

func main() {
	conf := config{
		addr: "http://localhost:8080",
		env:  "Development",
	}

	app := &application{
		configs: conf,
	}

	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
