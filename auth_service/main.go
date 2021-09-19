package main

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/gorilla/mux"
	"log"
	"manoamaro.github.com/auth_service/internal"
	"net/http"
	"os"
	"time"
)

func main() {
	internal.ConnectMongoDB(os.Getenv("MONGO_URL"))
	defer func() {
		if err := internal.DisconnectMongoDB(); err != nil {
			log.Println(err)
		}
	}()


	r := mux.NewRouter()
	r.StrictSlash(true)

	s := r.PathPrefix("/users").Subrouter()
	s.Path("/login").Methods("POST").HandlerFunc(loginHandler)
	s.Path("/signup").Methods("POST").HandlerFunc(signupHandler)
	s.Path("/").Methods("GET").Handler(authenticateMiddleware(getProfileHandler))
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

	err := srv.ListenAndServe()
	internal.FailOnError(err)
}

type UserInfo struct {
	Email string
	Access []string
}

type UserClaims struct {
	*jwt.StandardClaims
	UserInfo
}

var loginHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		&jwt.StandardClaims{
			Id: "1",
			ExpiresAt:  time.Now().Add(time.Hour * 168).Unix(),
		},
		UserInfo{
			Email:  "example@email.com",
			Access: []string{""},
		},
	})
	if signedString, err := token.SignedString([]byte("My Secret")); err != nil {
		handleError(err, w, r)
	} else {
		w.Header().Add("Authorization", "bearer " + signedString)
		w.WriteHeader(http.StatusOK)
	}
})

var signupHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

})

var getProfileHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

})

var updateProfileHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

})

var deleteProfileHandler = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

})


func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte("My Secret"), nil
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

func authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		if user != nil {
			userValues := user.(*jwt.Token).Claims.(jwt.MapClaims)
			if userValues != nil {
				access := userValues["access"].([]interface{})
				if access != nil {
					for _, v := range access {
						if v.(string) == "users" {
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
