package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/router"
)

func CreateChaincode() *router.Chaincode {
	r := router.New("root")

	r.Init(Init)

	// method for debug chaincode state
	debug.AddHandlers(r, "debug")

	CreateRouter(r)

	return router.NewChaincode(r)
}

func main() {
	chaincode := CreateChaincode()
	if err := shim.Start(chaincode); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
