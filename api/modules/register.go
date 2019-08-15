package modules

import (
	"encoding/json"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"io/ioutil"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type RegisterResponse struct {
	Email string `json:"email,omitempty"`
}

func (app *App) Register(w http.ResponseWriter, req *http.Request) {
	data := &RegisterRequest{}
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

	err = app.Client.Register(data.Email, data.Password)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
		return
	}

	res := &RegisterResponse{
		Email: data.Email,
	}

	shared.WriteResponse(w, http.StatusCreated, res)
}
