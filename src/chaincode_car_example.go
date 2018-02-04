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
	"errors"
	"strings"
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

type TransArgs struct {
	ShipID string `json:"ship_id"`
	Sender string `json:"sender"`
	Amount string `json:"amount"`
	Type string `json:"type"`
	Receiver string `json:"receiver"`
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
		IdHolder{Id:"ID0", Balance:"98053", Name: "Philipp der Erste"},
		IdHolder{Id:"ID1", Balance:"235", Name: "Philipp der Zweite"},
		IdHolder{Id:"ID2", Balance:"5534", Name: "Philipp der Dritte"},
		IdHolder{Id:"ID3", Balance:"2366", Name: "Philipp der Vierte"},
		IdHolder{Id:"ID4", Balance:"334", Name: "Philipp der Erste Junior"},
		IdHolder{Id:"ID5", Balance:"542", Name: "Philipp S."},

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
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	shipID := args[0]
	shipAsBytes , err := APIstub.GetState(shipID)
	if err != nil {
		return shim.Error(fmt.Sprintf("Couldn't find Shipment. Error: %s " , err.Error()))
	}

	ship := new(Shipment)
	err = json.Unmarshal(shipAsBytes, ship)
	ship.Status = args[1]
	if ship.Carrier == "" {
		ship.Carrier = args[2]
	}
	if ship.Status == "Accepted"{
		transArgs := TransArgs{Receiver: "", ShipID: ship.Id, Sender: ship.Carrier, Amount: ship.Price, Type: "deposite"}
		err =  s.createTransaction(APIstub, transArgs)
		if err != nil {
			return shim.Error(err.Error())
		}
	} else if ship.Status == "Approved"{
		transArgs := TransArgs{Receiver: ship.Carrier, ShipID: ship.Id, Sender: ship.Recipient, Amount: ship.Price, Type: "payment"}
		err = s.createTransaction(APIstub, transArgs)
		if err != nil {
			return shim.Error(err.Error())
		}
		trans := s.queryTrans(APIstub, ship.Id, "deposite")
		if trans == nil {
			return shim.Error("Alles scheiße!")
		}
		s.updateTransaction(APIstub, transArgs, trans.Id)
	} else if ship.Status == "not delivered" {
		transArgs := TransArgs{Receiver: ship.Retailer, ShipID: ship.Id, Sender: ship.Carrier, Amount: ship.Price, Type: "deposite"}
		trans := s.queryTrans(APIstub, ship.Id, "deposite")
		if trans == nil {
			return shim.Error("Alles kacke!")
		}
		s.updateTransaction(APIstub, transArgs, trans.Id)
	}
	ship.StatusUpdateTime = time.Now().String()
	ship.StatusChanger = args[2]
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

func (s *SmartContract) createTransaction (APIstub shim.ChaincodeStubInterface, args TransArgs) error {

	var id int = 0
	var transString string = "TRANS"
	var transID string = ""
	var stringID string = ""


	startKey := "TRANS000"
	endKey := "TRANS999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return err
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


	trans := Transaction{Id:transID, Amount: args.Amount, Receiver: args.Receiver, Sender: args.Sender, ShipId: args.ShipID, Type: args.Type}


	if trans.Receiver != "" {
		var argument = s.changeBalance(APIstub, trans)
		if argument != nil {
			return argument
		}
	}
	transAsBytes , err := json.Marshal(trans)
	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't marshal Transaction. Error: %s " , err.Error()))
	}

	if err := APIstub.PutState(transID, transAsBytes); err != nil{
		return err
	}else {
		return nil
	}
}

func (s *SmartContract) updateTransaction(APIstub shim.ChaincodeStubInterface, args TransArgs, id string) peer.Response  {

	transID := id
	transAsBytes , err := APIstub.GetState(transID)
	if err != nil {
		return shim.Error(fmt.Sprintf("Couldn't find Transaction. Error: %s " , err.Error()))
	}

	trans := new(Transaction)
	json.Unmarshal(transAsBytes, trans)

	trans.Receiver = args.Receiver


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

func (s *SmartContract) queryTrans(APIstub shim.ChaincodeStubInterface, shipID string, typ string) *Transaction{
	startKey := "TRANS000"
	endKey := "TRANS999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return nil
	}
	defer resultsIterator.Close()

	if resultsIterator.HasNext() != true {
		return nil
	}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil
		}
		trans := new(Transaction)
		err = json.Unmarshal(queryResponse.Value, trans)
		if err != nil {
			return nil
		}


		if trans.ShipId == shipID && trans.Type == typ{
			s.changeBalance(APIstub, *trans)
			return trans

		}
	}
	return nil
}

func (s *SmartContract) changeBalance(APIstub shim.ChaincodeStubInterface, trans Transaction) error{
	startKey := "ID0"
	endKey := "ID999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return errors.New("Blah1")
	}
	defer resultsIterator.Close()

	if resultsIterator.HasNext() != true {
		return errors.New("Blah2")
	}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return errors.New("Blah3")
		}
		idHolder := new(IdHolder)
		err = json.Unmarshal(queryResponse.Value, idHolder)
		if err != nil {
			return errors.New("Blah4")
		}
		//return errors.New("idHolder.Id: " + idHolder.Id + " trans.Receiver: " + trans.Receiver + " trans.Sender: " + trans.Sender)
		if strings.Compare(idHolder.Id, trans.Receiver) == 0{

			bal := idHolder.Balance

			balFl, err :=  strconv.ParseFloat(bal, 64)
			if err != nil {
				return errors.New("Parsing Error 1")
			}

			amountFl, err := strconv.ParseFloat(trans.Amount, 64)

			if err != nil {
				return errors.New("Parsing Error 2")
			}

			balFl = balFl + amountFl


			idHolder.Balance = strconv.FormatFloat(balFl, 'f', 2, 64)

			idHolderAsByte, _ := json.Marshal(idHolder)

			err = APIstub.PutState(idHolder.Id, idHolderAsByte)
			if err != nil {
				return err
			}

		}else if strings.Compare(idHolder.Id, trans.Sender) == 0 {
			bal := idHolder.Balance

			balFl, err :=  strconv.ParseFloat(bal, 64)

			if err != nil {
				return errors.New("Parsing Error 3")
			}

			amountFl, err := strconv.ParseFloat(trans.Amount, 64)

			if err != nil {
				return errors.New("Parsing Error 4")
			}

			balFl = balFl - amountFl

			idHolder.Balance = strconv.FormatFloat(balFl, 'f', 2, 64)

			idHolderAsByte, _ := json.Marshal(idHolder)

			err = APIstub.PutState(idHolder.Id, idHolderAsByte)
			if err != nil {
				return err
			}
		}
	}
	return nil
}



/// / The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}