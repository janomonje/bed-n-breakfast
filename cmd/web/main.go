package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/janomonje/bed-n-breakfast/internal/config"
	"github.com/janomonje/bed-n-breakfast/internal/handlers"
	"github.com/janomonje/bed-n-breakfast/internal/models"
	"github.com/janomonje/bed-n-breakfast/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function of the application
func main() {
	// what will be put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	// Sessions
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
	}

	app.UsedCache = false
	app.TemplateCache = templateCache

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Starting application in port %s", portNumber)
	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()

}
