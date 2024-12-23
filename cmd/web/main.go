package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/janomonje/bed-n-breakfast/pkg/config"
	"github.com/janomonje/bed-n-breakfast/pkg/handlers"
	"github.com/janomonje/bed-n-breakfast/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function of the application
func main() {

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
		Handler: routes(),
	}

	err = serve.ListenAndServe()

}
