/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
"bytes"
"encoding/json"
"fmt"
"strconv"

"github.com/hyperledger/fabric/core/chaincode/shim"
"github.com/hyperledger/fabric/protos/peer"
"time"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Content struct {
	Amount int64 `json:"amount"`
	Description string `json:"description"`
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Shipment struct {
	Id string `json:"id"`
	CreaterId string `json:"creater_id"`
	Status string `json:"status"`
	StatusUpdateTime string `json:"status_update_time"`
	StatusChanger string `json:"status_changer"`
	Carrier string `json:"carrier"`
	Recipient string `json:"recipient"`
	Retailer string `json:"retailer"`
	Price string `json:"price"`
	PickUp string `json:"pick_up"`
	Destination string `json:"destination"`
	ContentList string `json:"content_list"`
	Space string `json:"space"`
	Startpoint string `json:"startpoint"`
	Endpoint string `json:"endpoint"`
}

type IdHolder struct {
	Id string `json:"id"`
	Balance string `json:"balance"`
	Name string `json:"name"`
}

type Transaction struct {
	Id string `json:"id"`
	Receiver string `json:"receiver"`
	Sender string `json:"sender"`
	Amount string `json:"amount"`
	Type string `json:"type"`
	ShipId string `json:"ship_id"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createShipment" {
		return s.createShipment(APIstub, args)
	} else if function == "updateStatus" {
		return s.updateStatus(APIstub, args)
	}else if function == "queryAllShips" {
		return s.queryAllShips(APIstub)
	}else if function == "queryId" {
		return s.queryId(APIstub, args)
	}else if function == "initLedger" {
		return s.initLedger(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) peer.Response {

	idHolders := []IdHolder{
		IdHolder{Id:"1", Balance:"500", Name: "Philipp der Erste"},
		IdHolder{Id:"2", Balance:"500", Name: "Philipp der Zweite"},
		IdHolder{Id:"3", Balance:"500", Name: "Philipp der Dritte"},
		IdHolder{Id:"4", Balance:"500", Name: "Philipp der Vierte"},
		IdHolder{Id:"5", Balance:"500", Name: "Philipp der Erste Junior"},
		IdHolder{Id:"6", Balance:"500", Name: "Philipp S."},

	}

	i := 0
	for i < len(idHolders) {
		idsAsBytes, err := json.Marshal(idHolders[i])
		if err != nil {
			return shim.Error(err.Error())
		}
		APIstub.PutState("ID"+strconv.Itoa(i), idsAsBytes)
		fmt.Println("Added ", idHolders[i])
		i = i + 1

	}


	return shim.Success(nil)
}

func (s *SmartContract) createShipment(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	var id int = 0
	var shipString string = "SHIP"
	var shipID string = ""
	var stringID string = ""


	startKey := "SHIP000"
	endKey := "SHIP999"

	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments. Expecting 11")
	}
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf(string(queryResponse.Key))
		id = id + 1
	}
	if id < 10 {
		shipString = "SHIP00"
	} else if id < 100 {
		shipString = "SHIP0"
	}

	stringID = strconv.Itoa(id)
	shipID = shipString + stringID

	currentTime := time.Now().String()


	var ship = Shipment{Id: shipID,CreaterId: args[0], Status: "Created",  StatusUpdateTime: currentTime, StatusChanger: args[0], Carrier: args[1], Recipient: args[2], Retailer: args[3], Price: args[4], PickUp: args[5], Destination: args[6], ContentList: args[7], Space: args[8], Startpoint: args[9], Endpoint:args[10]}

	shipAsBytes , err := json.Marshal(ship);
	if err != nil {
		return shim.Error(fmt.Sprintf("Couldn't marshal Shipment. Error: %s " , err.Error()))
	}

	if err := APIstub.PutState(shipID, shipAsBytes); err != nil{
		return shim.Error(err.Error())
	}else {
		return shim.Success(shipAsBytes)
	}
}

func (s *SmartContract) updateStatus(APIstub shim.ChaincodeStubInterface, args []string) peer.Response  {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	shipID := args[0]
	shipAsBytes , err := APIstub.GetState(shipID)
	if err != nil {
		return shim.Error(fmt.Sprintf("Couldn't find Shipment. Error: %s " , err.Error()))
	}

	ship := Shipment{}
	json.Unmarshal(shipAsBytes, &ship)
	ship.Status = args[1]
	ship.StatusUpdateTime = time.Now().String()
	ship.StatusChanger = args[2]
	if ship.Carrier == "" {
		ship.Carrier = args[2]
	}
	if  args[3] == ""{
		ship.Space = args[3]
	}

	shipAsBytes, err = json.Marshal(ship)
	if err != nil {
		return shim.Error(fmt.Sprintf("Couldn't marshal Shipment. Error: %s " , err.Error()))
	}

	if err := APIstub.PutState(shipID, shipAsBytes); err != nil{
		return shim.Error(err.Error())
	}else {
		return shim.Success(shipAsBytes)
	}
}

func (s *SmartContract) queryAllShips(APIstub shim.ChaincodeStubInterface) peer.Response {

	startKey := "SHIP000"
	endKey := "SHIP999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllShipments:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryId(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1{
		return shim.Error("Philipp hats verbockt -.-")
	}

	idAsBytes,_ := APIstub.GetState(args[0])

	return shim.Success(idAsBytes)
}

func (s *SmartContract) createTransaction (APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	var id int = 0
	var transString string = "TRANS"
	var transID string = ""
	var stringID string = ""


	startKey := "TRANS000"
	endKey := "TRANS999"

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf(string(queryResponse.Key))
		id = id + 1
	}
	if id < 10 {
		transString = "TRANS00"
	} else if id < 100 {
		transString = "TRANS0"
	}

	stringID = strconv.Itoa(id)
	transID = transString + stringID


	var trans = Transaction{Id:transID, Amount: args[0], Receiver: args[1], Sender: args[2], ShipId: args[3], Type: args[4]}

	transAsBytes , err := json.Marshal(trans);
	if err != nil {
		return shim.Error(fmt.Sprintf("Couldn't marshal Transaction. Error: %s " , err.Error()))
	}

	if err := APIstub.PutState(transID, transAsBytes); err != nil{
		return shim.Error(err.Error())
	}else {
		return shim.Success(transAsBytes)
	}
}

func (s *SmartContract) updateTransaction(APIstub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	transID := args[0]
	transAsBytes , err := APIstub.GetState(transID)
	if err != nil {
		return shim.Error(fmt.Sprintf("Couldn't find Transaction. Error: %s " , err.Error()))
	}

	trans := Transaction{}
	json.Unmarshal(transAsBytes, &trans)

	trans.Receiver = args[1]


	transAsBytes, err = json.Marshal(trans)
	if err != nil {
		return shim.Error(fmt.Sprintf("Couldn't marshal Transaction. Error: %s " , err.Error()))
	}

	if err := APIstub.PutState(transID, transAsBytes); err != nil{
		return shim.Error(err.Error())
	}else {
		return shim.Success(transAsBytes)
	}
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}