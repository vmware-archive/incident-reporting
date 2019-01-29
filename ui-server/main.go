//go:generate solc --abi contracts/IncidentLog.sol  -o generate/
//go:generate abigen --abi generate/IncidentLog.abi --type IncidentLog --pkg main --out IncidentLog.go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// IncidentLogAddress is the deployed address of the incident log contract
var IncidentLogAddress = common.HexToAddress(os.Getenv("CLIENT_CONTRACT_ADDRESS"))
var keyfile = os.Getenv("CLIENT_KEYFILE")
var passphrase = os.Getenv("CLIENT_PASSPHRASE")

var user = os.Getenv("CLIENT_USER")
var password = os.Getenv("CLIENT_PASSWORD")
var url = os.Getenv("CLIENT_URL")
var ilog *IncidentLog
var transactor *bind.TransactOpts
var templateEngine *Template

func init() {

	file, err := os.Open(keyfile)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", keyfile, err)
	}

	transactor, err = bind.NewTransactor(file, passphrase)
	if err != nil {
		log.Fatalf("Failed to bind a new transctor using keyfile: %v", err)
	}

	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial(fmt.Sprintf("https://%s:%s@%s", user, password, url))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Instantiate the contract and display its name
	ilog, err = NewIncidentLog(IncidentLogAddress, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate the IncidentLog contract: %v", err)
	}

	templateEngine = &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Renderer = templateEngine

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "main", template.HTML(`
		<H2>Welcome to the great Incident Reporting tool.</H2>
		<p>This application leverages blockchain to enable co-auditing across a multi party
		system.  Any member can whistle-blow on any other without fear of losing that record.
		</p>
		`))
	})
	e.GET("/log", reportIncidentForm)
	e.POST("/log", reportIncidentHTML)
	e.GET("/log/:id", getIncidentHTML)
	e.GET("/logs", getIncidents)

	e.GET("/rest/log/:id", getIncidentJSON)
	e.POST("/rest/log", reportIncidentJSON)

	e.Logger.Fatal(e.Start(":80"))
}

func reportIncident(c echo.Context) (Incident, error) {
	// collect input as an incident
	incident, err := bindInput(c)
	if err != nil {
		return incident, err
	}

	// file the report
	_, err = ilog.ReportIncident(transactor, common.HexToAddress(incident.Reporter), incident.Message)
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
