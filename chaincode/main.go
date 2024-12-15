package main

import (
	"log"

	"supply/contracts"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	drugContract := new(contracts.DrugContract)
	smartContract := new(contracts.SmartContract)

	chaincode, err := contractapi.NewChaincode(drugContract, smartContract)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}
}
