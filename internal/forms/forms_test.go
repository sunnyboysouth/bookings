package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/test", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Errorf("Expected form to be valid, but it is not")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/test", nil)
	form := New(r.PostForm)

	form.Required("first_name", "last_name")

	if form.Valid() {
		t.Errorf("Expected form to be valid, but it is not")
	}

	postedData := url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Doe")
	postedData.Add("email", "e@e.org")
	postedData.Add("phone", "1234567890")

	r, _ = http.NewRequest("POST", "/test", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")

	if !form.Valid() {
		t.Errorf("Expected form to be invalid, but it is valid")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/test", nil)
	form := New(r.PostForm)

	form.Has("first_name")

	if form.Valid() {
		t.Errorf("Expected form to be valid, but it is not")
	}

	postedData := url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Doe")
	postedData.Add("email", "e@e.org")
	postedData.Add("phone", "1234567890")
	form = New(postedData)

	has := form.Has("first_name")
	if !has {
		t.Errorf("Expected form to have field 'first_name', but it does not")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/test", nil)
	form := New(r.PostForm)

	form.MinLength("first_name", 3)

	if form.Valid() {
		t.Errorf("Expected form to be valid, but it is not")
	}

	isError := form.Errors.Get("first_name")
	if isError == "" {
		t.Errorf("Expected error for 'first_name', but got none")
	}

	postedData := url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Doe")

	form = New(postedData)
	form.MinLength("first_name", 3)
	if !form.Valid() {
		t.Errorf("Expected form to be valid, but it is not")
	}

	postedData = url.Values{}
	postedData.Add("first_name", "Jo123456")
	form = New(postedData)
	form.MinLength("first_name", 1)
	if !form.Valid() {
		t.Errorf("Expected form to be valid, but it is not")
	}

	isError = form.Errors.Get("first_name")
	if isError != "" {
		t.Errorf("Expected no error for 'first_name', but got none")
	}

}
func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("email")

	if form.Valid() {
		t.Errorf("Expected form to be valid, but it is not")
	}

	postedData = url.Values{}
	postedData.Add("email", "invalid-email")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Errorf("Expected form to be valid, but it is not")
	}

	postedData = url.Values{}
	postedData.Add("email", "we@we.org")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Errorf("Expected form to be valid, but it is not")
	}
}
