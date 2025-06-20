package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/faizisyellow/gobali/internal/auth"
	"github.com/faizisyellow/gobali/internal/mailer"
	"github.com/faizisyellow/gobali/internal/repository"
	"github.com/faizisyellow/gobali/internal/uploader"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type application struct {
	configs        config
	repository     repository.Repository
	mailer         mailer.Client
	upload         uploader.Uploader
	authentication auth.Authenticator
}

type config struct {
	addr      string
	env       string
	db        dbConfig
	mail      mailConfig
	upload    uploadConfig
	clientURL string
	auth      authConfig
}

type tokenConfig struct {
	privateKey string
	exp        time.Duration
	iss        string
	sub        string
}

type authConfig struct {
	token tokenConfig
}

type uploadConfig struct {
	baseDir string
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
	r.Use(middleware.Timeout(60 * time.Second))

	// TODO: FIX ROUTING AUTH
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthHandler)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("http:%v/swagger/doc.json", app.configs.addr)),
		))
		r.Get("/debug/vars", expvar.Handler().ServeHTTP)

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.ActivateUserHandler)

			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.CreateUserHandler)
		})

		r.Route("/categories", func(r chi.Router) {

			r.Post("/", app.CreateCategoryHandler)
			r.Get("/", app.GetCategoriesHandler)

			r.Route("/{ID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.GetCategoryByIDHandler)
				r.Put("/", app.UpdateCategoryHandler)
				r.Delete("/", app.DeleteCategoryHandler)
			})
		})

		r.Route("/locations", func(r chi.Router) {

			r.Post("/", app.CreateLocationHandler)
			r.Get("/", app.GetLocationsHandler)

			r.Route("/{ID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.GetLocationByIdHandler)
				r.Put("/", app.UpdateLocationHandler)
				r.Delete("/", app.DeleteLocationHandler)
			})
		})

		r.Route("/types", func(r chi.Router) {

			r.Get("/", app.GetTypesHandler)
			r.Post("/", app.CreateTypeHandler)

			r.Route("/{ID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.GetTypeByIDHandler)
				r.Put("/", app.UpdateTypeHandler)
				r.Delete("/", app.DeleteTypeHandler)
			})
		})

		r.Route("/amenities", func(r chi.Router) {

			r.Get("/", app.GetAmenitiesHandler)
			r.Post("/", app.CreateAmenityHandler)

			r.Route("/{ID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.GetAmenityByIDHandler)
				r.Put("/", app.UpdateAmenityHandler)
				r.Delete("/", app.DeleteAmenityHandler)
			})
		})

		r.Route("/villas", func(r chi.Router) {

			// public
			r.Get("/", app.GetVillasHandler)
			r.Get("/{villaID}", app.VillaContentMiddleware(app.GetVillaByIdHandler))

			// authenticated
			r.Route("/", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Post("/", app.UploadImagesMiddleware(app.CreateVillaHandler, "villas"))

				r.Put("/{villaID}", app.VillaContentMiddleware(app.UploadImagesMiddleware(app.UpdateVillaHandler, "villas")))
				r.Delete("/{villaID}", app.VillaContentMiddleware(app.DeleteVillaByIdHandler))
			})
		})

		r.Route("/bookings", func(r chi.Router) {

			r.Get("/", app.GetBookingsHandler)
			r.Post("/", app.CreateBookingHandler)

			r.Route("/{bookingID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Use(app.BookingContentMiddleware)

				r.Get("/", app.GetBookingByIdHandler)
				r.Patch("/check-in", app.CheckInHandler)
				r.Patch("/check-out", app.CheckOutHandler)
				r.Delete("/", app.DeleteBookingHandler)
			})
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/register", app.RegisterHandler)
			r.Post("/login", app.LoginHandler)
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
