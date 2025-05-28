package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"html/template"

	"github.com/alexedwards/scs/v2"
	"github.com/dsundar/bookings/internal/config"
	"github.com/dsundar/bookings/internal/models"
	"github.com/dsundar/bookings/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = false // set to true in production
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// create a template cache
	tc, err := CreateTestTemplateCache()
	if err != nil {
		fmt.Println("Error creating template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = true // set to true for production

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(LoadCSRF)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}

// LoadCSRF adds CSRF protection to all routes
func LoadCSRF(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction, // true in production
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// we will use this function to create a template cache from scratch
	// first declare a map to hold the template cache
	myCache := map[string]*template.Template{} //this is same as templateSet := make(map[string]*template.Template)

	// we will use the template.ParseGlob function to get all the templates in the templates directory
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	// loop through the pages and parse them
	for _, page := range pages {
		// get the name of the page using filepath.Base which returns the last element of the path
		name := filepath.Base(page)
		// parse the page
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// parse the base layout template
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			//pray that there is only one match
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		// update the cache
		myCache[name] = ts
	}
	//return the cache
	return myCache, nil
}
