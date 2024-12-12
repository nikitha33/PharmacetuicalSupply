package main

import (
	"log"

	"supply/contracts"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	drugContract := new(contracts.DrugContract)

	chaincode, err := contractapi.NewChaincode(drugContract)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}
}
