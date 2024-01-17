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
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Post("/search-availability-json", handlers.Repo.AvailabityJSON) // making a Post request
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Post("/search-availability", handlers.Repo.PostAvailabity)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// mux.Route()

	return mux
}
