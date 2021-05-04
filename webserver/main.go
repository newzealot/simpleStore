package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	. "simpleStore/webserver/route"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/add_product", AddProductGET).Methods("GET")
	r.HandleFunc("/add_product", AddProductPOST).Methods("POST")
	r.HandleFunc("/login", LoginGET).Methods("GET")
	r.HandleFunc("/login", LoginPOST).Methods("POST")

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/register.html")
	})

	r.HandleFunc("/forgot_password", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/forgot_password.html")
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	r.HandleFunc("/verification_code", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/verification_code.html")
	})

	r.HandleFunc("/change_password", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/change_password.html")
	})
	r.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/error.html")
	})
	if err := http.ListenAndServe(":5000", csrf.Protect([]byte("32-byte-long-auth-key"))(r)); err != nil {
		log.Fatalln(err)
	}
}
