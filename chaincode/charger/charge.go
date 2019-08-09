package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	"math/rand"
	"strconv"
	"time"
)

func QueryById(c router.Context) (interface{}, error) {
	var (
		contractor = c.ParamString("contractor")
		chargeId   = c.ParamString("charge_id")
	)

	return c.State().Get(&Entity{
		Contractor: contractor,
		ChargeId:   chargeId,
	})
}

func QueryAll(c router.Context) (interface{}, error) {
	return c.State().List(TypeName, &Entity{})
}

func InvokeStartChargeTransaction(c router.Context) (res interface{}, err error) {
	var startTransaction = c.Param("startTransaction").(StartTransaction) // Assert the chaincode parameter

	// Validate current owner
	user, err := identity.FromStub(c.Stub())
	if user == nil || user.GetSubject() == "" {
		return nil, fmt.Errorf("missing identity email for this user")
	}

	// TODO: check if user can start transaction (has outstanding balance? already ongoing transaction?)
	var chargeTransaction = &Entity{
		Contractor: startTransaction.Contractor,
		ChargeId:   strconv.Itoa(rand.Int()),
		User:       SerializeIdentity(user.MspID, user.Cert),
		StartTime:  time.Now(),
		State:      ChargeStateStarted, // Initial state
	}

	return chargeTransaction, c.State().Insert(chargeTransaction)
}

func InvokeStopChargeTransaction(c router.Context) (interface{}, error) {
	var (
		chargeTransaction Entity
		// Buy transaction payload
		data = c.Param("stopTransaction").(StopTransaction)

		// Get the current commercial paper state
		cp, err = c.State().Get(&Entity{
			Contractor: data.Contractor,
			ChargeId:   data.ChargeId,
		}, &Entity{})
	)

	if err != nil {
		return nil, errors.Wrap(err, "charge transaction not found")
	}

	chargeTransaction = cp.(Entity)

	// Validate current owner
	user, err := identity.FromStub(c.Stub())
	if user != nil && chargeTransaction.User != SerializeIdentity(user.MspID, user.Cert) {
		return nil, fmt.Errorf(
			"transaction %s %s is not owned by this user", chargeTransaction.Contractor, chargeTransaction.ChargeId)
	}

	// Check if transaction is in progress
	if chargeTransaction.State == ChargeStateStarted {
		chargeTransaction.State = ChargeStateStopped
		chargeTransaction.EndTime = time.Now()
	} else {
		return nil, fmt.Errorf("transaction can not be stopped from %s state", chargeTransaction.State)
	}

	// TODO: calculate price

	return chargeTransaction, c.State().Put(chargeTransaction)
}

func InvokeCompleteChargeTransaction(c router.Context) (interface{}, error) {
	var (
		chargeTransaction Entity
		// Buy transaction payload
		data = c.Param("completeTransaction").(CompleteTransaction)

		// Get the current commercial paper state
		cp, err = c.State().Get(&Entity{
			Contractor: data.Contractor,
			ChargeId:   data.ChargeId,
		}, &Entity{})
	)

	if err != nil {
		return nil, errors.Wrap(err, "charge transaction not found")
	}

	chargeTransaction = cp.(Entity)

	// TODO: Validate current invoker has this right

	// Check if transaction is in progress
	if chargeTransaction.State == ChargeStateStopped {
		chargeTransaction.State = ChargeStateCompleted
	} else {
		return nil, fmt.Errorf("transaction can not be stopped from %s state", chargeTransaction.State)
	}

	return chargeTransaction, c.State().Put(chargeTransaction)
}
