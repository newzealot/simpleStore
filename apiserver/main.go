package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	. "simpleStore/apiserver/api"
	. "simpleStore/apiserver/middleware"
)

func main() {
	r := mux.NewRouter()
	protected := r.Host("localhost:8080").Subrouter()
	protected.HandleFunc("/api/v1/changepassword", ChangePassword).Methods("POST")
	protected.HandleFunc("/api/v1/product", AddProduct).Methods("POST")
	protected.Use(ValidationMiddleWare)

	r.HandleFunc("/api/v1/refresh", RefreshUser).Methods("POST")
	r.HandleFunc("/api/v1/register", AddUser).Methods("POST")
	r.HandleFunc("/api/v1/login", LoginUser).Methods("POST")
	r.HandleFunc("/api/v1/forgotpassword", ForgotPassword).Methods("POST")
	r.HandleFunc("/api/v1/verificationcode", VerifyEmail).Methods("POST")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}
