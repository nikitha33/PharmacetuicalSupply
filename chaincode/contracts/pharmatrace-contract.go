package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type DrugContract struct {
	contractapi.Contract
}

type HistoryQueryResult struct {
	Record    *Drug  `json:"record"`
	TxId      string `json:"txId"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"isDelete"`
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

// UpdateDrugStatus allows dealers and pharmacies to update the status of a drug as it moves through the supply chain

func (d *DrugContract) UpdateDrugStatus(ctx contractapi.TransactionContextInterface, serialNumber string, newStatus string) error {
	drugExists, err := d.DrugExists(ctx, serialNumber)
	if err != nil {
		return err
	}
	if !drugExists {
		return fmt.Errorf("drug with serial number %s does not exist", serialNumber)
	}
	// Fetch the drug details
	drugJSON, err := ctx.GetStub().GetState(serialNumber)
	if err != nil {
		return fmt.Errorf("failed to read drug from world state: %v", err)
	}
	var drug Drug
	err = json.Unmarshal(drugJSON, &drug)
	if err != nil {
		return fmt.Errorf("failed to unmarshal drug: %v", err)
	}
	// Update the drug status
	drug.Status = newStatus
	// Save the updated drug record
	updatedDrugJSON, err := json.Marshal(drug)
	if err != nil {
		return fmt.Errorf("failed to marshal updated drug: %v", err)
	}
	return ctx.GetStub().PutState(serialNumber, updatedDrugJSON)
}

// QueryDrug allows querying a drug by its serial number
func (d *DrugContract) QueryDrug(ctx contractapi.TransactionContextInterface, serialNumber string) (*Drug, error) {
	drugJSON, err := ctx.GetStub().GetState(serialNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to read drug from world state: %v", err)
	}
	if drugJSON == nil {
		return nil, fmt.Errorf("drug with serial number %s does not exist", serialNumber)
	}
	var drug Drug
	err = json.Unmarshal(drugJSON, &drug)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal drug: %v", err)
	}
	return &drug, nil
}

// queryDrugsWithQueryString executes a CouchDB rich query and returns matching drugs

func (d *DrugContract) queryDrugsWithQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Drug, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var drugs []*Drug
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var drug Drug
		err = json.Unmarshal(queryResponse.Value, &drug)
		if err != nil {
			return nil, err
		}
		drugs = append(drugs, &drug)
	}
	return drugs, nil
}

// Rich query to get the history of the drug

func (c *DrugContract) GetDrugHistory(ctx contractapi.TransactionContextInterface, serialNumber string) ([]*HistoryQueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(serialNumber)
	if err != nil {
		return nil, fmt.Errorf("could not get the data. %s", err)
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("could not get the value of resultsIterator. %s", err)
		}
		var drug Drug
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &drug)
			if err != nil {
				return nil, err
			}
		} else {
			drug = Drug{
				SerialNumber: serialNumber,
			}
		}
		timestamp := response.Timestamp.AsTime()
		formattedTime := timestamp.Format(time.RFC1123)
		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &drug,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}
	return records, nil
}

// QueryDrugsByStatus queries drugs by their status (e.g., "Shipped", "Manufactured")

func (d *DrugContract) QueryDrugsByStatus(ctx contractapi.TransactionContextInterface, status string) ([]*Drug, error) {
	queryString := fmt.Sprintf(`{"selector":{"status":"%s"}}`, status)
	return d.queryDrugsWithQueryString(ctx, queryString)
}
