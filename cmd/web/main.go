package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dsundar/bookings/pkg/config"
	"github.com/dsundar/bookings/pkg/handlers"
	"github.com/dsundar/bookings/pkg/render"
)

const WebPort = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = false // set to true in production
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// create a template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Println("Error creating template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false // set to true for production

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	//	http.HandleFunc("/", handlers.Repo.Home)
	//	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Starting web server on port", WebPort)

	srv := &http.Server{
		Addr:    WebPort,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
