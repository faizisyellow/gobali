package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/charmbracelet/log"
	"github.com/faizisyellow/gobali/docs"
	"github.com/faizisyellow/gobali/internal/db"
	"github.com/faizisyellow/gobali/internal/env"
	"github.com/faizisyellow/gobali/internal/mailer"
	"github.com/faizisyellow/gobali/internal/repository"
)

const version = "0.1"

//	@title			Gobali Restful API
//	@version		1.0
//	@description	Restful API Documentation for Gobali app.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes	http https
//
// @BasePath	/v1
func main() {

	e := &env.Env{}
	err := e.Set()
	if err != nil {
		log.Fatal(err)
	}

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = e.GetString("ADDRESS", "localhost:8080")

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

	// metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
