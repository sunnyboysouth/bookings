package models

import "github.com/dsundar/bookings/internal/forms"

// TemplateData holds the data sent from handlers to templates
type TemplateData struct {
	// add any data you want to pass to the template
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float64
	Data      map[string]interface{}
	CSRFToken string      // add Cross-Site Request Forgery token that is used to protect against CSRF attacks
	Flash     string      // add flash message that is used to display messages to the user
	Warning   string      // add warning message that is used to display warning messages to the user
	Error     string      // add error message that is used to display error messages to the user
	Form      *forms.Form // add form struct that is used to handle form data
}
