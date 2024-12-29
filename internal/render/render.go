package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/janomonje/bed-n-breakfast/internal/config"
	"github.com/janomonje/bed-n-breakfast/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData passes data that is wanted available in every page
func AddDefaultData(templateData *models.TemplateData, req *http.Request) *models.TemplateData {
	templateData.CSRFToken = nosurf.Token(req)
	return templateData
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData, req *http.Request) {

	// gets the template cache fom the app config
	var templateCache map[string]*template.Template

	if app.UsedCache {
		// get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	template, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buffer := new(bytes.Buffer)

	templateData = AddDefaultData(templateData, req)

	err := template.Execute(buffer, templateData)
	if err != nil {
		log.Println(err)
	}

	// rendering the template
	_, err = buffer.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser. ", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, nil
	}

	// range through all the files ending with *.page.tmpl
	for _, page := range pages {
		fileName := filepath.Base(page)
		templateSet, err := template.New(fileName).ParseFiles(page)
		if err != nil {
			return myCache, nil
		}

		// looking for layouts
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, nil
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[fileName] = templateSet

	}
	return myCache, nil
}
