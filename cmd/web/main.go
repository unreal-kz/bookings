package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/unreal-kz/bookings/internal/config"
	"github.com/unreal-kz/bookings/internal/handlers"
	"github.com/unreal-kz/bookings/internal/render"
)

const portNumber = ":8080"

var (
	app     config.AppConfig
	session *scs.SessionManager
)

func main() {

	app.InProduction = false
	session = scs.New()
	session.Lifetime = 2 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // in prod should set to true
	app.Session = session

	ts, err := render.CreateTempalteCache()
	if err != nil {
		log.Fatal("Can not create Template", err)
	}
	app.TmplCache = ts
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	log.Println("Starting server on port 8080")

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	log.Fatal(srv.ListenAndServe())

}
