package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// PrivateDrug represents the structure for a private drug asset
type PrivateDrug struct {
	SerialNumber string `json:"serialNumber"` // Unique ID
	DrugName     string `json:"drugName"`     // Name of the drug
	Manufacturer string `json:"manufacturer"` // Manufacturer of the drug
	BatchNumber  string `json:"batchNumber"`  // Batch number
	ExpiryDate   string `json:"expiryDate"`   // Expiry date
	Quantity     int    `json:"quantity"`     // Quantity
	Status       string `json:"status"`       // Drug status (e.g., Manufactured, Shipped, Sold)
}

// SmartContract for managing private drugs
type SmartContract struct {
	contractapi.Contract
}

const collectionName = "PrivateDrugCollection"

// DrugExists checks if a drug exists in the private data collection
func (s *SmartContract) DrugExists(ctx contractapi.TransactionContextInterface, serialNumber string) (bool, error) {
	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, serialNumber)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve private data hash: %v", err)
	}
	return data != nil, nil
}

// CreateDrug creates a new private drug asset
func (s *SmartContract) CreateDrug(ctx contractapi.TransactionContextInterface, serialNumber string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to get MSP ID: %v", err)
	}
	// Ensure only DealerMSP can create drugs
	if clientOrgID == "DealerMSP" {
		exists, err := s.DrugExists(ctx, serialNumber)
		if err != nil {
			return "", fmt.Errorf("failed to check existence: %v", err)
		}
		if exists {
			return "", fmt.Errorf("drug with serial number %s already exists", serialNumber)
		}
		var privateDrug PrivateDrug
		// Retrieve transient data
		transientData, err := ctx.GetStub().GetTransient()
		if err != nil {
			return "", fmt.Errorf("failed to retrieve transient data: %v", err)
		}
		if len(transientData) == 0 {
			return "", fmt.Errorf("please provide private data for all required fields")
		}
		// Assign transient data fields to the private drug struct
		// if serial, exists := transientData["serialNumber"]; exists {
		// 	privateDrug.SerialNumber = string(serial)
		// } else {
		// 	return "", fmt.Errorf("serialNumber is missing in transient data")
		// }
		if drugName, exists := transientData["drugName"]; exists {
			privateDrug.DrugName = string(drugName)
		} else {
			return "", fmt.Errorf("drugName is missing in transient data")
		}
		if manufacturer, exists := transientData["manufacturer"]; exists {
			privateDrug.Manufacturer = string(manufacturer)
		} else {
			return "", fmt.Errorf("manufacturer is missing in transient data")
		}
		if batchNumber, exists := transientData["batchNumber"]; exists {
			privateDrug.BatchNumber = string(batchNumber)
		} else {
			return "", fmt.Errorf("batchNumber is missing in transient data")
		}
		if expiryDate, exists := transientData["expiryDate"]; exists {
			privateDrug.ExpiryDate = string(expiryDate)
		} else {
			return "", fmt.Errorf("expiryDate is missing in transient data")
		}
		if quantity, exists := transientData["quantity"]; exists {
			privateDrug.Quantity = int(quantity[0])
		} else {
			return "", fmt.Errorf("quantity is missing in transient data")
		}
		if status, exists := transientData["status"]; exists {
			privateDrug.Status = string(status)
		} else {
			return "", fmt.Errorf("status is missing in transient data")
		}
		// Serialize and store the private drug
		privateDrugBytes, _ := json.Marshal(privateDrug)
		err = ctx.GetStub().PutPrivateData(collectionName, serialNumber, privateDrugBytes)
		if err != nil {
			return "", fmt.Errorf("failed to store private drug: %v", err)
		}
		return fmt.Sprintf("Drug with serial number %s successfully created", serialNumber), nil
	}
	return "", fmt.Errorf("client organization %s is not authorized to create drugs", clientOrgID)
}

// ** ReadDrug retrieves a drug by its serial number from the private data collection **
func (s *SmartContract) ReadDrug(ctx contractapi.TransactionContextInterface, serialNumber string) (*PrivateDrug, error) {
	exists, err := s.DrugExists(ctx, serialNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("drug with serial number %s does not exist", serialNumber)
	}
	privateDrugBytes, err := ctx.GetStub().GetPrivateData(collectionName, serialNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private data: %v", err)
	}
	var privateDrug PrivateDrug
	err = json.Unmarshal(privateDrugBytes, &privateDrug)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private drug: %v", err)
	}
	return &privateDrug, nil
}
