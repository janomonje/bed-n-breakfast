package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/janomonje/bed-n-breakfast/internal/config"
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
	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{}, req)
}

// SearchAvailability is the search-availability page handler
func (m *Repository) SearchAvailability(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{}, req)
}

// Contact is the contact page handler
func (m *Repository) Contact(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{}, req)
}

// MakeReservation is the make-reservation page handler
func (m *Repository) MakeReservation(w http.ResponseWriter, req *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{}, req)
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
