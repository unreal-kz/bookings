package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/unreal-kz/bookings/pkg/config"
	"github.com/unreal-kz/bookings/pkg/models"
)

var (
	functions = template.FuncMap{}
	app       *config.AppConfig
)

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefualtData(td *models.TemplateData) *models.TemplateData {

	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var ts map[string]*template.Template
	// get template cache from the app config
	if app.UseCache {
		ts = app.TmplCache
	} else {
		ts, _ = CreateTempalteCache()
	}

	//get requested templete
	t, ok := ts[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer) // optional

	td = AddDefualtData(td)

	_ = t.Execute(buf, td) // if not use line of code, w used insted os buf

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing templates to browser", err)
	}
}

func CreateTempalteCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)
	// get all of the files *page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	//range through all files
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
