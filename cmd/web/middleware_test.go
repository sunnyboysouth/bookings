package main

import (
	"net/http"
	"testing"
)

func TestLoadCSRF(t *testing.T) {
	var myH myHandler
	h := LoadCSRF(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing as test passes
	default:
		t.Errorf("h is not of type http.Handler, got %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing as test passes
	default:
		t.Errorf("h is not of type http.Handler, got %T", v)
	}
}
