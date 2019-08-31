/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
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

	if function == "initTicket" { //create a new ticket
		return cc.initTicket(stub, args)
	} else if function == "transferTicket" { //change owner of a ticket
		return cc.transferTicket(stub, args)
	} else if function == "readTicket" { //read ticket
		return cc.readTicket(stub, args)
	} else if function == "redeemTicket" { //redeem ticket
		return cc.redeemTicket(stub, args)
	} else if function == "deleteTicket" { //delete ticket
		return cc.deleteTicket(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")

}

type TicketsChaincode struct {
}

type ticket struct {
	ObjectType string `json:"docType"`
	EventName  string `json:"eventName"`
	Location   string `json:"location"`
	EventDate  string `json:"eventDate"`
	Holder     string `json:"holder"`
	Redeemed   bool   `json:"redeemed"`
}
