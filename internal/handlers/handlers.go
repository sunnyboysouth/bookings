package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	repository "github.com/dsundar/bookings/internal/Repository"
	"github.com/dsundar/bookings/internal/Repository/dbrepo"
	"github.com/dsundar/bookings/internal/config"
	"github.com/dsundar/bookings/internal/driver"
	"github.com/dsundar/bookings/internal/forms"
	"github.com/dsundar/bookings/internal/helpers"
	"github.com/dsundar/bookings/internal/models"
	"github.com/dsundar/bookings/internal/render"
	"github.com/go-chi/chi"
)

// we are goign to use repository pattern to create a handler

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo is the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresDBRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// we are building a web page that has two pages
// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// Reservations renders the makr-reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	// get the room id from the URL
	room, err := m.DB.GetRoomById(res.RoomID)
	if err != nil {
		log.Println("Error getting room by id: ", err)
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := map[string]string{
		"start_date": sd,
		"end_date":   ed,
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	//remoteIP := r.RemoteAddr
	// m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservations handles posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		log.Println("Error converting room id to int: ", err)
		helpers.ServerError(w, err)
		return
	}

	userId, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		log.Println("Error converting user id to int: ", err)
		helpers.ServerError(w, err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		UserId:    userId,
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// insert the reservation into the database
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		UserId:        userId,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// Post Availability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// parse the start and end dates into time.Time objects
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

	//w.Write([]byte(fmt.Sprintf("Starte date is %s and end date is %s", start, end)))
}

type JSONResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		OK:      false,
		Message: "Available",
	}
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// set the header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Generals renders the generals room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the Majors-suite room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// About is there about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic here

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// get the reservation data from the session
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//log.Println("cannot get reservation from session")
		m.App.ErrorLog.Println("cannot get reservation from session")
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		//m.App.ErrorLog.Println("cannot get reservation from session")
		//m.App.ErrorLog.Println("redirecting to home page")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation") // remove the reservation data from the session

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// first get the room id from the URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Println("Room id is empty =>> :", idStr)
		idStr = "1" //this is hardcoded for now - I am having trouble getting the id from the URL
	}
	roomID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting id to int: ", err, "-->", roomID)
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	log.Println("Room id is: ", roomID, " and reservation is: ", res, " redirecting to make reservation page")
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
