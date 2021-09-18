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
	mongoDBClient := internal.ConnectMongoDB(os.Getenv("MONGO_URL"))
	defer mongoDBClient.Disconnect(nil)
	internal.StartMQ(os.Getenv("AMQP_URL"))

	r := mux.NewRouter()
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

func ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	products, err := internal.DB.ListProducts()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := json.Marshal(&products)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Println(err)
	}
}

func PostProductsHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newProduct := models.Product{}

	err = json.Unmarshal(bytes, &newProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	savedProduct, err := internal.DB.InsertProduct(newProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	bytes, err = json.Marshal(&savedProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := internal.DB.FetchProduct(objId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := json.Marshal(product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write(bytes)
}

func UpdateProductsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}

	product := models.Product{}
	err = json.Unmarshal(bytes, &product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}

	updated, err := internal.DB.UpdateProduct(objId, product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}

	if updated {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func DeleteProductsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deleted, err := internal.DB.DeleteProduct(objId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if deleted {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
