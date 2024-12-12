package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type DrugContract struct {
	contractapi.Contract
}

type Drug struct {
	SerialNumber string `json:"serialNumber"`
	DrugName     string `json:"drugName"`
	Manufacturer string `json:"manufacturer"`
	Status       string `json:"status"` // e.g., "Manufactured", "Shipped", "Received", "Sold"
}

// CreateDrug allows the manufacturer to create a new drug with a serial number

func (d *DrugContract) CreateDrug(ctx contractapi.TransactionContextInterface, serialNumber string, drugName string, manufacturer string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}
	if clientOrgID == "ManufacturerMSP" {

		drugExists, err := d.DrugExists(ctx, serialNumber)
		if err != nil {
			return "", err
		}
		if drugExists {
			return "", fmt.Errorf("drug with serial number %s already exists", serialNumber)
		}
		drug := Drug{
			SerialNumber: serialNumber,
			DrugName:     drugName,
			Manufacturer: manufacturer,
			Status:       "Manufactured",
		}

		drugJSON, err := json.Marshal(drug)
		if err != nil {
			return "", err
		}

		err = ctx.GetStub().PutState(serialNumber, drugJSON)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("successfully added drug %v", serialNumber), nil

	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

// DrugExists checks if a drug exists
func (d *DrugContract) DrugExists(ctx contractapi.TransactionContextInterface, serialNumber string) (bool, error) {
	drugJSON, err := ctx.GetStub().GetState(serialNumber)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return drugJSON != nil, nil
}

// // UpdateDrugStatus allows dealers and pharmacies to update the status of a drug as it moves through the supply chain

// func (d *DrugContract) UpdateDrugStatus(ctx contractapi.TransactionContextInterface, serialNumber string, newStatus string) error {
//    drugExists, err := s.DrugExists(ctx, serialNumber)
//    if err != nil {
//        return err
//    }
//    if !drugExists {
//        return fmt.Errorf("drug with serial number %s does not exist", serialNumber)
//    }
//    // Fetch the drug details
//    drugJSON, err := ctx.GetStub().GetState(serialNumber)
//    if err != nil {
//        return fmt.Errorf("failed to read drug from world state: %v", err)
//    }
//    var drug Drug
//    err = json.Unmarshal(drugJSON, &drug)
//    if err != nil {
//        return fmt.Errorf("failed to unmarshal drug: %v", err)
//    }
//    // Update the drug status
//    drug.Status = newStatus
//    // Save the updated drug record
//    updatedDrugJSON, err := json.Marshal(drug)
//    if err != nil {
//        return fmt.Errorf("failed to marshal updated drug: %v", err)
//    }
//    return ctx.GetStub().PutState(serialNumber, updatedDrugJSON)
// }

// // QueryDrug allows querying a drug by its serial number
// func (d *DrugContract) QueryDrug(ctx contractapi.TransactionContextInterface, serialNumber string) (*Drug, error) {
// 	drugJSON, err := ctx.GetStub().GetState(serialNumber)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read drug from world state: %v", err)
// 	}
// 	if drugJSON == nil {
// 		return nil, fmt.Errorf("drug with serial number %s does not exist", serialNumber)
// 	}
// 	var drug Drug
// 	err = json.Unmarshal(drugJSON, &drug)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal drug: %v", err)
// 	}
// 	return &drug, nil
//  }
