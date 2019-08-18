package modules

import (
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"github.com/gorilla/mux"
	"net/http"
)

func (app *App) GetCharges(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	users, err := app.Client.AllCharges(identity, "")
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, users)
}

func (app *App) GetChargeByProvider(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	pathVars := mux.Vars(req)
	user, err := app.Client.FindChargeById(identity, pathVars["provider"], "")
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, user)
}

func (app *App) GetChargeById(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	pathVars := mux.Vars(req)
	chargeTransaction, err := app.Client.FindChargeById(identity, pathVars["provider"], pathVars["id"])
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, chargeTransaction)
}

func (app *App) StartCharge(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	if err := req.ParseForm(); err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	provider := req.FormValue("provider")
	users, err := app.Client.StartCharge(identity, provider)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, users)
}

func (app *App) StopCharge(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	if err := req.ParseForm(); err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	provider := req.FormValue("provider")
	chargeId := req.FormValue("chargeId")
	users, err := app.Client.StopCharge(identity, provider, chargeId)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, users)
}

func (app *App) CompleteCharge(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	if err := req.ParseForm(); err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	provider := req.FormValue("provider")
	chargeId := req.FormValue("chargeId")
	users, err := app.Client.CompleteCharge(identity, provider, chargeId)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, users)
}
