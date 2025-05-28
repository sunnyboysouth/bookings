package main

import (
	"testing"

	"github.com/dsundar/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//Do nothing as test is successful
	default:
		t.Errorf("mux is not of type *chi.Mux, got %T", v)
	}
}
