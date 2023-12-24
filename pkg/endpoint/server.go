package endpoint

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	dbPkg "github.com/Arash-Afshar/gohtmx/pkg/db"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/caarlos0/env/v10"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	DB *sql.DB
}

type Config struct {
	Address      string `env:"GOHTMX_BACKEND_HOST" envDefault:"127.0.0.1:7000"`
	ReadTimeout  int    `env:"GOHTMX_BACKEND_READ_TIMEOUT" envDefault:"5"`
	WriteTimeout int    `env:"GOHTMX_BACKEND_WRITE_TIMEOUT" envDefault:"0"`
	DbURL        string `env:"GOHTMX_BACKEND_DB_URL" envDefault:"sample.sqlite"`
}

// Run runs a new HTTP endpoint with the loaded environment variables.
func Run() error {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	// Create a new chi router.
	router := chi.NewRouter()

	// Use chi middlewares.
	router.Use(middleware.Logger)

	// Handle static files (with a custom handler).
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	var err error
	db, err := dbPkg.NewDB(config.DbURL)
	if err != nil {
		return fmt.Errorf("db connection: %v", err)
	}

	h := Handler{DB: db}

	// Handle pages
	router.Get("/", h.indexViewHandler)

	// Handle apis
	router.Get("/api/sample", h.apiListSampleHandler)
	router.Post("/api/sample", h.apiNewSampleHandler)
	router.Delete("/api/sample/{name}", h.apiDeleteSampleHandler)
	server := &http.Server{
		Addr:         config.Address,
		Handler:      router, // handle all chi routes
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
	}

	return server.ListenAndServe()
}
