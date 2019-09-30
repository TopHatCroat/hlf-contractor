package modules

import (
	"encoding/json"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type UserFilter struct {
	Ids []string `json:"id"`
}

func (app *App) GetUsers(w http.ResponseWriter, req *http.Request) {
	identity, err := shared.ExpectIdentity(req)
	if err != nil {
		shared.WriteErrorResponse(w, 403, err)
		return
	}

	filterQry := req.URL.Query().Get("filter")
	var filter UserFilter
	if err := json.Unmarshal([]byte(filterQry), &filter); err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	if len(filter.Ids) != 0 {
		var res []interface{}
		for _, userId := range filter.Ids {
			idParts := strings.Split(userId, ":")
			usr, err := app.Client.FindUserById(identity, idParts[0], idParts[1])
			if err != nil {
				shared.WriteErrorResponse(w, 400, err)
				return
			}

			res = append(res, usr)
		}

		shared.WriteResponse(w, 200, res)
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
	idParts := strings.Split(pathVars["id"], ":")
	user, err := app.Client.FindUserById(identity, idParts[0], idParts[1])
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

	user, err := app.Client.FindUserById(identity, identity.Msp, identity.Id)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	shared.WriteResponse(w, 200, user)
}
