package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"manoamaro.github.com/auth_service/models"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/gorilla/mux"
	"manoamaro.github.com/auth_service/internal"
)

func getMongoDBUrl() string {
	if value, exists := os.LookupEnv("MONGO_URL"); exists {
		return value
	}
	return "mongodb://localhost:27017"
}

var db *internal.MongoDB

func main() {

	db = internal.ConnectMongoDB(getMongoDBUrl())

	defer func() {
		if err := db.DisconnectMongoDB(); err != nil {
			log.Println(err)
		}
	}()

	r := mux.NewRouter()
	r.StrictSlash(true)

	s := r.PathPrefix("/users").Subrouter()
	s.Path("/login").Methods("POST").HandlerFunc(loginHandler)
	s.Path("/signup").Methods("POST").HandlerFunc(signupHandler)
	s.Path("/logout").Methods("POST").HandlerFunc(signupHandler)
	s.Path("/").Methods("GET").Handler(getProfileHandler)
	s.Path("/").Methods("PUT").HandlerFunc(updateProfileHandler)
	s.Path("/").Methods("DELETE").HandlerFunc(deleteProfileHandler)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
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
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		if err := json.Unmarshal(body, request); err != nil {
			handleError(err, w, r)
		} else if user, err := db.LoginUser(request.Email, request.Password); err != nil {
			handleError(err, w, r)
		} else if signedString, err := internal.GetTokenSigned(user.Id.Hex(), user.Email); err != nil {
			handleError(err, w, r)
		} else if response, err := json.Marshal(user); err != nil {
			handleError(err, w, r)
		} else {
			w.Header().Add("Authorization", "bearer "+signedString)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(response)
		}
	}
})

var signupHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		handleError(err, w, r)
	} else {
		request := &struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		if err := json.Unmarshal(body, request); err != nil {
			handleError(err, w, r)
		} else if user, err := db.CreateUser(models.User{FullName: request.Name, Email: request.Email}, request.Password); err != nil {
			handleError(err, w, r)
		} else if signedString, err := internal.GetTokenSigned(user.Id.Hex(), user.Email); err != nil {
			handleError(err, w, r)
		} else if response, err := json.Marshal(user); err != nil {
			handleError(err, w, r)
		} else {
			w.Header().Add("Authorization", "bearer "+signedString)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(response)
		}
	}
})

var getProfileHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

})

var updateProfileHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

})

var deleteProfileHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

})

func keyFunc(_ *jwt.Token) (interface{}, error) {
	return internal.GetJWTSecret(), nil
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, keyFunc, request.WithClaims(jwt.MapClaims{})); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r.Clone(context.WithValue(r.Context(), "user", token)))
		}
	})
}

func authenticateMiddleware(next http.Handler, requiredRole string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := r.Context().Value("user"); user != nil {
			if userValues := user.(*jwt.Token).Claims.(jwt.MapClaims); userValues != nil {
				if access := userValues["access"].([]interface{}); access != nil {
					for _, v := range access {
						if v.(string) == requiredRole {
							next.ServeHTTP(w, r)
							return
						}
					}
				}
			}
		}
		w.WriteHeader(http.StatusForbidden)
	})
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
}
