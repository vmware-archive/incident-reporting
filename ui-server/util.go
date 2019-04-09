// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: BSD-2
package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo"
)

func getIndexLargestIncident() (int64, error) {
	id, err := Session.GetCount()
	if err != nil {
		return 0, fmt.Errorf("Failed to get count of incidents: %v", err)
	}
	count := id.Sub(id, big.NewInt(1))
	return count.Int64(), nil
}

func lookupLatestIncident() (Incident, error) {
	id, err := getIndexLargestIncident()
	if err != nil {
		log.Printf("Failed to get count of incidents: %v", err)
		return Incident{}, err
	}

	incident, err := lookupIncident(id)
	if err != nil {
		return Incident{}, err
	}

	return incident, nil
}

func lookupIncident(id int64) (Incident, error) {
	incident := Incident{}

	// log.Printf("trying to get incident id %d", id)
	sender, message, timestamp, location, resolved, err := Session.GetIncident(big.NewInt(id))
	if err != nil {
		log.Printf("Failed to get an incident with id %d: %v", id, err)
		return incident, fmt.Errorf("Failed to get an incident with id %d: %v", id, err)
	}

	incident.Reporter = sender.String()
	incident.Message = message
	incident.Timestamp = timestamp.Uint64()
	incident.Location = location
	incident.Resolved = resolved
	return incident, nil
}

func bindInput(c echo.Context) (Incident, error) {
	// incident message to report
	incident := new(Incident)
	err := c.Bind(incident)
	if err != nil {
		log.Printf("Failed to bind an incident: %v", err)
		return *incident, fmt.Errorf("{error:%v}", err)
	}
	return *incident, nil
}

func processOldBlockchainEvents(client *ethclient.Client) {
	blockcount, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	endblock := blockcount.Number.Uint64()

	opt := &bind.FilterOpts{
		Start: 0,
		End:   &endblock,
	}
	pastFireIncidents, err := Session.Contract.FilterFireIncident(opt)

	if err != nil {
		log.Fatalf("Failed to filter past logs: %v", err)
	}

	notDone := true
	for notDone {
		notDone = pastFireIncidents.Next()
		if notDone {
			log.Println("FireIncident: ", pastFireIncidents.Event.Message)
		}
	}
}

func watchOptsAtCurrentHead(client *ethclient.Client) *bind.WatchOpts {
	blockcount, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		log.Fatal("Failed to get latest block id: ", err)
	}
	start := blockcount.Number.Uint64()

	return &bind.WatchOpts{
		Context: context.Background(),
		Start:   &start,
	}
}
func initBlockchainEventChannels(client *ethclient.Client) {
	// Channels to watch
	FireIncidentChan = make(chan *IncidentLogFireIncident)
	GotCalledChan = make(chan *IncidentLogGotCalled)

	opts := watchOptsAtCurrentHead(client)
	var err error
	FireEventSubscription, err = Session.Contract.WatchFireIncident(opts, FireIncidentChan)
	if err != nil {
		log.Fatalf("Failed WatchFireIncident: %v", err)
	}

	CalledEventSubscription, err = Session.Contract.WatchGotCalled(opts, GotCalledChan)
	if err != nil {
		log.Fatalf("Failed GotCalledChan: %v", err)
	}
}

func connectConcord(user string, password string, url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(fmt.Sprintf("https://%s:%s@%s", user, password, url))
	client.
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
