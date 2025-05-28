package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dsundar/bookings/internal/config"
	"github.com/dsundar/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	// this function is used to create a new template cache
	// we will use this function to create a new template cache
	// and assign it to the app config
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// this function is used to add default data to the template data
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// we will use this function to render the template using the html/template package to render the template
	var templateCache map[string]*template.Template
	if app.UseCache {
		//get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	//get requested template from the template cache
	t, ok := templateCache[tmpl]
	if !ok {
		//log.Fatal("Could not get template from cache")
		return errors.New("could not get template from cache")
	}

	// this step is completely arbitary but good to track the program if there is an error
	// declare a buffer to hold the template
	buf := new(bytes.Buffer)

	// add default data to the template data
	td = AddDefaultData(td, r)

	// execute the template and write it to the buffer
	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal("Error executing template", err)
	}

	// Render the template by writing the buffer to the response writer
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing buffer to response writer", err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
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
