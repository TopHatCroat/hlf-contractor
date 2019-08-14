package main

import (
	"github.com/TopHatCroat/hlf-contractor/api/modules"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	DEFAULT_CA_ULR = "localhost:7054"
)

func main() {
	app, err := modules.NewApp("./config.yml")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/register", app.Register).Methods("POST")
	r.HandleFunc("/login", app.Login).Methods("POST")

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
