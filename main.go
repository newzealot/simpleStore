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
	protected := r.Host("localhost:8080").Subrouter()
	protected.HandleFunc("/add_product", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/add_product.html")
	})
	protected.HandleFunc("/api/v1/product", AddProduct).Methods("POST")
	protected.Use(ValidationMiddleWare)

	r.HandleFunc("/api/v1/register", AddUser).Methods("POST")
	r.HandleFunc("/api/v1/login", LoginUser).Methods("POST")
	r.HandleFunc("/api/v1/refresh", RefreshUser).Methods("GET")
	r.HandleFunc("/api/v1/changepassword", ChangePassword).Methods("POST")
	r.HandleFunc("/api/v1/forgotpassword", ForgotPassword).Methods("POST")
	r.HandleFunc("/api/v1/verificationcode", VerifyEmail).Methods("POST")
	//r.HandleFunc("/add_product", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, "static/add_product.html")
	//})
	r.HandleFunc("/change_password", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/change_password.html")
	})
	r.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/error.html")
	})
	r.HandleFunc("/forgot_password", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/forgot_password.html")
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/register.html")
	})
	r.HandleFunc("/verification_code", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/verification_code.html")
	})

	//r.PathPrefix("/add_product").Handler(ValidationMiddleWare(http.FileServer(http.Dir("./static/")), "/add_product.html"))
	//r.PathPrefix("/forgot_password").Handler(http.FileServer(http.Dir("./static/")))
	//r.PathPrefix("/change_password").Handler(ValidationMiddleWare(http.FileServer(http.Dir("./static/")), "/change_password.html"))
	//r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}
