package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/janomonje/bed-n-breakfast/internal/config"
	"github.com/janomonje/bed-n-breakfast/internal/forms"
	"github.com/janomonje/bed-n-breakfast/internal/models"
	"github.com/janomonje/bed-n-breakfast/internal/render"
)

// ******** Repository pattern *********

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, req *http.Request) {
	remoteIP := req.RemoteAddr
	m.App.Session.Put(req.Context(), "remoteIP", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{}, req)
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, req *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello World"

	remoteIP := m.App.Session.GetString(req.Context(), "remoteIP")
	stringMap["remote_ip"] = remoteIP
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	}, req)
}

// GeneralsQuarters is the generals-quarters page handler
func (m *Repository) GeneralsQuarters(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{}, req)
}

// MajorsSuite is the majors-suite page handler
func (m *Repository) MajorsSuite(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{}, req)
}

// Reservation is the make-reservation page handler
func (m *Repository) Reservation(w http.ResponseWriter, req *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, req)
}

// PostReservation handles the posting of a
func (m *Repository) PostReservation(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println("error parsing form", err)
		return
	}

	reservation := models.Reservation{
		FirstName: req.Form.Get("first_name"),
		LastName:  req.Form.Get("last_name"),
		Email:     req.Form.Get("email"),
		Phone:     req.Form.Get("phone"),
	}

	form := forms.New(req.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, req)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, req)
		return
	}

	m.App.Session.Put(req.Context(), "reservation", reservation)
	http.Redirect(w, req, "/reservation-details", http.StatusSeeOther)
}

// ReservationDetails renders all details of newly done reservation to the user
func (m *Repository) ReservationDetails(w http.ResponseWriter, req *http.Request) {
	reservation, ok := m.App.Session.Get(req.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("cannot get item from session")
		m.App.Session.Put(req.Context(), "error", "Reservation cannot be obtained from session")
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(req.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, "reservation-details.page.tmpl", &models.TemplateData{
		Data: data,
	}, req)
}

// SearchAvailability is the search-availability page handler
func (m *Repository) SearchAvailability(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{}, req)
}

// Contact is the contact page handles the posting of a reservation form
func (m *Repository) Contact(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{}, req)
}

// SearchAvailability is the search-availability page handler
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, req *http.Request) {
	start := req.Form.Get("start")
	end := req.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("arrival date is %s and departure date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send json response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, req *http.Request) {
	response := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	output, err := json.MarshalIndent(response, "", "     ")
	if err != nil {
		fmt.Println("Error marshalling data", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
