package modules

import (
	"encoding/json"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

type LoginRequest struct {
	Email    string `email:"-"`
	Password string `password:"-"`
}

type LoginResponse struct {
	Token string `token:"-"`
}

func (app *App) Login(w http.ResponseWriter, req *http.Request) {
	data := &LoginRequest{}
	raw, err := ioutil.ReadAll(req.Body)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
	}
	err = json.Unmarshal(raw, data)
	if err != nil {
		shared.WriteErrorResponse(w, 400, err)
	}

	err = app.client.Register(data.Email, data.Password)

	sessionToken := uuid.New().String()
	app.sessions[sessionToken] = data.Email

	res := &LoginResponse{
		Token: sessionToken,
	}

	shared.WriteResponse(w, http.StatusCreated, res)
}
