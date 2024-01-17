package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/unreal-kz/bookings/internal/config"
	"github.com/unreal-kz/bookings/internal/models"
)

var (
	functions       = template.FuncMap{}
	app             *config.AppConfig
	pathToTemplates = "./templates"
)

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefualtData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var ts map[string]*template.Template

	if app.UseCache {
		// get template cache from the app config
		ts = app.TmplCache
	} else {
		ts, _ = CreateTempalteCache()
	}

	//get requested templete
	t, ok := ts[tmpl]
	if !ok {
		// log.Fatal("Could not get template from template cache")
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer) // optional

	td = AddDefualtData(td, r)

	_ = t.Execute(buf, td) // if not use line of code, w used insted os buf

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing templates to browser", err)
		return err
	}
	return nil
}

func CreateTempalteCache() (map[string]*template.Template, error) {
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
