package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	. "simpleStore/apiserver/db"
	. "simpleStore/apiserver/middleware"
	. "simpleStore/apiserver/route"
)

func main() {
	dbCloser := SetupDB()
	defer dbCloser()

	r := mux.NewRouter()
	protected := r.Host("localhost:8080").Subrouter()
	protected.HandleFunc("/api/v1/product", ProductHandlerPOST).Methods("POST")
	protected.HandleFunc("/api/v1/product/{id}", ProductHandlerPUT).Methods("PUT")
	protected.HandleFunc("/api/v1/product/{id}", ProductHandlerDELETE).Methods("DELETE")
	protected.HandleFunc("/api/v1/media/{id}/{filename}", MediaHandlerDELETE).Methods("DELETE")
	protected.HandleFunc("/api/v1/checkout", CheckoutHandler).Methods("POST")
	protected.HandleFunc("/api/v1/checkout/{id}", CheckoutCancelHandler).Methods("PUT")
	protected.Use(Validation)
	r.HandleFunc("/api/v1/checkout-webhook", CheckoutWebhookHandler).Methods("POST")
	r.HandleFunc("/api/v1/data", DataHandler).Methods("GET")
	r.HandleFunc("/api/v1/collection", CollectionHandler).Methods("GET")
	r.HandleFunc("/api/v1/refresh", RefreshHandler).Methods("POST")
	r.HandleFunc("/api/v1/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", LoginHandler).Methods("POST")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}
