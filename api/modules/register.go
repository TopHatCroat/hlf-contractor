package modules

import (
	"encoding/json"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"io/ioutil"
	"net/http"
)

type RegisterRequest struct {
	Email    string `email:"-"`
	Password string `password:"-"`
}

type RegisterResponse struct {
	Email string `email:"-"`
}

func (app *App) Register(w http.ResponseWriter, req *http.Request) {
	data := &RegisterRequest{}
	raw, err := ioutil.ReadAll(req.Body)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
	}
	err = json.Unmarshal(raw, data)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
	}

	err = app.client.Register(data.Email, data.Password)

	res := &RegisterResponse{
		Email: data.Email,
	}

	shared.WriteResponse(w, http.StatusCreated, res)
}
