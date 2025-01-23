package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/janomonje/bed-n-breakfast/internal/config"
	driver "github.com/janomonje/bed-n-breakfast/internal/driverdb"
	"github.com/janomonje/bed-n-breakfast/internal/forms"
	"github.com/janomonje/bed-n-breakfast/internal/helpers"
	"github.com/janomonje/bed-n-breakfast/internal/models"
	"github.com/janomonje/bed-n-breakfast/internal/render"
	"github.com/janomonje/bed-n-breakfast/internal/repository"
	"github.com/janomonje/bed-n-breakfast/internal/repository/dbrepo"
)

// ******** Repository pattern *********

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPosgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (model *Repository) Home(w http.ResponseWriter, req *http.Request) {
	model.DB.AllUsers()

	render.Template(w, "home.page.tmpl", &models.TemplateData{}, req)
}

// About is the about page handler
func (model *Repository) About(w http.ResponseWriter, req *http.Request) {
	render.Template(w, "about.page.tmpl", &models.TemplateData{}, req)
}

// GeneralsQuarters is the generals-quarters page handler
func (model *Repository) GeneralsQuarters(w http.ResponseWriter, req *http.Request) {
	render.Template(w, "generals.page.tmpl", &models.TemplateData{}, req)
}

// MajorsSuite is the majors-suite page handler
func (model *Repository) MajorsSuite(w http.ResponseWriter, req *http.Request) {
	render.Template(w, "majors.page.tmpl", &models.TemplateData{}, req)
}

// Reservation is the make-reservation page handler
func (model *Repository) Reservation(w http.ResponseWriter, req *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, req)
}

// PostReservation handles the posting of a
func (model *Repository) PostReservation(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := req.Form.Get("start_date")
	ed := req.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700
	// parsing date
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(req.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: req.Form.Get("first_name"),
		LastName:  req.Form.Get("last_name"),
		Email:     req.Form.Get("email"),
		Phone:     req.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(req.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, req)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, req)
		return
	}

	newReservationID, err := model.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = model.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
	}

	model.App.Session.Put(req.Context(), "reservation", reservation)
	http.Redirect(w, req, "/reservation-details", http.StatusSeeOther)
}

// ReservationDetails renders all details of newly done reservation to the user
func (model *Repository) ReservationDetails(w http.ResponseWriter, req *http.Request) {
	reservation, ok := model.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		model.App.ErrorLog.Println("Can't get error from session")
		log.Println("cannot get item from session")
		model.App.Session.Put(req.Context(), "error", "Reservation cannot be obtained from session")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}
	model.App.Session.Remove(req.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, "reservation-details.page.tmpl", &models.TemplateData{
		Data: data,
	}, req)
}

// SearchAvailability is the search-availability page handler
func (model *Repository) SearchAvailability(w http.ResponseWriter, req *http.Request) {
	render.Template(w, "search-availability.page.tmpl", &models.TemplateData{}, req)
}

// Contact is the contact page handles the posting of a reservation form
func (model *Repository) Contact(w http.ResponseWriter, req *http.Request) {
	render.Template(w, "contact.page.tmpl", &models.TemplateData{}, req)
}

// SearchAvailability is the search-availability page handler
func (model *Repository) PostSearchAvailability(w http.ResponseWriter, req *http.Request) {
	start := req.Form.Get("start")
	end := req.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("arrival date is %s and departure date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send json response
func (model *Repository) AvailabilityJSON(w http.ResponseWriter, req *http.Request) {
	response := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	output, err := json.MarshalIndent(response, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
