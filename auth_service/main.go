package main

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4/request"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"manoamaro.github.com/auth_service/models"
	"net/http"
	"strconv"
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
	s.Path("/sign_up").Methods("POST").HandlerFunc(loginHandler)
	s.Path("/validate").Methods("GET").HandlerFunc(validateHandler)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	if err := srv.ListenAndServe(); err != nil {
		internal.FailOnError(err)
	}
}

var loginHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		handleError(err, w, r)
	} else {
		request := &struct {
			UserId uint64 `json:"user_id"`
		}{}
		if err := json.Unmarshal(body, request); err != nil {
			handleError(err, w, r)
		} else if auth, found := authService.FindAuth(request.UserId); !found {
			handleError(errors.New("auth not found"), w, r)
		} else if signedString, err := internal.GetTokenSigned(
			strconv.FormatUint(auth.UserId, 10),
			mapTo(auth.Roles, func(i models.Role) string { return i.Name }),
			mapTo(auth.Flags, func(i models.Flag) string { return i.Name })); err != nil {
			handleError(err, w, r)
		} else {
			w.Header().Add("Authorization", "bearer "+signedString)
			w.WriteHeader(http.StatusOK)
		}
	}
})

func mapTo[I interface{}, O interface{}](i []I, f func(I) O) []O {
	var output []O
	for _, el := range i {
		output = append(output, f(el))
	}
	return output
}

var validateHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, internal.GetJWTSecretFunc, request.WithClaims(&internal.UserClaims{})); err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
	} else if userValues := token.Claims.(*internal.UserClaims); userValues == nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else if userId, err := strconv.ParseUint(userValues.ID, 10, 64); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	} else if _, found := authService.FindAuth(userId); !found {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		response := struct {
			Roles []string `json:"roles"`
			Flags []string `json:"flags"`
		}{
			Roles: userValues.Roles,
			Flags: userValues.Flags,
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
