// Copyright 2019 VMware, Inc.
// SPDX-License-Identifier: BSD-2
package main

import (
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// DApp is a type representing a smart contract in Ethereum or Concord
// and a contract to interact with
type DApp struct {
	User            string
	Password        string
	URL             string
	ChainType       string
	ChainID         *big.Int
	ContractHex     string
	ContractAddress common.Address
	Client          *ethclient.Client

	Keystore *keystore.KeyStore
	Account  accounts.Account
	Session  *IncidentLogSession
}

var VMwareChainID = big.NewInt(12349876)

// NewDApp creates a new dapp
func NewDApp(url, user, password string, contractAddress common.Address) *DApp {
	dapp := new(DApp)
	dapp.URL = url
	dapp.User = user
	dapp.Password = password
	dapp.ContractAddress = contractAddress

	dapp.Keystore, dapp.Account = initializeAndUnlockKeystore(dapp.Password)

	err := dapp.Connect()
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}

	return dapp
}

// Connect the DApp to an ethereum or Concord endpoint
func (dapp *DApp) Connect() error {
	var err error

	// with a username and password, we're probably a Concord endpoint
	if dapp.User != "" && dapp.Password != "" {
		dapp.Client, err = connectConcord(dapp.User, dapp.Password, dapp.URL)
		dapp.ChainType = "concord"
		dapp.ChainID = VMwareChainID

	} else {
		dapp.Client, err = connectStandard(dapp.URL)
		dapp.ChainType = "ethereum"
	}
	if err != nil {
		return err
	}
	return nil
}

// CreateSession will initiate a full session for this dapp
func (dapp *DApp) CreateSession() error {

	// Instantiate the contract
	ilog, err := NewIncidentLog(dapp.ContractAddress, dapp.Client)
	if err != nil {
		log.Printf("Failed to instantiate the IncidentLog contract: %v", err)
		return err
	}

	dapp.Session = &IncidentLogSession{
		Contract: ilog,
		TransactOpts: bind.TransactOpts{
			From: dapp.Account.Address,
			Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				if address != dapp.Account.Address {
					return nil, errors.New("not authorized to sign this account")
				}
				signature, err := dapp.Keystore.SignTx(dapp.Account, tx, dapp.ChainID)
				if err != nil {
					return nil, err
				}
				return signature, nil
			},
		},
	}
	return nil
}
