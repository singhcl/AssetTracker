package main

import (
	"fmt"
	"errors"
	"encoding/json"
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
	fmt.Println("starting write asset")

	if len(args) < 18 {
		return nil, errors.New("Incorrect number of arguments. Expecting at least 18")
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
	str := `{"assetId": "` + args[0] + `",
		"assetType": "` + args[1] + `",
		"category": "` + args[2] + `",
		"assetClass": "` + args[3] + `",
		"assetTraceData":[ {
			"owner": "` + args[4] + `", 
			"status": "` + args[5] + `",
			"assetTraceData": "` + args[6] + `",
			"location": "` + args[7] + `",
			"geoLocation": "` + args[8] + `"
		}],		
		"assetData": {
			"information": {
				"assetName": "` + args[9] + `",
				"company": "` + args[10] + `",
				"packingType": "` + args[11] + `",
				"packageSize": "` + args[12] + `",
				"mfgDate": "` + args[13] + `",
				"lotNumber": "` + args[14] + `",
				"expiryDate": "` + args[15] + `"
			},
			"children": [
				{
					"assetId": "` + args[16] + `", 
					"assetType": "` + args[17] + `"							
				}
			]
		}			
	}`
	
	fmt.Println("Input PharmaAsset Object - " + str)
	var pharmaAsset PharmaAsset
	err = json.Unmarshal([]byte(str), &pharmaAsset)
	if err != nil {
		fmt.Println("Error while unmarshalling "+err.Error())
		return nil, errors.New(err.Error())
	}
	fmt.Println("PharmaAsset Object after marshalling - ")
	if len(args) > 18 {
		for i := 18; i < len(args); i=i+2 {
			var child AssetChildren
			child.assetId=args[i]
			child.assetType=args[i+1]
			fmt.Println("New child to be appended - ")
			pharmaAsset.assetData.children = append(pharmaAsset.assetData.children, child)
		}
		fmt.Println("PharmaAsset Object after appending children - ")
	}
	
	inputByteStr, err := json.Marshal(pharmaAsset)
	if err != nil {
		fmt.Println("Error while marshalling "+err.Error())
		return nil, errors.New(err.Error())
	}
	
	err = stub.PutState(id, inputByteStr)                         //store asset with id as key
	if err != nil {
		return nil, errors.New(err.Error())
	}

	fmt.Println("- end write asset")
	return nil, nil
}

func update_asset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("starting update asset")

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}

	id := args[0]
	// get the asset
	var asset PharmaAsset
	assetAsBytes, err := stub.GetState(id)                  //getState retreives a key/value from the ledger
	if err != nil {                                          //this seems to always succeed, even if key didn't exist
		return nil, errors.New("Failed to find asset - " + id)
	}
	json.Unmarshal(assetAsBytes, &asset)                   //un stringify it aka JSON.parse()

	if asset.assetId != id {                                     //test if marble is actually here or just nil
		return nil, errors.New("Asset does not exist - " + id)
	}
	
	var traceData AssetTraceData
	traceData.owner=args[1]
	traceData.status=args[2]
	traceData.moveDateTime=args[3]
	traceData.location=args[4]
	traceData.geoLocation=args[5]
	fmt.Println("New trace data to be appended - ")
	asset.assetTraceData = append(asset.assetTraceData, traceData)
	fmt.Println("PharmaAsset object post update of trace data - ")

	inputByteStr, err := json.Marshal(asset)
	if err != nil {
		fmt.Println("Error while marshalling in update "+err.Error())
		return nil, errors.New(err.Error())
	}
	
	err = stub.PutState(id, inputByteStr)                         //store asset with id as key
	if err != nil {
		return nil, errors.New(err.Error())
	}

	fmt.Println("- end update asset")
	return nil, nil
}