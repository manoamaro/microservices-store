package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"manoamaro.github.com/products_service/internal"
)

func main() {
	mongoDBClient := internal.ConnectMongoDB(os.Getenv("MONGO_URL"))
	defer mongoDBClient.Disconnect(nil)
	internal.StartMQ(os.Getenv("AMQP_URL"))

	r := mux.NewRouter()
	s := r.PathPrefix("/products").Subrouter()

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

func ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	products, err := internal.DB.ListProducts()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(&products)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(bytes)
	w.WriteHeader(http.StatusOK)
}

func PostProductsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func UpdateProductsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func DeleteProductsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}
