package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	if err != nil {
		fmt.Println("Failed to find asset by id " + id)
		return nil, errors.New(err.Error())
	}

	// remove the asset
	err = stub.DelState(id) //remove the key from chaincode state
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

	if len(args) < 16 {
		return nil, errors.New("Incorrect number of arguments. Expecting at least 16")
	}

	id := args[0]
	//check if asset id already exists
	_, err = get_asset(stub, id)
	if err == nil {
		fmt.Println("This asset already exists - " + id)
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
			"moveDateTime": "` + args[6] + `",
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
			}
		}			
	}`

	fmt.Println("Input PharmaAsset Object - " + str)
	var pharmaAsset PharmaAsset
	err = json.Unmarshal([]byte(str), &pharmaAsset)
	if err != nil {
		fmt.Println("Error while unmarshalling " + err.Error())
		return nil, errors.New(err.Error())
	}
	//fmt.Printf("PharmaAsset Object after un-marshalling:\n%s", pharmaAsset)
	if len(args) > 16 {
		for i := 16; i < len(args); i = i + 2 {
			var child AssetChildren
			child.AssetId = args[i]
			child.AssetType = args[i+1]
			pharmaAsset.AssetData.Children = append(pharmaAsset.AssetData.Children, child)
		}
	}

	inputByteStr, err := json.Marshal(pharmaAsset)
	if err != nil {
		fmt.Println("Error while marshalling " + err.Error())
		return nil, errors.New(err.Error())
	}

	err = stub.PutState(id, inputByteStr) //store asset with id as key
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
	assetAsBytes, err := stub.GetState(id) //getState retreives a key/value from the ledger
	if err != nil {                        //this seems to always succeed, even if key didn't exist
		return nil, errors.New("Failed to find asset - " + id)
	}
	json.Unmarshal(assetAsBytes, &asset) //un stringify it aka JSON.parse()

	if asset.AssetId != id { //test if marble is actually here or just nil
		return nil, errors.New("Asset does not exist - " + id)
	}

	var traceData AssetTraceData
	traceData.Owner = args[1]
	traceData.Status = args[2]
	traceData.MoveDateTime = args[3]
	traceData.Location = args[4]
	traceData.GeoLocation = args[5]
	asset.AssetTraceData = append(asset.AssetTraceData, traceData)
	fmt.Println("PharmaAsset object post update of trace data - %s",asset)

	inputByteStr, err := json.Marshal(asset)
	if err != nil {
		fmt.Println("Error while marshalling in update " + err.Error())
		return nil, errors.New(err.Error())
	}

	err = stub.PutState(id, inputByteStr) //store asset with id as key
	if err != nil {
		return nil, errors.New(err.Error())
	}

	fmt.Println("- end update asset")
	return nil, nil
}
