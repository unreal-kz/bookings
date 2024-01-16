package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache     bool
	TmplCache    map[string]*template.Template
	InfoLog      *log.Logger
	InProduction bool
	Session      *scs.SessionManager
}
