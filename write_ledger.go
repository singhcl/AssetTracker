package main

import (
	"fmt"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
// delete_asset() - remove a asset from state and from asset index
// 
// Shows Off DelState() - "removing"" a key/value from the ledger
//
// Inputs - Array of strings
//      0      ,         1
//     id      ,  authed_by_company
// "m999999999", "united assets"
// ============================================================================================================================
func delete_asset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("starting delete_asset")

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	id := args[0]
	// get the asset
	_, err := get_asset(stub, id)
	if err != nil{
		fmt.Println("Failed to find asset by id " + id)
		return nil, errors.New(err.Error())
	}

	// remove the asset
	err = stub.DelState(id)                                                 //remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	fmt.Println("- end delete_asset")
	return nil, nil
}

// ============================================================================================================================
// Write PharmaAsset - create a new asset, store into chaincode state
//
// Shows off building a key's JSON value manually
//
// Inputs - Array of strings
//      0      ,    1  ,  2  ,      3          ,       4
//     id      ,  color, size,     owner id    ,  authing company
// "m999999999", "blue", "35", "o9999999999999", "united assets"
// ============================================================================================================================
func write_asset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	fmt.Println("starting init_asset")

	if len(args) != 17 {
		return nil, errors.New("Incorrect number of arguments. Expecting 17")
	}

	id := args[0]
	//check if asset id already exists
	asset, err := get_asset(stub, id)
	if err == nil {
		fmt.Println("This asset already exists - " + id)
		fmt.Println(asset)
		return nil, errors.New("This asset already exists - " + id)
	}

	//build the asset json string manually
	str := `{
				"id": "` + args[0] + `", 
				"type": "` + args[1] + `", 
				"category": "` + args[2] + `",
				"assetClass": "` + args[3] + `",
				"assetTraceData": {
					"owner": "` + args[4] + `", 
					"status": "` + args[5] + `",
					"eventDateTime": "` + args[6] + `",
					"location": "` + args[7] + `",
					"geoLocation": "` + args[8] + `"
				},		
				"assetData": {
					"info": {
						"name": "` + args[9] + `", 
						"type": "` + args[10] + `",
						"pkgSize": "` + args[11] + `",
						"mfgDate": "` + args[12] + `",
						"lotNo": "` + args[13] + `",
						"expiryDate": "` + args[14] + `"
					},
					"children": [
						{
							"id": "` + args[15] + `", 
							"type": "` + args[16] + `"							
						}
					]
				}			
			}`
	err = stub.PutState(id, []byte(str))                         //store asset with id as key
	if err != nil {
		return nil, errors.New(err.Error())
	}

	fmt.Println("- end init_asset")
	return nil, nil
}