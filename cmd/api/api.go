package main

import (
	"context"
	"errors"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/faizisyellow/gobali/internal/mailer"
	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	configs    config
	repository repository.Repository
	mailer     mailer.Client
}

type config struct {
	addr      string
	env       string
	db        dbConfig
	mail      mailConfig
	clientURL string
}

type mailConfig struct {
	sendGrid  sendgridConfig
	fromEmail string
	exp       time.Duration
}

type sendgridConfig struct {
	apiKey string
}

type dbConfig struct {
	addr        string
	maxOpenConn int
	maxIdleConn int
	maxIdleTime string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthHandler)

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.ActivateUserHandler)
			r.Post("/", app.CreateUserHandler)
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/register", app.RegisterHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.configs.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Info("signal cought, %v", s.String())

		shutdown <- srv.Shutdown(ctx)

	}()

	log.Info("server has started at", "addr", app.configs.addr, "env", app.configs.env)

	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	log.Info("server has stopped at", "addr", app.configs.addr, "env", app.configs.env)

	return nil

}
