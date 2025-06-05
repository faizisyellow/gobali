package main

import (
	"github.com/charmbracelet/log"
	"github.com/faizisyellow/gobali/internal/db"
	"github.com/faizisyellow/gobali/internal/env"
	"github.com/faizisyellow/gobali/internal/repository"
)

func main() {

	e := &env.Env{}
	err := e.Set()
	if err != nil {
		log.Fatal(err)
	}

	conf := config{
		addr: e.GetString("ADDRESS", "http://localhost:8080"),
		env:  e.GetString("ENVIRONMENT", "Development"),
		db:   dbConfig{addr: e.GetString("DB_ADDRESS", "nil"), maxOpenConn: 30, maxIdleConn: 30, maxIdleTime: "15m"},
	}

	db, err := db.New(conf.db.addr, conf.db.maxOpenConn, conf.db.maxIdleConn, conf.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Info("database connection pool established")

	app := &application{
		configs:    conf,
		repository: repository.NewRepository(db),
	}

	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
