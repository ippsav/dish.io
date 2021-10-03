package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dish.io/cmd/api/handlers/user"
	"github.com/dish.io/internal/database/postgres"
	"github.com/dish.io/internal/services/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	//Setting up the context
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	//Setting up a logger
	stdout := zerolog.NewConsoleWriter()
	log := zerolog.New(stdout)
	//Getting the environment variables
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Err(err).Msg("could not parse port")
	}
	dbUri := os.Getenv("POSTGRESQL_URL")
	if dbUri == "" {
		log.Error().Msg("Please set the db URI as env variable")
		os.Exit(0)
	}
	// Setting up postgres store
	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Err(err).Msg("could not open connection")
	}
	err = db.PingContext(ctx)
	if err != nil {
		log.Err(err).Msg("could not ping database")
	}
	store := &postgres.Store{DB: db}

	//Setting up the Services
	us := &user.Service{
		Store: store,
	}

	//Setting up the handlers
	uh := &handlers.UserHandler{
		Service: us,
		Log:     &log,
	}

	//Setting up the router
	router := chi.NewMux()
	// Middlewares
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	//User Routes
	router.Post("/users", uh.RegisterHandler)
	router.Get("/users", uh.LoginHandler)

	//Server mux
	log.Info().Msgf("Server is running on port %d", port)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", port), router); err != nil {
		log.Err(err).Msg("could not serve mux")
	}
}
