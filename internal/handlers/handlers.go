package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/unreal-kz/bookings/internal/config"
	"github.com/unreal-kz/bookings/internal/models"
	"github.com/unreal-kz/bookings/internal/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rp.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello"
	remoteIP := rp.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["rIP"] = remoteIP
	//send it to browser
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
func (rp *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}
func (rp *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}
func (rp *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (rp *Repository) Availabity(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}
func (rp *Repository) PostAvailabity(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Start data is %s. End data is %s. Enjoy!", start, end)))
	// render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{})
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (rp *Repository) AvailabityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      false,
		Message: "Availabel!",
	}

	out, err := json.MarshalIndent(&resp, "", "    ")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

func (rp *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}
