package shared

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/msp/api"
	"github.com/pkg/errors"
	"net/http"
)

type Role string

const (
	Admin = "admin"
	User  = "user"
)

type Identity struct {
	Msp  string
	Id   string
	Role Role
}

func ExpectIdentity(req *http.Request) (*Identity, error) {
	identity := req.Context().Value("identity").(*api.IdentityResponse)
	if identity == nil {
		return nil, errors.New("you must be logged in")
	}

	return &Identity{
		Id:   identity.ID,
		Msp:  identity.Affiliation,
		Role: User,
	}, nil
}
