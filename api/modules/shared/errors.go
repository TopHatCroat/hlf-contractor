package shared

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Status  int    `json:"code,omitempty"`
}

func WriteErrorResponse(w http.ResponseWriter, code int, err error) {
	print(err)
	res, err := json.Marshal(&ErrorResponse{Message: err.Error()})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(code)
	_, err = w.Write(res)
	if err != nil {
		panic(err)
	}
}
