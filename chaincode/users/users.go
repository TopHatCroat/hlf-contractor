package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	"reflect"
)

func getUser(ctx router.Context, mspId, username string) (*Entity, error) {
	result, err := ctx.State().Get(&Entity{
		MSP:   mspId,
		Email: username,
	}, &Entity{})

	if err != nil {
		return nil, err
	}

	user := result.(Entity)
	return &user, nil
}

func QueryAll(c router.Context) (interface{}, error) {
	certIdentity, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	if !IdentityIsAdmin(certIdentity) {
		user, err := getUser(c, certIdentity.MspID, GetCertificateSubject(certIdentity.Cert))
		if err != nil {
			return nil, err
		}

		return []Entity{*user}, nil
	}

	users, err := c.State().List(TypeName, &Entity{})
	if err != nil {
		return nil, err
	}

	if reflect.ValueOf(users).IsNil() {
		return []Entity{}, nil
	}

	return users, nil
}

func QueryById(c router.Context) (interface{}, error) {
	mspId := c.Param("mspId")
	username := c.Param("email")
	if mspId == nil || mspId == "" {
		return nil, errors.New("MSP must be specified")
	}

	if username == nil || username == "" {
		return nil, errors.New("Email must be specified")
	}

	user, err := getUser(c, mspId.(string), username.(string))
	if err != nil {
		return nil, err
	}

	certIdentity, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}
	// Return full user data if user has admin role or is current user
	if IdentityIsAdmin(certIdentity) || IdentityIsEqual(certIdentity, mspId.(string), username.(string)) {
		return user, nil
	}

	return &RestrictedResponse{
		MSP:   user.MSP,
		Email: user.Email,
		State: user.State,
	}, nil
}

func InvokeCreateUser(c router.Context) (interface{}, error) {

	mspId := c.Param("mspId")
	username := c.Param("email")
	if mspId == nil || mspId == "" {
		return nil, errors.New("MSP must be specified")
	}

	if username == nil || username == "" {
		return nil, errors.New("Email must be specified")
	}

	user := &Entity{
		MSP:     mspId.(string),
		Email:   username.(string),
		State:   UserStateActive,
		Balance: 0,
	}

	if err := c.State().Insert(user); err != nil {
		return nil, errors.Wrap(err, "Could not create user")
	}

	return user, nil
}

func InvokeBlockUserTransaction(c router.Context) (interface{}, error) {
	certIdentity, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	if !IdentityIsAdmin(certIdentity) {
		return nil, nil
	}

	mspId := c.Param("mspId")
	username := c.Param("email")
	if mspId == nil || mspId == "" {
		return nil, errors.New("MSP must be specified")
	}

	if username == nil || username == "" {
		return nil, errors.New("Email must be specified")
	}

	user, err := getUser(c, mspId.(string), username.(string))
	if err != nil {
		return nil, err
	}

	if user.State == UserStateActive {
		user.State = UserStateBlocked
	} else {
		return nil, fmt.Errorf("user can not be blocked from %s state", user.State)
	}

	return user, c.State().Put(user)
}

func InvokeUnblockUserTransaction(c router.Context) (interface{}, error) {
	certIdentity, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	if !IdentityIsAdmin(certIdentity) {
		return nil, nil
	}

	mspId := c.Param("mspId")
	username := c.Param("email")
	if mspId == nil || mspId == "" {
		return nil, errors.New("MSP must be specified")
	}

	if username == nil || username == "" {
		return nil, errors.New("Email must be specified")
	}

	user, err := getUser(c, mspId.(string), username.(string))
	if err != nil {
		return nil, err
	}

	if user.State == UserStateBlocked {
		user.State = UserStateActive
	} else {
		return nil, fmt.Errorf("user can not be blocked from %s state", user.State)
	}

	return user, c.State().Put(user)
}
