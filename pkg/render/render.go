package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dsundar/bookings/pkg/config"
	"github.com/dsundar/bookings/pkg/models"
)

var app *config.AppConfig

//NewTemplates sets the config for the template package

func NewTemplates(a *config.AppConfig) {
	// this function is used to create a new template cache
	// we will use this function to create a new template cache
	// and assign it to the app config
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	// this function is used to add default data to the template data
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
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
		log.Fatal("Could not get template from cache")
	}

	// this step is completely arbitary but good to track the program if there is an error
	// declare a buffer to hold the template
	buf := new(bytes.Buffer)

	// add default data to the template data
	td = AddDefaultData(td)

	// execute the template and write it to the buffer
	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal("Error executing template", err)
	}

	// Render the template by writing the buffer to the response writer
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing buffer to response writer", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// we will use this function to create a template cache from scratch
	// first declare a map to hold the template cache
	myCache := map[string]*template.Template{} //this is same as templateSet := make(map[string]*template.Template)

	// we will use the template.ParseGlob function to get all the templates in the templates directory
	pages, err := filepath.Glob("./templates/*.page.tmpl")
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
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			//pray that there is only one match
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
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
