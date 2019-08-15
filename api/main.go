package main

import (
	"context"
	"github.com/TopHatCroat/hlf-contractor/api/modules"
	"github.com/TopHatCroat/hlf-contractor/api/modules/shared"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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

	router := mux.NewRouter()

	router.HandleFunc("/register", app.Register).Methods("POST")
	router.HandleFunc("/login", app.Login).Methods("POST")

	router.HandleFunc("/users", Authenticated(app, app.GetUsers)).Methods("GET")
	router.HandleFunc("/me", Authenticated(app, app.GetMe)).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err = srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func Authenticated(app *modules.App, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionToken := r.Header.Get("Authorization")
		if sessionToken == "" {
			shared.WriteErrorResponse(w, 403, errors.New("not authorized"))
			return
		}

		userIdentity := app.GetSession(sessionToken)
		if userIdentity == "" {
			shared.WriteErrorResponse(w, 403, errors.New("not authorized"))
			return
		}

		identity, err := app.Client.CA.GetIdentity(userIdentity, "")
		if err != nil {
			shared.WriteErrorResponse(w, 403, errors.Wrap(err, "not authorized"))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "identity", identity)))
	})
}
