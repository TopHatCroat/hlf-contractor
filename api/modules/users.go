package modules

import (
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *App) GetUsers(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	users, err := app.Client.AllUsers(identity)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, users)
}

func (app *App) GetUserById(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	pathVars := mux.Vars(req)
	user, err := app.Client.FindUserById(identity, pathVars["id"])
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, user)
}

func (app *App) GetMe(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	user, err := app.Client.FindUserById(identity, identity.Id)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, user)
}
