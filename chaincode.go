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
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the vote structure, with 5 properties.  Structure tags are used by encoding/json library
type Asset_Votes struct {
	Security   string `json:"security"`
	Factor  string `json:"factor"`
	Vote string `json:"vote"`
	OwnerId  string `json:"ownerid"`
	OwnerDesc string `json:"ownerdesc"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryVoter" {
		return s.queryVoter(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAllVotes" {
		return s.queryAllVotes(APIstub)
	} else if function == "doVoting" {
		return s.doVoting(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryVoter(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	voterAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(voterAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	voters := []Asset_Votes{
	      	Asset_Votes{Security: "PIN", Factor: "100", Vote: "SINVOTAR", OwnerId: "12345", OwnerDesc: "Homero Simpson"},
      		Asset_Votes{Security: "PIN", Factor: "120", Vote: "SINVOTAR", OwnerId: "67890", OwnerDesc: "Peter Griffin"},
      		Asset_Votes{Security: "PIN", Factor: "200", Vote: "SINVOTAR", OwnerId: "654321", OwnerDesc: "Jhon Smith"},
      		Asset_Votes{Security: "PIN", Factor: "1000", Vote: "SINVOTAR", OwnerId: "098765", OwnerDesc: "Chaparron Bonaparte"},
      		Asset_Votes{Security: "PIN", Factor: "50", Vote: "SINVOTAR", OwnerId: "A23421", OwnerDesc: "Sun Wukong"},
      		Asset_Votes{Security: "PIN", Factor: "80", Vote: "SINVOTAR", OwnerId: "98765", OwnerDesc: "Seiya Shiryu"},
      		Asset_Votes{Security: "PIN", Factor: "250", Vote: "SINVOTAR", OwnerId: "765890", OwnerDesc: "Ned Flanders"},
      		Asset_Votes{Security: "PIN", Factor: "40", Vote: "SINVOTAR", OwnerId: "777888", OwnerDesc: "Carlos Donoso"},
      		Asset_Votes{Security: "PIN", Factor: "5", Vote: "SINVOTAR", OwnerId: "877899", OwnerDesc: "Ramon Valdes"},
      		Asset_Votes{Security: "PIN", Factor: "500", Vote: "SINVOTAR", OwnerId: "131313", OwnerDesc: "Roberto Gomez"},
	}

	i := 0
	for i < len(voters) {
		fmt.Println("i is ", i)
		voterAsBytes, _ := json.Marshal(voters[i])
		APIstub.PutState("VOTER"+strconv.Itoa(i), voterAsBytes)
		fmt.Println("Added", voters[i])
		i = i + 1
	}

	return shim.Success(nil)
}


func (s *SmartContract) queryAllVotes(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "VOTER0"
	endKey := "VOTER999"

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

	fmt.Printf("- queryAllVotes:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) doVoting(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	voterAsBytes, _ := APIstub.GetState(args[0])
	vote := Asset_Votes{}



	json.Unmarshal(voterAsBytes, &vote)
	if (vote.Vote == "SINVOTAR"){
		fmt.Println("You havent voted before, vote accepted!")
		vote.Vote = args[1]
	}else{
		return shim.Error("Already voted before!, you can only vote once")
	}


	voterAsBytes, _ = json.Marshal(vote)
	APIstub.PutState(args[0], voterAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
