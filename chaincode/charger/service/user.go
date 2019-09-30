package service

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/pkg/errors"
)

type User struct {
	MSP     string `json:"msp,omitempty"`
	Email   string `json:"email,omitempty"`
	Balance int    `json:"balance,omitempty"`
	State   string `json:"state,omitempty"`
}

func GetUser(stub shim.ChaincodeStubInterface, msp, userName string) (*User, error) {
	args := [][]byte{[]byte("QueryById"), []byte(msp), []byte(userName)}
	res := stub.InvokeChaincode("users", args, "default")

	if res.Status != 200 {
		return nil, errors.Errorf("User %s:%s does not exist: %s", msp, userName, res.Message)
	}

	var user User
	err := json.Unmarshal(res.Payload, &user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response from users::QueryById function")
	}

	return &user, nil
}
