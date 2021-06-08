package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"strings"
)

type SimpleChaincode struct {
}

type opinionRecord struct {
	DocType     string `json:"docType"`
	Id          string `json:"id"`
	Department  string `json:"department"`
	Name        string `json:"name"`
	Object      string `json:"object"`
	Type        string `json:"type"`
	OpinionTime string `json:"opinionTime"`
	DoneTime    string `json:"doneTime"`
	Content     string `json:"content"`
}

type reviewRecord struct {
	DocType    string `json:"docType"`
	Id         string `json:"id"`
	Department string `json:"department"`
	Name       string `json:"name"`
	Object     string `json:"object"`
	From       string `json:"from"`
	ReviewTime string `json:"reviewTime"`
	Result     string `json:"result"`
}

type userRecord struct {
	DocType    string `json:"docType"`
	Id         string `json:"id"`
	Department string `json:"department"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Content    string `json:"content"`
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple Chaincode: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Invoking is running" + function)
	switch function {
	case "createOpinionRecord":
		return t.createOpinionRecord(stub, args)
	case "createReviewRecord":
		return t.createReviewRecord(stub, args)
	case "createUserRecord":
		return t.createUserRecord(stub, args)
	case "modifyOpinionRecord":
		return t.modifyOpinionRecord(stub, args)
	case "modifyUserRecord":
		return t.modifyUserRecord(stub, args)
	case "queryWithQueryString":
		return t.queryWithQueryString(stub, args)
	case "queryRecordById":
		return t.queryRecordById(stub, args)
	case "queryRecordByObject":
		return t.queryRecordByObject(stub, args)
	case "queryRecordByUser":
		return t.queryRecordByUser(stub, args)
	default:
		return shim.Error("Received unknown function invocation")
	}
}

func (t *SimpleChaincode) createOpinionRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 6 && len(args) != 8 {
		return shim.Error("Incorrect number of arguments: " +
			"expecting 6 for opinions need to be reviewed; " +
			"expecting 8 for opinions without being reviewed")
	}

	for i := 0; i < len(args); i += 1 {
		if len(args[i]) == 0 {
			errMessage := fmt.Sprintf("Number %d member of the arguments must be a non-empty string", i)
			return shim.Error(errMessage)
		}
	}

	docType := "opinionRecord"
	id := strings.ToLower(args[0])
	department := strings.ToLower(args[1])
	name := strings.ToLower(args[2])
	object := strings.ToLower(args[3])
	_type := strings.ToLower(args[4])
	opinionTime := strings.ToLower(args[5])
	var doneTime, content string
	doneTime = ""
	content = ""

	if len(args) == 8 {
		doneTime = strings.ToLower(args[6])
		content = strings.ToLower(args[7])
	}

	dataRecordJsonAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get opinion record: " + err.Error())
	} else if dataRecordJsonAsBytes != nil {
		return shim.Error("Duplicated opinion record found: " + id)
	}

	dataRecord := &opinionRecord{
		docType,
		id,
		department,
		name,
		object,
		_type,
		opinionTime,
		doneTime,
		content,
	}

	dataRecordJsonAsBytes, err = json.Marshal(dataRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(id, dataRecordJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) createReviewRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments: expecting 7.")
	}

	for i := 0; i < len(args); i += 1 {
		if len(args[i]) == 0 {
			errMessage := fmt.Sprintf("Number %d member of arguments must be a non-empty string", i)
			return shim.Error(errMessage)
		}
	}

	docType := "reviewRecord"
	id := strings.ToLower(args[0])
	department := strings.ToLower(args[1])
	name := strings.ToLower(args[2])
	object := strings.ToLower(args[3])
	from := strings.ToLower(args[4])
	reviewTime := strings.ToLower(args[5])
	result := strings.ToLower(args[6])

	reviewRecordJsonAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get review record: " + err.Error())
	} else if reviewRecordJsonAsBytes != nil {
		return shim.Error("Duplicated review record found: " + id)
	}

	reviewRecord := &reviewRecord{
		docType,
		id,
		department,
		name,
		object,
		from,
		reviewTime,
		result,
	}

	reviewRecordJsonAsBytes, err = json.Marshal(reviewRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(id, reviewRecordJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) createUserRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments for user, expecting 5")
	}
	for i := 0; i < 5; i += 1 {
		if len(args[i]) <= 0 {
			errMessage := fmt.Sprintf("No.%d argument must be a non-empty string", i)
			return shim.Error(errMessage)
		}
	}

	docType := "userRecord"
	id := strings.ToLower(args[0])
	department := strings.ToLower(args[1])
	name := strings.ToLower(args[2])
	role := strings.ToLower(args[3])
	content := strings.ToLower(args[4])

	userJsonAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userJsonAsBytes != nil {
		return shim.Error("Duplicated user found")
	}

	userRecord := &userRecord{
		docType,
		id,
		department,
		name,
		role,
		content,
	}

	userJsonAsBytes, err = json.Marshal(userRecord)
	if err != nil {
		shim.Error(err.Error())
	}
	err = stub.PutState(id, userJsonAsBytes)
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) modifyOpinionRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments: expecting 3")
	}
	id := strings.ToLower(args[0])
	opinionRecordJsonAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get opinion record: " + err.Error())
	} else if opinionRecordJsonAsBytes == nil {
		return shim.Error("opinion record does not exist")
	}

	opinionRecordInstance := opinionRecord{}
	err = json.Unmarshal(opinionRecordJsonAsBytes, &opinionRecordInstance)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(opinionRecordInstance.DoneTime) != 0 || len(opinionRecordInstance.Content) != 0 {
		return shim.Error("Unable to modify a opinion record that has been done")
	}

	doneTime := strings.ToLower(args[1])
	content := strings.ToLower(args[2])
	opinionRecordInstance.DoneTime = doneTime
	opinionRecordInstance.Content = content

	opinionRecordJsonAsBytes, _ = json.Marshal(opinionRecordInstance)
	err = stub.PutState(id, opinionRecordJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) modifyUserRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments, expecting 5")
	}
	id := strings.ToLower(args[0])
	department := strings.ToLower(args[1])
	name := strings.ToLower(args[2])
	role := strings.ToLower(args[3])
	content := strings.ToLower(args[4])

	userRecordJsonAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get user record: " + err.Error())
	} else if userRecordJsonAsBytes == nil {
		return shim.Error("user record does not exist")
	}
	userRecordInstance := userRecord{}
	err = json.Unmarshal(userRecordJsonAsBytes, &userRecordInstance)
	if err != nil {
		return shim.Error(err.Error())
	}
	userRecordInstance.Name = name
	userRecordInstance.Department = department
	userRecordInstance.Role = role
	userRecordInstance.Content = content

	userRecordJsonAsBytes, err = json.Marshal(userRecordInstance)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(id, userRecordJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("{\"list\":[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]}")

	return &buffer, nil
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		_ = resultsIterator.Close()
	}(resultsIterator)

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *SimpleChaincode) queryWithQueryString(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect arguments, must be a query string")
	}
	queryString := args[0]
	queryResult, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResult)
}

func (t *SimpleChaincode) queryRecordById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}
	id := strings.ToLower(args[0])
	recordJsonAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get record: " + err.Error())
	} else if recordJsonAsBytes == nil {
		return shim.Error("Record does not exist")
	}
	return shim.Success(recordJsonAsBytes)
}

func (t *SimpleChaincode) queryRecordByObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments, expecting 2")
	}
	docType := args[0]
	object := strings.ToLower(args[1])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"object\":\"%s\"}}", docType, object)
	queryResult, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResult)
}

func (t *SimpleChaincode) queryRecordByUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments, expecting 3")
	}
	docType := args[0]
	department := strings.ToLower(args[1])
	name := strings.ToLower(args[2])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"department\":\"%s\",\"name\":\"%s\"}}", docType, department, name)
	queryResult, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResult)
}
