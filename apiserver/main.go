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
	r.HandleFunc("/api/v1/product", ProductHandler).Methods("POST")
	protected.HandleFunc("/api/v1/media", MediaHandler).Methods("POST")
	protected.Use(Validation)
	r.HandleFunc("/api/v1/media", MediaHandler).Methods("POST")
	r.HandleFunc("/api/v1/refresh", RefreshHandler).Methods("POST")
	r.HandleFunc("/api/v1/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", LoginHandler).Methods("POST")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}