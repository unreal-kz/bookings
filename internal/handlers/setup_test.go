package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/unreal-kz/bookings/internal/config"
	"github.com/unreal-kz/bookings/internal/models"
	"github.com/unreal-kz/bookings/internal/render"
)

var (
	app             config.AppConfig
	session         *scs.SessionManager
	pathToTemplates = "./../../templates"
	functions       = template.FuncMap{}
)

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 2 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // in prod should set to true

	app.Session = session

	ts, err := CreateTestTempalteCache()
	if err != nil {
		log.Fatal("Can not create Template", err)
	}
	app.TmplCache = ts

	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/search-availability", Repo.Availabity)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/make-reservation", Repo.Reservation)

	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Post("/search-availability", Repo.PostAvailabity)
	mux.Post("/search-availability-json", Repo.AvailabityJSON) // making a Post request

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTempalteCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)
	// get all of the files *page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	//range through all files
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
