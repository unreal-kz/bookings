package render

import (
	"net/http"
	"testing"

	"github.com/unreal-kz/bookings/internal/models"
)

func TestAddDefualtData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefualtData(&td, r)

	if result.Flash != "123" {
		t.Error("flash with value of 123 not found")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Sesssion"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTempalteCache()
	if err != nil {
		t.Error(err)
	}

	app.TmplCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing tempalte to browser")
	}

	err = RenderTemplate(&ww, r, "none-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered tempalte that does NOT exists")
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTempalteCache()
	if err != nil {
		t.Error(err)
	}
}
