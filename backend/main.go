package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/product", postProduct).Methods("POST")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}
