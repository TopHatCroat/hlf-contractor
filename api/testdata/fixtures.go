package testdata

import (
	"fmt"
	"github.com/TopHatCroat/hlf-contractor/api/fabric"
	"github.com/TopHatCroat/hlf-contractor/api/modules"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"time"
)

const (
	userMsp = "AwesomeAgencyMSP"

	admin          = "admin@mail.com"
	adminPassword = "asdf"

	firstUser     = "user@mail.com"
	userPassword = "asdf"

	chargeProvider = "Pharmatic"
)

func expectUserIdentity(app *modules.App, username, password string, role shared.Role) shared.Identity {
	// ignore error, probably means that user is already registered
	err := app.Client.Register(username, password, role)
	if err != nil {
		println(err)
	}

	err = app.Client.Login(username, password)
	if err != nil {
		panic(err)
	}

	userIdentityResponse, err := app.Client.CA.GetIdentity(username, "")
	if err != nil {
		panic(err)
	}

	return shared.Identity{
		Msp:  userIdentityResponse.Affiliation,
		Id:   userIdentityResponse.ID,
		Role: shared.User,
	}
}

func InitFixtures(app *modules.App) {
	adminUserIdentity := expectUserIdentity(app, admin, userPassword, shared.Admin)
	_ = expectUserIdentity(app, firstUser, userPassword, shared.User)

	existingCharges, err := app.Client.AllCharges(&adminUserIdentity, chargeProvider)
	if err != nil {
		panic(err)
	}

	if (len(existingCharges) != 0) {
		fmt.Printf("Data already exists. Skipping fixtures...")
		return
	}

	adminEntity, err := app.Client.CreateUser(&adminUserIdentity, userMsp, admin)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created user: %+v\n", adminEntity)

	regularEntity, err := app.Client.CreateUser(&adminUserIdentity, userMsp, firstUser)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created user: %+v\n", regularEntity)

	// Wait for identity to propagate
	time.Sleep(2 * time.Second)

	chargeTransactions := make([]*fabric.ChargeTransaction, 5)
	for i := 0; i < 1; i++ {
		chargeTransaction, err := app.Client.StartCharge(&adminUserIdentity, chargeProvider)
		if err != nil {
			panic(err)
		}

		chargeTransactions[i] = chargeTransaction
	}

	createdChargeTransaction, err := app.Client.AllCharges(&adminUserIdentity, chargeProvider)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created transaction count: %d", len(createdChargeTransaction))

	for i := 0; i < 1; i++ {
		chargeTransaction, err := app.Client.StopCharge(&adminUserIdentity, chargeProvider, chargeTransactions[i].ChargeId)
		if err != nil {
			panic(err)
		}
		chargeTransactions[i] = chargeTransaction
	}

	adminIdentity := expectUserIdentity(app, firstUser, userPassword, shared.Admin)

	for i := 0; i < 1; i++ {
		chargeTransaction, err := app.Client.CompleteCharge(&adminIdentity, chargeProvider, chargeTransactions[i].ChargeId)
		if err != nil {
			panic(err)
		}

		chargeTransactions[i] = chargeTransaction
	}
}
