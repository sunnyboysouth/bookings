package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/dsundar/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTest = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"genq", "/generals-quarters", "GET", http.StatusOK},
	{"majs", "/majors-suite", "GET", http.StatusOK},
	{"search", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},

	// {"makeres", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// {"ressummary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	// {"post-search-avail", "/search-availability", "POST", []postData{
	// 	{"start", "2023-10-01"},
	// 	{"end", "2023-10-02"},
	// }, http.StatusOK},
	// {"post-search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{"start", "2023-10-01"},
	// 	{"end", "2023-10-02"},
	// }, http.StatusOK},
	// {"post-make-res", "/make-reservation", "POST", []postData{
	// 	{"first_name", "John"},
	// 	{"last_name", "Doe"},
	// 	{"email", "me@email.org"},
	// 	{"phone", "1234567890"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes) // ts standas for test server
	defer ts.Close()

	for _, e := range theTest {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Logf("Error making GET request: %v", err)
			t.Fatal()
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("For test %s, Expected status code %d, got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Generals Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()                  // create a new response recorder, it fakes the place of http.ResponseWriter
	session.Put(ctx, "reservation", reservation)  // put the reservation in thcoveragee session
	handler := http.HandlerFunc(Repo.Reservation) // create a new handler

	handler.ServeHTTP(rr, req) // serve the request
	if rr.Code != http.StatusOK {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusOK, rr.Code)
	}

	// test case where reservation is not in session i.e. reset everything
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)  // put back the context in the request
	rr = httptest.NewRecorder() // create a new response recorder

	handler.ServeHTTP(rr, req) // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusSeeOther, rr.Code)
	}

	//test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)  // put back the context in the request
	rr = httptest.NewRecorder() // create a new response recorder
	reservation.RoomID = 100

	session.Put(ctx, "reservation", reservation) // put the reservation in the session
	handler.ServeHTTP(rr, req)                   // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusSeeOther, rr.Code)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	// reqBody := "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Doe")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=jdoe@doe.org")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "user_id=1")

	postData := url.Values{}
	postData.Add("start_date", "2050-01-01")
	postData.Add("end_date", "2050-01-02")
	postData.Add("room_id", "1")
	postData.Add("first_name", "John")
	postData.Add("last_name", "Doe")
	postData.Add("email", "jdoe@doe.org")
	postData.Add("phone", "1234567890")
	postData.Add("user_id", "1")

	//test with valid data ***********************************************************
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()                      // create a new response recorder
	handler := http.HandlerFunc(Repo.PostReservation) // create a new handler
	handler.ServeHTTP(rr, req)                        // serve the request
	if rr.Code != http.StatusSeeOther {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusSeeOther, rr.Code)
	}

	//test with no data ***********************************************************
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                      // create a new response recorder
	handler = http.HandlerFunc(Repo.PostReservation) // create a new handler
	handler.ServeHTTP(rr, req)                       // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}

	//test with invalid start data ***********************************************************
	postData.Set("start_date", "invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                      // create a new response recorder
	handler = http.HandlerFunc(Repo.PostReservation) // create a new handler
	handler.ServeHTTP(rr, req)                       // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For invalid start date %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}

	//test with invalid end data ***********************************************************
	postData.Set("start_date", "2050-01-01")
	postData.Set("end_date", "invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                      // create a new response recorder
	handler = http.HandlerFunc(Repo.PostReservation) // create a new handler
	handler.ServeHTTP(rr, req)                       // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For invalid end date %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}

	//test with invalid room id data ***********************************************************
	postData.Set("end_date", "2050-01-02")
	postData.Set("room_id", "invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                      // create a new response recorder
	handler = http.HandlerFunc(Repo.PostReservation) // create a new handler
	handler.ServeHTTP(rr, req)                       // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For invalid Room ID %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}

	//test with invalid first name data ***********************************************************
	postData.Set("room_id", "1")
	postData.Set("first_name", "j")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                      // create a new response recorder
	handler = http.HandlerFunc(Repo.PostReservation) // create a new handler
	handler.ServeHTTP(rr, req)                       // serve the request
	if rr.Code != http.StatusSeeOther {
		t.Errorf("For invalid First name %s, expected status code %d, got %d", req.URL.Path, http.StatusSeeOther, rr.Code)
	}

	//test for failure to insert reservation into database
	postData.Set("room_id", "2")
	postData.Set("first_name", "john")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                      // create a new response recorder
	handler = http.HandlerFunc(Repo.PostReservation) // create a new handler
	handler.ServeHTTP(rr, req)                       // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("this has terribly failed to insert to DB %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}

	//test for failure to insert reservation into database
	postData.Set("room_id", "1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                      // create a new response recorder
	handler = http.HandlerFunc(Repo.PostReservation) // create a new handler
	handler.ServeHTTP(rr, req)                       // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("failed to insert to DB rID 1000 %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println("Error loading session: ", err)
		return nil
	}
	return ctx
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	// test with valid data
	postData := url.Values{}
	postData.Add("start_date", "2050-01-01")
	postData.Add("end_date", "2050-01-02")
	postData.Add("room_id", "1")

	// create a new request with the body created above
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()                       // create a new response recorder
	handler := http.HandlerFunc(Repo.AvailabilityJSON) // create a new handler
	handler.ServeHTTP(rr, req)                         // serve the request

	var j JSONResponse
	err := json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	// Test with no data
	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                       // create a new response recorder
	handler = http.HandlerFunc(Repo.AvailabilityJSON) // create a new handler
	handler.ServeHTTP(rr, req)                        // serve the request

	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	// Test with invalid start date
	postData.Set("start_date", "invalid")

	// create a new request with the body created above
	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                       // create a new response recorder
	handler = http.HandlerFunc(Repo.AvailabilityJSON) // create a new handler
	handler.ServeHTTP(rr, req)                        // serve the request

	err = json.Unmarshal(rr.Body.Bytes(), &j)
	if err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	// test with valid data
	postData := url.Values{}
	postData.Add("start_date", "2050-01-01")
	postData.Add("end_date", "2050-01-02")

	// create a new request with the body created above
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()                       // create a new response recorder
	handler := http.HandlerFunc(Repo.PostAvailability) // create a new handler
	handler.ServeHTTP(rr, req)                         // serve the request

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}

	// Test with no data
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                       // create a new response recorder
	handler = http.HandlerFunc(Repo.PostAvailability) // create a new handler
	handler.ServeHTTP(rr, req)                        // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}

	// Test with invalid start date
	postData.Set("start_date", "invalid")

	// create a new request with the body created above
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                       // create a new response recorder
	handler = http.HandlerFunc(Repo.PostAvailability) // create a new handler
	handler.ServeHTTP(rr, req)                        // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}

	// Test with invalid end date
	postData.Set("start_date", "2050-01-01")
	postData.Set("end_date", "invalid")

	// create a new request with the body created above
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                       // create a new response recorder
	handler = http.HandlerFunc(Repo.PostAvailability) // create a new handler
	handler.ServeHTTP(rr, req)                        // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}

	// Test with invalid start date
	postData.Set("start_date", "2020-01-01")
	postData.Set("end_date", "2020-01-02")

	// create a new request with the body created above
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx) // put back the context in the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()                       // create a new response recorder
	handler = http.HandlerFunc(Repo.PostAvailability) // create a new handler
	handler.ServeHTTP(rr, req)                        // serve the request
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("For %s, expected status code %d, got %d", req.URL.Path, http.StatusTemporaryRedirect, rr.Code)
	}
}
