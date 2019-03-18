//go:generate solc --abi contracts/IncidentLog.sol  -o generate/
//go:generate sh -c "abigen --abi generate/*IncidentLog.abi --type IncidentLog --pkg main --out IncidentLog.go"
// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: BSD-2
package main

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"html"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// IncidentLogAddress is the deployed address of the incident log contract
var IncidentLogAddress = common.HexToAddress(os.Getenv("CLIENT_CONTRACT_ADDRESS"))
var privateKey *ecdsa.PrivateKey

var user = html.EscapeString(os.Getenv("CLIENT_USER"))
var password = html.EscapeString(os.Getenv("CLIENT_PASSWORD"))
var url = os.Getenv("CLIENT_URL")
var session *IncidentLogSession
var templateEngine *Template
var chainID = big.NewInt(12349876)

func init() {
	var err error
	var client *ethclient.Client

	ks, account := initializeAndUnlockKeystore(password)

	// Create an IPC based RPC connection to standard url or a permissioned Concord endpoint
	if user == "" || password == "" {
		client, err = connectStandard(url)
	} else {
		client, err = connectConcord(user, password, url)
	}
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Instantiate the contract
	ilog, err := NewIncidentLog(IncidentLogAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate the IncidentLog contract: %v", err)
	}

	session = &IncidentLogSession{
		Contract: ilog,
		TransactOpts: bind.TransactOpts{
			From: account.Address,
			Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				if address != account.Address {
					return nil, errors.New("not authorized to sign this account")
				}
				signature, err := ks.SignTx(account, tx, chainID)
				if err != nil {
					return nil, err
				}
				return signature, nil
			},
		},
	}

	templateEngine = &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Renderer = templateEngine

	e.Static("/", "public/assets")

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

// func recieveEvents() {
// 	query := ethereum.FilterQuery{
// 		FromBlock: big.NewInt(0),
// 		Addresses: []common.Address{
// 			IncidentLogAddress,
// 		},
// 	}
// 	// client.Filter
// }
func reportIncident(c echo.Context) (Incident, error) {
	// collect input as an incident
	incident, err := bindInput(c)
	if err != nil {
		return incident, err
	}

	// file the report
	_, err = session.ReportIncident(common.HexToAddress(incident.Reporter), incident.Message)
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

func connectConcord(user string, password string, url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(fmt.Sprintf("https://%s:%s@%s", user, password, url))
	if err != nil {
		return nil, err
	}
	return client, nil
}
func connectStandard(url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return client, nil
}
func initializeAndUnlockKeystore(password string) (*keystore.KeyStore, accounts.Account) {
	ks := keystore.NewKeyStore("./keydir", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	err = ks.Unlock(account, password)
	if err != nil {
		log.Fatal(err)
	}
	return ks, account
}
