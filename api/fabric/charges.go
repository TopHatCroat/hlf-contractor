package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/pkg/errors"
	"time"
)

type ChargeTransaction struct {
	Id string    `json:"id,omitempty"`
	Contractor string    `json:"contractor,omitempty"`
	ChargeId   string    `json:"charge_id,omitempty"`
	User       string    `json:"user_email,omitempty"`
	Price      int       `json:"price,omitempty"`
	StartTime  time.Time `json:"start_date,omitempty"`
	EndTime    time.Time `json:"end_date,omitempty"`
	State      string    `json:"state,omitempty"`
}

type ChargePrice struct {
	Price int `json:"price,omitempty"`
}

type StartTransaction struct {
	Contractor string `json:"contractor,omitempty"`
}

type StopTransaction struct {
	Contractor string `json:"contractor,omitempty"`
	ChargeId   string `json:"charge_id,omitempty"`
}

type CompleteTransaction StopTransaction

func (c *Client) AllCharges(identity *shared.Identity, chargeProvider string) ([]ChargeTransaction, error) {
	channelName := "default"

	channelClient, err := c.GetChannelClient(identity, channelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	res, err := channelClient.Query(channel.Request{
		ChaincodeID: "charger",
		Fcn:         "QueryAll",
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from charger::QueryAll function")
	}

	var charges []ChargeTransaction
	err = json.Unmarshal(res.Payload, &charges)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response from charger::QueryAll function")
	}

	for i := range charges {
		charges[i].Id = fmt.Sprintf("%s:%s", charges[i].Contractor, charges[i].ChargeId)
	}

	return charges, nil
}

func (c *Client) FindChargeById(identity *shared.Identity, chargeProvider, chargeId string) (*ChargeTransaction, error) {
	channelName := "default"

	channelClient, err := c.GetChannelClient(identity, channelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	args := [][]byte{[]byte(chargeProvider), []byte(chargeId)}
	res, err := channelClient.Query(channel.Request{
		ChaincodeID: "charger",
		Fcn:         "QueryById",
		Args:        args,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from charger::QueryById function")
	}

	var chargeTransaction ChargeTransaction
	err = json.Unmarshal(res.Payload, &chargeTransaction)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response from charger::QueryById function")
	}

	chargeTransaction.Id = fmt.Sprintf("%s:%s", chargeTransaction.Contractor, chargeTransaction.ChargeId)

	return &chargeTransaction, nil
}

func (c *Client) StartCharge(identity *shared.Identity, chargeProvider string) (*ChargeTransaction, error) {
	channelName := "default"

	channelClient, err := c.GetChannelClient(identity, channelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	startChargeRequest := &StartTransaction{Contractor: chargeProvider}
	startChargeBytes, err := json.Marshal(startChargeRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	args := [][]byte{startChargeBytes}
	res, err := channelClient.Execute(channel.Request{
		ChaincodeID: "charger",
		Fcn:         "InvokeStartChargeTransaction",
		Args:        args,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from charger::InvokeStartChargeTransaction function")
	}

	var chargeTransaction ChargeTransaction
	err = json.Unmarshal(res.Payload, &chargeTransaction)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response from charger::InvokeStartChargeTransaction function")
	}

	chargeTransaction.Id = fmt.Sprintf("%s:%s", chargeTransaction.Contractor, chargeTransaction.ChargeId)

	return &chargeTransaction, nil
}

func (c *Client) StopCharge(identity *shared.Identity, chargeProvider, chargeId string) (*ChargeTransaction, error) {
	channelName := "default"

	channelClient, err := c.GetChannelClient(identity, channelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	stopChargeRequest := &StopTransaction{Contractor: chargeProvider, ChargeId: chargeId}
	stopChargeBytes, err := json.Marshal(stopChargeRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	args := [][]byte{stopChargeBytes}
	res, err := channelClient.Execute(channel.Request{
		ChaincodeID: "charger",
		Fcn:         "InvokeStopChargeTransaction",
		Args:        args,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from charger::InvokeStopChargeTransaction function")
	}

	var chargeTransaction ChargeTransaction
	err = json.Unmarshal(res.Payload, &chargeTransaction)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response from charger::InvokeStopChargeTransaction function")
	}

	chargeTransaction.Id = fmt.Sprintf("%s:%s", chargeTransaction.Contractor, chargeTransaction.ChargeId)

	return &chargeTransaction, nil
}

func (c *Client) CompleteCharge(identity *shared.Identity, chargeProvider, chargeId string) (*ChargeTransaction, error) {
	channelName := "default"

	channelClient, err := c.GetChannelClient(identity, channelName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	completeChargeRequest := &CompleteTransaction{Contractor: chargeProvider, ChargeId: chargeId}
	completeChargeBytes, err := json.Marshal(completeChargeRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute chaincode call")
	}

	args := [][]byte{completeChargeBytes}
	res, err := channelClient.Execute(channel.Request{
		ChaincodeID: "charger",
		Fcn:         "InvokeCompleteChargeTransaction",
		Args:        args,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from charger::InvokeCompleteChargeTransaction function")
	}

	var chargeTransaction ChargeTransaction
	err = json.Unmarshal(res.Payload, &chargeTransaction)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response from charger::InvokeCompleteChargeTransaction function")
	}

	chargeTransaction.Id = fmt.Sprintf("%s:%s", chargeTransaction.Contractor, chargeTransaction.ChargeId)

	return &chargeTransaction, nil
}
