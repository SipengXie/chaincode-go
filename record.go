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

type dataRecord struct {
	DocType          string `json:"docType"`
	Id               string `json:"id"`
	Department       string `json:"department"`
	User             string `json:"user"`
	Object           string `json:"object"`
	Type             string `json:"type"`
	OpinionTime      string `json:"opinionTime"`
	Reviewer         string `json:"reviewer"`
	ReviewTime       string `json:"reviewTime"`
	ReviewResult     string `json:"reviewResult"`
	ReviewDepartment string `json:"reviewDepartment"`
	OperateTime      string `json:"operateTime"`
	Content          string `json:"content"`
}

type userRecord struct {
	DocType     string `json:"docType"`
	Department  string `json:"department"`
	UserName    string `json:"userName"`
	UserAddress string `json:"userAddress"`
	Role        string `json:"role"`
}

func (t *dataRecord) UnReviewed() bool {
	return t.Reviewer == "" && t.ReviewTime == "" && t.ReviewResult == ""
}

func (t *dataRecord) UnOperated() bool {
	return t.OperateTime == ""
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
	case "initDataRecord":
		return t.initDataRecord(stub, args)
	case "initUserRecord":
		return t.initUserRecord(stub, args)
	case "modifyDataRecord":
		return t.modifyDataRecord(stub, args)
	case "modifyUserRecord":
		return t.modifyUserRecord(stub, args)
	case "queryDataRecordByObject":
		return t.queryDataRecordByObject(stub, args)
	case "queryDataRecordById":
		return t.queryDataRecordById(stub, args)
	case "queryDataRecordByUser":
		return t.queryDataRecordByUser(stub, args)
	case "queryUserRecordByDept":
		return t.queryUserRecordByDept(stub, args)
	case "queryUserRecordByAddr":
		return t.queryUserRecordByAddr(stub, args)
	case "queryWithQueryString":
		return t.queryWithQueryString(stub, args)
	default:
		return shim.Error("Received unknown function invocation")
	}
}

func (t *SimpleChaincode) initDataRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0    1           2    3      4   5
	// uuid Department1 Jack object add opinionTime
	// 0    1           2    3      4   5           6
	// uuid Department1 Jack object add operateTime [contentHash]
	var err error
	if len(args) != 6 && len(args) != 7 {
		return shim.Error("Incorrect number of arguments, expecting 6 for records with opinion; expecting 7 for records without opinion")
	}
	// basic requirements of a data record
	for i := 0; i < len(args); i += 1 {
		if len(args[i]) == 0 {
			errMessage := fmt.Sprintf("Number %d member of arguments must be a non-empty string", i)
			return shim.Error(errMessage)
		}
	}

	docType := "dataRecord"
	uuid := strings.ToLower(args[0])
	department := strings.ToLower(args[1])
	user := strings.ToLower(args[2])
	object := strings.ToLower(args[3])
	_type := strings.ToLower(args[4])
	var opinionTime, operateTime, content string
	opinionTime = ""
	operateTime = ""
	content = ""

	if len(args) == 6 {
		opinionTime = strings.ToLower(args[5])
	} else {
		operateTime = strings.ToLower(args[5])
		content = strings.ToLower(args[6])
	}

	dataRecordJsonAsBytes, err := stub.GetState(uuid)
	if err != nil {
		return shim.Error("Failed to get record: " + err.Error())
	} else if dataRecordJsonAsBytes != nil {
		return shim.Error("Duplicated record found: " + uuid)
	}

	dataRecord := &dataRecord{
		docType,
		uuid,
		department,
		user,
		object,
		_type,
		opinionTime,
		"",
		"",
		"",
		"",
		operateTime,
		content,
	}

	dataRecordJsonAsBytes, err = json.Marshal(dataRecord)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(uuid, dataRecordJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) initUserRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//  0           1    2           3
	//  department1 Jack userAddress Role
	var err error
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments for user, expecting 5")
	}
	for i := 0; i < 4; i += 1 {
		if len(args[i]) <= 0 {
			errMessage := fmt.Sprintf("No.%d argument must be a non-empty string", i)
			return shim.Error(errMessage)
		}
	}

	docType := "userRecord"
	department := strings.ToLower(args[0])
	userName := strings.ToLower(args[1])
	userAddress := strings.ToLower(args[2])
	role := strings.ToLower(args[3])

	userJsonAsBytes, err := stub.GetState(userAddress)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userJsonAsBytes != nil {
		return shim.Error("Duplicated user found")
	}

	userRecord := &userRecord{
		docType,
		department,
		userName,
		userAddress,
		role,
	}
	userJsonAsBytes, err = json.Marshal(userRecord)
	if err != nil {
		shim.Error(err.Error())
	}
	err = stub.PutState(userAddress, userJsonAsBytes)
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) modifyDataRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0     1        2		     3             4
	// uuid  reviewer reviewTime reviewResult  reviewDepartment   for review
	// 0     1           2
	// uuid  operateTime [contentHash]   						  for operate
	var err error
	if len(args) != 5 && len(args) != 3 {
		return shim.Error("Incorrect number of arguments: expecting 5 for review, 3 for operate")
	}
	fmt.Println("Starting modifyDataRecord")
	uuid := strings.ToLower(args[0])

	dataRecordJsonAsBytes, err := stub.GetState(uuid)
	if err != nil {
		return shim.Error("Failed to get dataRecord: " + err.Error())
	} else if dataRecordJsonAsBytes == nil {
		return shim.Error("dataRecord does not exist")
	}

	dataRecordInstance := dataRecord{}
	err = json.Unmarshal(dataRecordJsonAsBytes, &dataRecordInstance)
	if err != nil {
		return shim.Error(err.Error())
	}

	if len(args) == 5 {
		if !(dataRecordInstance.UnReviewed() && dataRecordInstance.UnOperated()) {
			return shim.Error("The data record has been reviewed or operated, dataID: " + uuid)
		}
		reviewer := strings.ToLower(args[1])
		reviewerTime := strings.ToLower(args[2])
		reviewResult := strings.ToLower(args[3])
		reviewDepartment := strings.ToLower(args[4])
		dataRecordInstance.Reviewer = reviewer
		dataRecordInstance.ReviewTime = reviewerTime
		dataRecordInstance.ReviewResult = reviewResult
		dataRecordInstance.ReviewDepartment = reviewDepartment
	} else {
		if !dataRecordInstance.UnOperated() || dataRecordInstance.UnReviewed() {
			return shim.Error("The data record has been operated or has not been reviewed, dataID: " + uuid)
		}
		if dataRecordInstance.ReviewResult == "false" {
			return shim.Error("The opinion of this data record did not be approved, dataID: " + uuid)
		}
		operateTime := strings.ToLower(args[1])
		content := strings.ToLower(args[2])
		dataRecordInstance.OperateTime = operateTime
		dataRecordInstance.Content = content
	}
	dataRecordJsonAsBytes, _ = json.Marshal(dataRecordInstance)
	err = stub.PutState(uuid, dataRecordJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) modifyUserRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0          1        2           3
	// department username userAddress role
	var err error
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments, expecting 4")
	}
	department := strings.ToLower(args[0])
	userName := strings.ToLower(args[1])
	userAddress := strings.ToLower(args[2])
	role := strings.ToLower(args[3])

	userRecordJsonAsBytes, err := stub.GetState(userAddress)
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
	userRecordInstance.Department = department
	userRecordInstance.UserName = userName
	userRecordInstance.Role = role

	userRecordJsonAsBytes, err = json.Marshal(userRecordInstance)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(userAddress, userRecordJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

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

func (t *SimpleChaincode) queryDataRecordById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}
	uuid := strings.ToLower(args[0])
	dataRecordJsonAsBytes, err := stub.GetState(uuid)
	if err != nil {
		return shim.Error("Failed to get dataRecord: " + err.Error())
	} else if dataRecordJsonAsBytes == nil {
		return shim.Error("dataRecord does not exist")
	}
	return shim.Success(dataRecordJsonAsBytes)
}

func (t *SimpleChaincode) queryUserRecordByAddr(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}
	userAddress := strings.ToLower(args[0])
	userRecordJsonAsBytes, err := stub.GetState(userAddress)
	if err != nil {
		return shim.Error("Failed to get user record: " + err.Error())
	} else if userRecordJsonAsBytes == nil {
		return shim.Error("user record does not exist")
	}
	return shim.Success(userRecordJsonAsBytes)
}

func (t *SimpleChaincode) queryDataRecordByObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0
	// objectName
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}
	object := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"object\":\"%s\"}}", object)
	queryResult, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResult)
}

func (t *SimpleChaincode) queryDataRecordByUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0 		  1
	// department user
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments, expecting 2")
	}
	department := strings.ToLower(args[0])
	user := strings.ToLower(args[1])
	queryString := fmt.Sprintf("{\"selector\":{\"department\":\"%s\",\"user\":\"%s\"}}", department, user)
	queryResult, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResult)
}

func (t *SimpleChaincode) queryUserRecordByDept(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0
	// department
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}
	department := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"userRecord\",\"department\":\"%s\"}}", department)
	queryResult, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResult)
}
