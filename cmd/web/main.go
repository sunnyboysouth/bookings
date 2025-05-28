package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dsundar/bookings/internal/config"
	dbDriver "github.com/dsundar/bookings/internal/driver"
	"github.com/dsundar/bookings/internal/handlers"
	"github.com/dsundar/bookings/internal/helpers"
	"github.com/dsundar/bookings/internal/models"
	"github.com/dsundar/bookings/internal/render"
)

const WebPort = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal("Cannot start application", err)
	}
	defer db.SQL.Close()

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

func run() (*dbDriver.DB, error) {
	// moved all the code to new "run()" function to emable testing of main()
	// ------Coode begins ------
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Restriction{})

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

	// connect to database
	log.Println("Connecting to database...")
	db, err := dbDriver.ConnectSQL("host=localhost port=5432 user=sundareswarandevarajan password='' dbname=bookings sslmode=disable")
	if err != nil {
		log.Fatal("Cannot connect to database so dying...", err)
		return nil, err
	}
	log.Println("Connected to database")

	// create a template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Println("Error creating template cache", err)
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false // set to true for production

	repo := handlers.NewRepo(&app, db)

	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)

	return db, nil
}
