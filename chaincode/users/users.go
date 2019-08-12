package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	"reflect"
)

func getOrCreateUser(ctx router.Context, mspId, username string) (*Entity, error) {
	result, err := ctx.State().Get(&Entity{
		MSP:   mspId,
		Email: username,
	}, &Entity{})

	// If user not found, create it
	if result == nil {
		user := &Entity{
			MSP:     mspId,
			Email:   username,
			State:   UserStateActive,
			Balance: 0,
		}
		if err = ctx.State().Insert(user); err != nil {
			return nil, errors.Wrap(err, "Could not create user")
		}

		return user, nil
	} else {
		user := result.(Entity)
		return &user, nil
	}
}

func QueryAll(c router.Context) (interface{}, error) {
	certIdentity, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	if !IdentityIsAdmin(certIdentity) {
		return nil, nil
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

	user, err := getOrCreateUser(c, mspId.(string), username.(string))
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

	user, err := getOrCreateUser(c, mspId.(string), username.(string))
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

	user, err := getOrCreateUser(c, mspId.(string), username.(string))
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
