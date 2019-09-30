package shared

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, code int, data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(code)
	_, err = w.Write(res)
	if err != nil {
		panic(err)
	}
}

type ListResponse struct {
	Data interface{} `json:"data"`
}

func WriteResponseList(w http.ResponseWriter, code int, data interface{}) {
	wrapped := ListResponse{
		Data: data,
	}
	res, err := json.Marshal(wrapped)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(code)
	_, err = w.Write(res)
	if err != nil {
		panic(err)
	}
}
