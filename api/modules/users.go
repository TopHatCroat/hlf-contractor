package modules

import (
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp/api"
	"github.com/pkg/errors"
	"net/http"
)

func (app *App) GetUsers(w http.ResponseWriter, req *http.Request) {
	identity := req.Context().Value("identity").(*api.IdentityResponse)
	if identity == nil {
		shared.WriteErrorResponse(w, 403, errors.New("you must be logged in"))
		return
	}

	users, err := app.Client.AllUsers(identity)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, users)
}

func (app *App) GetMe(w http.ResponseWriter, req *http.Request) {
	identity := req.Context().Value("identity").(*api.IdentityResponse)
	if identity == nil {
		shared.WriteErrorResponse(w, 403, errors.New("you must be logged in"))
		return
	}

	user, err := app.Client.FindUserById(identity, identity.ID)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, user)
}
