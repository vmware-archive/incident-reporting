// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: BSD-
package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

// Incident is an entry on the IncidentLog blockchain
type Incident struct {
	Reporter  string `json:"Reporter" form:"Reporter" query:"Reporter"` // what public id reported this incident?
	Message   string `json:"Message" form:"Message" query:"Message"`    // log message for the incident
	Timestamp uint64 `json:"Timestamp" form:"Timestamp" query:"Timestamp"`
}

// Template is a pointer to html templates
type Template struct {
	templates *template.Template
}

// Render is used to render a Template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
