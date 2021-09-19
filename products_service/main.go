package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"manoamaro.github.com/products_service/internal"
	"manoamaro.github.com/products_service/internal/models"
)

func main() {
	internal.StartMQ(os.Getenv("AMQP_URL"))
	internal.ConnectMongoDB(os.Getenv("MONGO_URL"))
	defer func() {
		if err := internal.DisconnectMongoDB(); err != nil {
			log.Println(err)
		}
	}()

	r := mux.NewRouter()
	r.StrictSlash(true)
	s := r.PathPrefix("/products").Subrouter()
	s.Use(contentTypeMiddleware)

	s.Path("/").Methods("GET").HandlerFunc(ListProductsHandler)
	s.Path("/").Methods("POST").HandlerFunc(PostProductsHandler)
	s.Path("/{id}").Methods("GET").HandlerFunc(GetProductsHandler)
	s.Path("/{id}").Methods("PUT").HandlerFunc(UpdateProductsHandler)
	s.Path("/{id}").Methods("DELETE").HandlerFunc(DeleteProductsHandler)

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

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	log.Println(err)
	w.WriteHeader(http.StatusBadRequest)
}

func ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	if products, err := internal.DB.ListProducts(); err != nil {
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
		} else if savedProduct, err := internal.DB.InsertProduct(newProduct); err != nil {
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
	} else if product, err := internal.DB.FetchProduct(objId); err != nil {
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
		} else if updated, err := internal.DB.UpdateProduct(objId, product); err != nil {
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
	} else if deleted, err := internal.DB.DeleteProduct(objId); err != nil {
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
