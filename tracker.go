package main

import (
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type PharmaTrackerChaincode struct {
}

type PharmaAsset struct {
	assetId         		string        		`json:"assetId"`      
	assetType       		string        		`json:"assetType"`
	category   				string        		`json:"category"`
	assetClass      		string        		`json:"assetClass"`
	assetTraceData  		[]AssetTraceData 	`json:"assetTraceData"`
	assetData       		AssetData     		`json:"assetData"`
}

type AssetData struct {
	information         	AssetInfo 		 	`json:"information"`
	children     			[]AssetChildren  	`json:"children"`    
}

type AssetTraceData struct {
	owner         		string `json:"owner"`
	status   		 	string `json:"status"`
	moveDateTime      	string `json:"moveDateTime"`
	location         	string `json:"location"`
	geoLocation   		string `json:"geoLocation"`
}

type AssetInfo struct {
	assetName         	string `json:"assetName"`
	company				string `json:"company"`
	packingType   		string `json:"packingType"`
	packageSize      	string `json:"packageSize"`
	mfgDate         	string `json:"mfgDate"`
	lotNumber   		string `json:"lotNumber"`
	expiryDate      	string `json:"expiryDate"`
}

type AssetChildren struct {
	assetId         string 	`json:"assetId"`
	assetType       string 	`json:"assetType"`    
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(PharmaTrackerChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}else{
		fmt.Printf("Started Simple chaincode successfully")
	}
	
}


// ============================================================================================================================
// Init - initialize the chaincode - No initialization required
// ============================================================================================================================
func (t *PharmaTrackerChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("PharmaTrackerChaincode Is Starting Up")
	return nil, nil
}


// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *PharmaTrackerChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	if function == "write" {            //generic writes to ledger
		return write_asset(stub, args)
	} else if function == "update" {           //update an asset from state
		return update_asset(stub, args)
	} else if function == "delete" {           //deletes an asset from state
		return delete_asset(stub, args)
	}
//	} else if function == "getHistory"{        //read history of an asset (audit)
//		return getHistory(stub, args)
//	}

	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return nil, errors.New("Received unknown invoke function name - '" + function + "'")
}


// ============================================================================================================================
// Query - legacy function
// ============================================================================================================================
func (t *PharmaTrackerChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "fetch" {             //generic read ledger
		return read(stub, args)
	}
	return nil, errors.New("Unknown supported call - Query()")
}
