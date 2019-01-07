package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

// Incident is an entry on the IncidentLog blockchain
type Incident struct {
	Reporter  string // what public id reported this incident?
	Message   string // log message for the incident
	Timestamp uint64
}

// Template is a pointer to html templates
type Template struct {
	templates *template.Template
}

// Render is used to render a Template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
