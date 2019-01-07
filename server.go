//go:generate abigen --sol contracts/IncidentLog.sol --pkg main --out IncidentLog.go
package main

import (
	"fmt"
	"html/template"
	"log"
	"math/big"
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
	e.GET("/log", reportIncidentSetup)
	e.POST("/log", reportIncident)
	e.GET("/log/:id", getIncident)
	e.GET("/logs", getIncidents)

	e.Logger.Fatal(e.Start(":1234"))
}

// e.GET("/log", reportIncidentSetup)
func reportIncidentSetup(c echo.Context) error {
	return c.Render(http.StatusOK, "main", template.HTML(`
		<h2>Report an Incident</h2>
	    <form method="POST">
  		What happened:<br>
		<input type="text" name="message">
		<button class="btn btn-primary">Report</button>
		</form>`))
}

// e.POST("/log", reportIncident)
func reportIncident(c echo.Context) error {
	// incident message to report
	message := c.FormValue("message")

	_, err := ilog.ReportIncident(transactor, transactor.From, message)
	if err != nil {
		log.Fatalf("Failed to report an incident: %v", err)
	}

	id, err := getIndexLargestIncident()
	if err != nil {
		log.Fatalf("Failed to get count of incidents: %v", err)
	}
	return c.Redirect(http.StatusMovedPermanently, "/log/"+strconv.FormatInt(id, 10))

	// return c.Render(http.StatusOK, "main", template.HTML(`Go to <a href="/log/`+strconv.FormatInt(id, 10)+`">the log entry</a>`))
}

func getIndexLargestIncident() (int64, error) {
	id, err := ilog.GetCount(&bind.CallOpts{})
	if err != nil {
		return 0, fmt.Errorf("Failed to get count of incidents: %v", err)
	}
	count := id.Sub(id, big.NewInt(1))
	return count.Int64(), nil
}

// e.GET("/log/:id", getIncident)
func getIncident(c echo.Context) error {
	// Incident ID from path `log/:id`
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert Atoi: %v", err)
	}

	sender, message, timestamp, err := ilog.GetIncident(&bind.CallOpts{}, big.NewInt(id))
	incident := Incident{
		Reporter:  sender.String(),
		Message:   message,
		Timestamp: timestamp.Uint64(),
	}
	if err != nil {
		log.Printf("Failed to get an incident with id %d: %v", id, err)
		return c.Render(http.StatusGone, "main", template.HTML("<p>There does not exist an incident with id "+c.Param("id")+"</p>"))
	}

	return c.Render(http.StatusOK, "incident", incident)
}

// e.GET("/log/:id", getIncidents)
func getIncidents(c echo.Context) error {
	var incidents []Incident
	var index int64

	count, err := getIndexLargestIncident()
	if err != nil {
		log.Fatalf("Failed to get count of incidents: %v", err)
	}

	log.Printf("got %v incidents", count)
	for ; index < count; index++ {
		sender, message, timestamp, err := ilog.GetIncident(&bind.CallOpts{}, big.NewInt(index))
		i := Incident{
			Reporter:  sender.String(),
			Message:   message,
			Timestamp: timestamp.Uint64(),
		}
		incidents = append(incidents, i)
		if err != nil {
			log.Printf("Failed to get an incident with id %d: %v", index, err)
			return c.Render(http.StatusGone, "main", template.HTML("<p>There does not exist an incident with id "+c.Param("id")+"</p>"))
		}

	}
	log.Printf("%v", incidents)
	return c.Render(http.StatusOK, "incidents", incidents)
}
