package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/janomonje/bed-n-breakfast/internal/config"
	driver "github.com/janomonje/bed-n-breakfast/internal/driverdb"
	"github.com/janomonje/bed-n-breakfast/internal/handlers"
	"github.com/janomonje/bed-n-breakfast/internal/helpers"
	"github.com/janomonje/bed-n-breakfast/internal/models"
	"github.com/janomonje/bed-n-breakfast/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function of the application
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Starting application in port %s", portNumber)
	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serve.ListenAndServe()

}

func run() (*driver.DB, error) {
	// what will be put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Sessions
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bed-breakfast user=janomonje password=1234")
	if err != nil {
		log.Fatal("Cannot connect to database!")
	}
	log.Println("Connected to database")

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
		return nil, err
	}

	app.UsedCache = false
	app.TemplateCache = templateCache

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, err
}
