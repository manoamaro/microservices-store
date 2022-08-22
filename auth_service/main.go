package main

import (
	"encoding/json"
	"errors"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"manoamaro.github.com/auth_service/internal"
)

var authService *internal.AuthService

func main() {

	authService = internal.NewAuthService()

	r := mux.NewRouter()
	r.StrictSlash(true)

	s := r.PathPrefix("/auth").Subrouter()
	s.Path("/sign_in").Methods("POST").HandlerFunc(signInHandler)
	s.Path("/sign_up").Methods("POST").HandlerFunc(signUpHandler)
	s.Path("/validate").Methods("GET").HandlerFunc(validateHandler)

	srv := &http.Server{
		Addr:         "0.0.0.0:8081",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

var signUpHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		handleError(err, w, r)
	} else {
		request := &struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		if err := json.Unmarshal(body, request); err != nil {
			handleError(err, w, r)
		} else if auth, err := authService.CreateAuth(request.Email, request.Password); err != nil {
			handleError(err, w, r)
		} else if signedString, err := authService.CreateToken(auth); err != nil {
			handleError(err, w, r)
		} else {
			w.Header().Add("Authorization", "bearer "+signedString)
			w.WriteHeader(http.StatusOK)
		}
	}
})

var signInHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		handleError(err, w, r)
	} else {
		request := &struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		if err := json.Unmarshal(body, request); err != nil {
			handleError(err, w, r)
		} else if auth, found := authService.Authenticate(request.Email, request.Password); !found {
			handleError(errors.New("auth not found"), w, r)
		} else if signedString, err := authService.CreateToken(auth); err != nil {
			handleError(err, w, r)
		} else {
			w.Header().Add("Authorization", "bearer "+signedString)
			w.WriteHeader(http.StatusOK)
		}
	}
})

var validateHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if userClaims, err := authService.GetTokenFromRequest(r); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		response := struct {
			Audiences []string `json:"audiences"`
			Flags     []string `json:"flags"`
		}{
			Audiences: userClaims.Audience,
			Flags:     userClaims.Flags,
		}
		if responseJson, err := json.Marshal(response); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(responseJson)
		}
	}
})

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
}
