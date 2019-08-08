package main

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/router"
)

func QueryAll(c router.Context) (interface{}, error) {
	return c.State().List(&RestrictedResponse{})
}

func QueryById(c router.Context) (interface{}, error) {
	id := c.Param("email")
	if id == nil || id == "" {
		return nil, errors.New("Empty email")
	}

	result, err := c.State().Get(id, &Entity{})

	if result == nil {
		user := &Entity{
			Email:   id.(string),
			State:   UserStateRegistered,
			Balance: 0,
		}
		if err = c.State().Insert(user); err != nil {
			return nil, errors.Wrap(err, "Could not create user")
		}

		return &RestrictedResponse{State: user.State}, err
	} else {
		user := result.(*Entity)
		return &RestrictedResponse{State: user.State}, err
	}
}
