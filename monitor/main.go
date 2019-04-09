package main

import (
	"crypto/ecdsa"
	"fmt"
	"html"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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

func main() {

	ks := keystore.NewKeyStore("./keydir", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	err = ks.Unlock(account, password)
	if err != nil {
		log.Fatal(err)
	}

	// Create an IPC based RPC connection to a remote node
	client, err := ethclient.Dial(fmt.Sprintf("https://%s:%s@%s", user, password, url))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	filterer := &bind.ContractFilterer{
        
    }

	// Instantiate the contract and display its name
	sub, err := NewIncidentLogFilterer(IncidentLogAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate the IncidentLog contract: %v", err)
	}
}
