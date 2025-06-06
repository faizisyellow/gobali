package main

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/faizisyellow/gobali/internal/db"
	"github.com/faizisyellow/gobali/internal/env"
	"github.com/faizisyellow/gobali/internal/mailer"
	"github.com/faizisyellow/gobali/internal/repository"
)

func main() {

	e := &env.Env{}
	err := e.Set()
	if err != nil {
		log.Fatal(err)
	}

	mailConf := mailConfig{
		sendGrid:  sendgridConfig{apiKey: e.GetString("API_URL_SENDGRID", "")},
		fromEmail: e.GetString("SENDER_EMAIL", ""),
		exp:       time.Hour * 24 * 3,
	}

	conf := config{
		addr:      e.GetString("ADDRESS", "localhost:8080"),
		env:       e.GetString("ENVIRONMENT", "Development"),
		db:        dbConfig{addr: e.GetString("DB_ADDRESS", "nil"), maxOpenConn: 30, maxIdleConn: 30, maxIdleTime: "15m"},
		mail:      mailConf,
		clientURL: e.GetString("CLIENT_URL", "localhost:5173"),
	}

	db, err := db.New(conf.db.addr, conf.db.maxOpenConn, conf.db.maxIdleConn, conf.db.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Info("database connection pool established")

	sendGridMail := mailer.NewSendGrid(conf.mail.sendGrid.apiKey, conf.mail.fromEmail)

	app := &application{
		configs:    conf,
		repository: repository.NewRepository(db),
		mailer:     sendGridMail,
	}

	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
