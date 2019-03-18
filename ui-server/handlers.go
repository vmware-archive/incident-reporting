// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: BSD-2
package main

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo"
)

// e.GET("/log", reportIncidentForm)
func reportIncidentForm(c echo.Context) error {
	return c.Render(http.StatusOK, "main", template.HTML(`
		<h2>Report an Incident</h2>
		<form method="POST">

  		<p>What happened:<br>
		<input type="text" name="Message"></p>

  		<p>Your account ID:<br>
		<input type="text" name="Reporter"></p>

		<button class="btn btn-primary">Report</button>
		</form>`))
}

func reportIncidentHTML(c echo.Context) error {
	incident, err := reportIncident(c)
	if err != nil {
		return c.Render(http.StatusBadRequest, "main", template.HTML("<p>Error "+err.Error()+"</p>"))
	}

	return c.Render(http.StatusOK, "incident", incident)
}

func reportIncidentJSON(c echo.Context) error {
	incident, err := reportIncident(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, incident)
}

// e.GET("/log/:id", getIncident)
func getIncidentHTML(c echo.Context) error {
	incident, err := getIncident(c)
	if err != nil {
		return c.Render(http.StatusGone, "main", template.HTML("<p>"+err.Error()+"</p>"))
	}

	return c.Render(http.StatusOK, "incident", incident)
}

// e.GET("/rest/log/:id", getIncidentJSON)
func getIncidentJSON(c echo.Context) error {
	incident, err := getIncident(c)
	if err != nil {
		return c.JSON(http.StatusGone, err)
	}

	return c.JSON(http.StatusOK, incident)
}

// e.GET("/log/:id", getIncidents)
func getIncidents(c echo.Context) error {
	var incidents []Incident
	var index int64

	count, err := getIndexLargestIncident()
	if err != nil {
		return c.Render(http.StatusGone, "main", template.HTML("<p>"+err.Error()+"</p>"))
	}

	for ; index <= count; index++ {
		i, err := lookupIncident(index)
		if err != nil {
			return c.Render(http.StatusGone, "main", template.HTML("<p>"+err.Error()+"</p>"))
		}
		incidents = append(incidents, i)
	}

	return c.Render(http.StatusOK, "incidents", incidents)
}
