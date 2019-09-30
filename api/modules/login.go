package modules

import (
	"encoding/json"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"io/ioutil"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

func (app *App) Login(w http.ResponseWriter, req *http.Request) {
	data := &LoginRequest{}
	raw, err := ioutil.ReadAll(req.Body)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	err = json.Unmarshal(raw, data)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	err = app.Client.Login(data.Email, data.Password)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	sessionToken := app.SetSession(data.Email)
	res := &LoginResponse{
		Token: sessionToken,
	}

	shared.WriteResponse(w, http.StatusCreated, res)
}
