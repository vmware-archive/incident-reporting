// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: BSD-2

//go:generate solc --abi contracts/IncidentLog.sol  -o generate/
//go:generate sh -c "abigen --abi generate/*IncidentLog.abi --type IncidentLog --pkg main --out IncidentLog.go"

package main

import (
	"crypto/ecdsa"
	"html"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/event"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var privateKey *ecdsa.PrivateKey

var FireOpts *bind.WatchOpts
var FireIncidentChan chan *IncidentLogFireIncident
var GotCalledChan chan *IncidentLogGotCalled
var client *ethclient.Client

// Instrumentation
var Calls int64
var FireEventSubscription event.Subscription
var CalledEventSubscription event.Subscription

var IncidentLogDapp *DApp

func init() {
	var user = html.EscapeString(os.Getenv("CLIENT_USER"))
	var password = html.EscapeString(os.Getenv("CLIENT_PASSWORD"))
	var url = os.Getenv("CLIENT_URL")
	var contractHex = os.Getenv("CLIENT_CONTRACT_ADDRESS")
	var contractAddress = common.HexToAddress(contractHex)

	if url == "" {
		log.Fatalln("Must set CLIENT_URL environment variable")
	}
	if contractHex == "" {
		log.Fatalln("Must specify contract to operate on using CLIENT_CONTRACT_ADDRESS")
	}

	IncidentLogDapp = NewDApp(url, user, password, contractAddress)
	IncidentLogDapp.CreateSession()

}

func main() {

	e := echo.New()

	// Prepare the template engine
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Debug = true
	e.Static("/", "public/assets")

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", getIndex)
	e.GET("/log", reportIncidentForm)
	e.POST("/log", reportIncidentHTML)
	e.GET("/log/:id", getIncidentHTML)
	e.GET("/logs", getIncidents)

	e.GET("/rest/log/:id", getIncidentJSON)
	e.POST("/rest/log", reportIncidentJSON)

	// go handleBlockchainEvents()

	e.Logger.Fatal(e.Start(":80"))
}

func reportIncident(c echo.Context) (Incident, error) {
	// collect input as an incident
	incident, err := bindInput(c)
	if err != nil {
		return incident, err
	}

	// file the report
	_, err = IncidentLogDapp.Session.ReportIncident(common.HexToAddress(incident.Reporter), incident.Message, incident.Location)
	if err != nil {
		log.Printf("Failed to report an incident: %v", err)
		return incident, err
	}

	// get the latest report and return it for response
	newIncident, err := lookupLatestIncident()
	if err != nil {
		return newIncident, err
	}

	return newIncident, nil
}

func getIncident(c echo.Context) (Incident, error) {
	// Incident ID from path `log/:id`
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Print(err)
		return Incident{}, err
	}
	incident, err := lookupIncident(id)
	if err != nil {
		return Incident{}, err
	}

	return incident, nil
}

func handleBlockchainEvents() {
	for {
		select {
		case err := <-FireEventSubscription.Err():
			log.Println("bailing from FireEventSubscription: ", err)
			FireEventSubscription.Unsubscribe()
			opts := watchOptsAtCurrentHead(client)
			FireEventSubscription, err = IncidentLogDapp.Session.Contract.WatchFireIncident(opts, FireIncidentChan)
			if err != nil {
				log.Fatalf("Failed WatchFireIncident: %v", err)
			}
		case err := <-CalledEventSubscription.Err():
			log.Println("bailing from CalledEventSubscription: ", err)
			CalledEventSubscription.Unsubscribe()
			opts := watchOptsAtCurrentHead(client)
			CalledEventSubscription, err = IncidentLogDapp.Session.Contract.WatchGotCalled(opts, GotCalledChan)
			if err != nil {
				log.Fatalf("Failed GotCalledChan: %v", err)
			}

		case fire := <-FireIncidentChan:
			log.Println("Got Fire Event: ", fire.Message)
		case _ = <-GotCalledChan:
			log.Println("Got Called ", Calls, " times!")
			Calls = Calls + 1
			// default:
			// 	log.Println("No events, resting.")
			// 	time.Sleep(1 * time.Second)
			// case <-time.After(time.Second * 1):
			// 	fmt.Println("timeout 1")
		}
		log.Println("wrap loop, ", len(GotCalledChan), cap(GotCalledChan))
	}
}
