package main

import (
	"fmt"
	"github.com/TopHatCroat/hlf-contractor/chaincode/charger/charge"
	"github.com/TopHatCroat/hlf-contractor/chaincode/charger/price"
	"github.com/TopHatCroat/hlf-contractor/chaincode/charger/service"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

func Init(c router.Context) (interface{}, error) {
	_, err := owner.InvokeSetFromCreator(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize chaincode")
	}

	initialPrice := c.ParamString("initialPrice")

	if initialPrice != "" {
		chargePrice, err := strconv.Atoi(initialPrice)
		if err != nil {
			return nil, errors.Wrap(err, "failed to initialize chaincode")
		}

		if chargePrice <= 0 {
			return nil, errors.New("failed to initialize chaincode due to invalid priceEntity")
		}

		priceEntity := &price.Entity{Price: chargePrice}
		return nil, c.State().Insert(priceEntity)
	} else {
		return nil, errors.New("failed to initialize chaincode due to incorrect number of arguments")
	}
}

func QueryById(c router.Context) (interface{}, error) {
	var (
		contractor = c.ParamString("contractor")
		chargeId   = c.ParamString("chargeId")
	)

	return c.State().Get(&charge.Entity{
		Contractor: contractor,
		ChargeId:   chargeId,
	})
}

func QueryAll(c router.Context) (interface{}, error) {
	res, err := c.State().List(charge.TypeName, &charge.Entity{})
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(res).IsNil() {
		return []charge.Entity{}, nil
	}

	return res, nil
}

func InvokeStartChargeTransaction(c router.Context) (res interface{}, err error) {
	var startTransaction = c.Param("startTransaction").(charge.StartTransaction) // Assert the chaincode parameter

	// Validate current owner
	user, err := identity.FromStub(c.Stub())
	if user == nil || user.GetSubject() == "" {
		return nil, fmt.Errorf("missing identity email for this user")
	}

	userState, err := service.GetUser(c.Stub(), user.GetMSPID(), GetCertificateSubject(user.Cert))
	if err != nil || userState == nil {
		return nil, errors.Wrap(err, "could not get user state")
	}

	if userState.State == "blocked" {
		return nil, errors.Errorf("User %s:%s is blocked", user.GetMSPID(), GetCertificateSubject(user.Cert))
	}

	var chargeTransaction = &charge.Entity{
		Contractor: startTransaction.Contractor,
		ChargeId:   strconv.Itoa(rand.Int()),
		User:        SerializeIdentity(user.MspID, user.Cert),
		StartTime:  time.Now(),
		State:       charge.ChargeStateStarted, // Initial state
	}

	return chargeTransaction, c.State().Insert(chargeTransaction)
}

func InvokeStopChargeTransaction(c router.Context) (interface{}, error) {
	var (
		chargeTransaction charge.Entity
		// Buy transaction payload
		data = c.Param("stopTransaction").(charge.StopTransaction)

		// Get the current commercial paper state
		cp, err = c.State().Get(&charge.Entity{
			Contractor: data.Contractor,
			ChargeId:   data.ChargeId,
		}, &charge.Entity{})
	)

	if err != nil {
		return nil, errors.Wrap(err, "charge transaction not found")
	}

	chargeTransaction = cp.(charge.Entity)

	// Validate current owner
	user, err := identity.FromStub(c.Stub())
	if user != nil && chargeTransaction.User != SerializeIdentity(user.MspID, user.Cert) {
		return nil, fmt.Errorf(
			"transaction %s %s is not owned by this user", chargeTransaction.Contractor, chargeTransaction.ChargeId)
	}

	currentPrice, err := c.State().Get(&price.Entity{}, &price.Entity{})
	if err != nil {
		return nil, errors.Wrap(err, "could not calculate price")
	}

	// Check if transaction is in progress
	if chargeTransaction.State == charge.ChargeStateStarted {
		chargeTransaction.State = charge.ChargeStateStopped
		chargeTransaction.EndTime = time.Now()
		chargeTransaction.Price = currentPrice.(price.Entity).Price * int(time.Since(chargeTransaction.StartTime).Minutes())
	} else {
		return nil, fmt.Errorf("transaction can not be stopped from %s state", chargeTransaction.State)
	}

	return chargeTransaction, c.State().Put(chargeTransaction)
}

func InvokeCompleteChargeTransaction(c router.Context) (interface{}, error) {
	var (
		chargeTransaction charge.Entity
		// Buy transaction payload
		data = c.Param("completeTransaction").(charge.CompleteTransaction)

		// Get the current commercial paper state
		cp, err = c.State().Get(&charge.Entity{
			Contractor: data.Contractor,
			ChargeId:   data.ChargeId,
		}, &charge.Entity{})
	)

	if err != nil {
		return nil, errors.Wrap(err, "charge transaction not found")
	}

	chargeTransaction = cp.(charge.Entity)

	// TODO: Validate current invoker has this right

	// Check if transaction is in progress
	if chargeTransaction.State == charge.ChargeStateStopped {
		chargeTransaction.State = charge.ChargeStateCompleted
	} else {
		return nil, fmt.Errorf("transaction can not be completed from %s state", chargeTransaction.State)
	}

	return chargeTransaction, c.State().Put(chargeTransaction)
}
