package main

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"io/ioutil"
	"log"
	"manoamaro.github.com/mongodb"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manoamaro.github.com/products_service/internal"
	"manoamaro.github.com/products_service/internal/models"
)

var db *internal.DB

func main() {
	//internal.StartMQ(os.Getenv("AMQP_URL"))
	db = &internal.DB{
		MongoDB: mongodb.ConnectMongoDB(os.Getenv("MONGO_URL"), internal.DATABASE),
	}
	defer func() {
		if err := db.DisconnectMongoDB(); err != nil {
			log.Println(err)
		}
	}()

	r := mux.NewRouter()
	r.StrictSlash(true)

	s := r.PathPrefix("/products").Subrouter()
	s.Path("/").Methods("GET").HandlerFunc(ListProductsHandler)
	s.Path("/").Methods("POST").HandlerFunc(PostProductsHandler)
	s.Path("/{id}").Methods("GET").HandlerFunc(GetProductsHandler)
	s.Path("/{id}").Methods("PUT").HandlerFunc(UpdateProductsHandler)
	s.Path("/{id}").Methods("DELETE").HandlerFunc(DeleteProductsHandler)

	s.Use(jwtMiddleware, authenticateMiddleware, contentTypeMiddleware)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	err := srv.ListenAndServe()
	internal.FailOnError(err)
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
}

func ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	if products, err := db.ListProducts(); err != nil {
		handleError(err, w, r)
	} else if bytes, err := json.Marshal(&products); err != nil {
		handleError(err, w, r)
	} else if _, err = w.Write(bytes); err != nil {
		handleError(err, w, r)
	}
}

func PostProductsHandler(w http.ResponseWriter, r *http.Request) {
	if bytes, err := ioutil.ReadAll(r.Body); err != nil {
		handleError(err, w, r)
	} else {
		newProduct := models.Product{}
		if err = json.Unmarshal(bytes, &newProduct); err != nil {
			handleError(err, w, r)
		} else if savedProduct, err := db.InsertProduct(newProduct); err != nil {
			handleError(err, w, r)
		} else if bytes, err = json.Marshal(&savedProduct); err != nil {
			handleError(err, w, r)
		} else if _, err := w.Write(bytes); err != nil {
			handleError(err, w, r)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if objId, err := primitive.ObjectIDFromHex(id); err != nil {
		handleError(err, w, r)
	} else if product, err := db.FetchProduct(objId); err != nil {
		handleError(err, w, r)
	} else if bytes, err := json.Marshal(product); err != nil {
		handleError(err, w, r)
	} else if _, err := w.Write(bytes); err != nil {
		handleError(err, w, r)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateProductsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if objId, err := primitive.ObjectIDFromHex(id); err != nil {
		handleError(err, w, r)
	} else if bytes, err := ioutil.ReadAll(r.Body); err != nil {
		handleError(err, w, r)
	} else {
		product := models.Product{}
		if err = json.Unmarshal(bytes, &product); err != nil {
			handleError(err, w, r)
		} else if updated, err := db.UpdateProduct(objId, product); err != nil {
			handleError(err, w, r)
		} else if !updated {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func DeleteProductsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if objId, err := primitive.ObjectIDFromHex(id); err != nil {
		handleError(err, w, r)
	} else if deleted, err := db.DeleteProduct(objId); err != nil {
		handleError(err, w, r)
	} else if !deleted {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
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
						if v.(string) == "products" {
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

func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte("My Secret"), nil
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _ := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, keyFunc, request.WithClaims(jwt.MapClaims{}))
		if token != nil {
			next.ServeHTTP(w, r.Clone(context.WithValue(r.Context(), "user", token)))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
