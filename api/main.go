package main

import (
	"encoding/json"
	"github.com/TopHatCroat/hlf-contractor/api/fabric"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	DEFAULT_CA_ULR = "localhost:7054"
)

var (
	client *fabric.Client
)

func main() {
	fabClient, err := fabric.New("./config.yml")
	if err != nil {
		panic(err)
	}

	client = fabClient

	r := mux.NewRouter()
	r.HandleFunc("/register", Register).Methods("POST").Headers()

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err = srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Message string
}

type RegisterRequest struct {
	Email    string
	Name     string
	Password string
}

func Register(w http.ResponseWriter, r *http.Request) {
	data := &RegisterRequest{}
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse(400, err, w)
	}
	err = json.Unmarshal(raw, data)
	if err != nil {
		errorResponse(400, err, w)
	}

	err = client.RegisterUser(data.Email, data.Password)
	if err != nil {
		errorResponse(400, err, w)
	}

	w.WriteHeader(201)
}

func errorResponse(code int, err error, w http.ResponseWriter) {
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
