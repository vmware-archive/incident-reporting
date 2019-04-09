// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: BSD-2
package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

// Incident is an entry on the IncidentLog blockchain
// it should strictly follow the same struct in the IncidentLog.sol contract
type Incident struct {
	Reporter  string `json:"Reporter" form:"Reporter" query:"Reporter"`
	Message   string `json:"Message" form:"Message" query:"Message"`
	Timestamp uint64 `json:"Timestamp" form:"Timestamp" query:"Timestamp"`
	Location  string `json:"Location" form:"Location" query:"Location"`
	Resolved  bool   `json:"Resolved" form:"Resolved" query:"Resolved"`
}

// Template is a pointer to html templates
type Template struct {
	templates *template.Template
}

// Render is used to render a Template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
