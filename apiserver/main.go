package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	. "simpleStore/apiserver/middleware"
	. "simpleStore/apiserver/route"
)

func main() {
	r := mux.NewRouter()
	protected := r.Host("localhost:8080").Subrouter()
	protected.HandleFunc("/api/v1/product", AddProduct).Methods("POST")
	protected.Use(Validation)

	r.HandleFunc("/api/v1/refresh", Refresh).Methods("POST")
	r.HandleFunc("/api/v1/register", Register).Methods("POST")
	r.HandleFunc("/api/v1/login", Login).Methods("POST")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}
