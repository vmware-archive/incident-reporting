package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/labstack/echo"
)

func getIndexLargestIncident() (int64, error) {
	id, err := session.GetCount()
	if err != nil {
		return 0, fmt.Errorf("Failed to get count of incidents: %v", err)
	}
	count := id.Sub(id, big.NewInt(1))
	return count.Int64(), nil
}

func lookupLatestIncident() (Incident, error) {
	id, err := getIndexLargestIncident()
	if err != nil {
		log.Fatalf("Failed to get count of incidents: %v", err)
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

	log.Printf("trying to get incident id %d", id)
	sender, message, timestamp, err := session.GetIncident(big.NewInt(id))
	if err != nil {
		log.Printf("Failed to get an incident with id %d: %v", id, err)
		return incident, fmt.Errorf("Failed to get an incident with id %d: %v", id, err)
	}

	incident.Reporter = sender.String()
	incident.Message = message
	incident.Timestamp = timestamp.Uint64()
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
