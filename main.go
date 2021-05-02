package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	. "simpleStore/api"
	. "simpleStore/middleware"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/product", AddProduct).Methods("POST")
	r.HandleFunc("/api/v1/register", AddUser).Methods("POST")
	r.HandleFunc("/api/v1/login", LoginUser).Methods("POST")
	r.HandleFunc("/api/v1/refresh", RefreshUser).Methods("GET")
	r.HandleFunc("/api/v1/changepassword", ChangePassword).Methods("POST")
	r.HandleFunc("/api/v1/forgotpassword", ForgotPassword).Methods("POST")
	r.HandleFunc("/api/v1/verificationcode", VerifyEmail).Methods("POST")
	r.PathPrefix("/add_product").Handler(ValidationMiddleWare(http.FileServer(http.Dir("./static/")), "/add_product.html"))
	r.PathPrefix("/forgot_password").Handler(http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/change_password").Handler(ValidationMiddleWare(http.FileServer(http.Dir("./static/")), "/change_password.html"))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}
