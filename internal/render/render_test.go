package render

import (
	"net/http"
	"testing"

	"github.com/dsundar/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	session.Put(r.Context(), "flash", "test-flash")
	result := AddDefaultData(&td, r)

	if result.Flash != "test-flash" {
		t.Errorf("test failed, expected flash value to be empty")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Errorf("Error creating template cache: %v", err)
	}
	app.TemplateCache = tc
	app.UseCache = true

	r, err := getSession()
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	var ww myWriter
	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Errorf("Error rendering template: %v", err)
	}

	err = Template(&ww, r, "non-existing.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Errorf("rendered templatye that does not exist %v", err)
	}
}

func TestNewTemplates(t *testing.T) {
	// this function is used to create a new template cache
	// we will use this function to create a new template cache
	// and assign it to the app config
	NewRenderer(app)
	if app.TemplateCache == nil {
		t.Errorf("app.TemplateCache is nil")
	}
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Errorf("Error creating template cache: %v", err)
	}
}
