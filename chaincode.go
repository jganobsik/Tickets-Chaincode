/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Init()", fcn, params)
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Invoke()", fcn, params)

	if fcn == "initTicket" { //create a new ticket
		return cc.initTicket(stub, args)
	} else if fcn == "transferTicket" { //change owner of a ticket
		return cc.transferTicket(stub, args)
	} else if fcn == "readTicket" { //read ticket
		return cc.readTicket(stub, args)
	} else if fcn == "redeemTicket" { //redeem ticket
		return cc.redeemTicket(stub, args)
	} else if fcn == "deleteTicket" { //delete ticket
		return cc.deleteTicket(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")

}

type TicketsChaincode struct {
}

type ticket struct {
	ObjectType string `json:"docType"`
	TicketID   string `json:"ticketId"`
	EventName  string `json:"eventName"`
	Location   string `json:"location"`
	EventDate  string `json:"eventDate"`
	Holder     string `json:"holder"`
	Redeemed   bool   `json:"redeemed"`
}

func (cc *TicketsChaincode) initTicket(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error

	// check for proper # of argument
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init ticket")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}

	ticketID := args[0]
	eventName := strings.ToLower(args[1])
	location := strings.ToLower(args[2])
	eventDate, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("2nd argument must be a numeric string")
	}
	holder := strings.ToLower(args[4])

	// ==== Check if ticket already exists ====
	ticketAsBytes, err := stub.GetState(ticketID)
	if err != nil {
		return shim.Error("Failed to get ticket: " + err.Error())
	} else if ticketAsBytes != nil {
		return shim.Error("This ticket already exists: " + ticketID)
	}

	// ==== Create ticket object and marshal to JSON ====
	objectType := "ticket"
	ticket := &ticket{objectType, ticketID, eventName, location, eventDate, holder}
	ticketJSONasBytes, err := json.Marshal(ticket)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save ticket to state ===
	err = stub.PutState(ticketID, ticketJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Ticket saved and indexed. Return success ====
	fmt.Println("- end init ticket")
	return shim.Success(nil)
}
