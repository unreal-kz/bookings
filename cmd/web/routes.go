package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/unreal-kz/bookings/internal/config"
	"github.com/unreal-kz/bookings/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)

	mux.Get("/about", handlers.Repo.About)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/search-availability", handlers.Repo.Availabity)
	mux.Post("/search-availability", handlers.Repo.PostAvailabity)
	mux.Post("/search-availability-json", handlers.Repo.AvailabityJSON) // making a Post request

	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/generals-quarters", handlers.Repo.Generals)

	mux.Get("/make-reservation", handlers.Repo.Reservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
