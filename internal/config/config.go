package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/dsundar/bookings/internal/models"
)

// AppConfig holds the application configuration
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
