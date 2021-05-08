package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	. "simpleStore/webserver/data"
	. "simpleStore/webserver/middleware"
	. "simpleStore/webserver/route"
)

func main() {
	D.GetData()
	r := mux.NewRouter()
	protected := r.Host("localhost:5000").Subrouter()
	protected.HandleFunc("/admin", AdminGET).Methods("GET")
	protected.HandleFunc("/admin/add_product", AddProductGET).Methods("GET")
	protected.HandleFunc("/admin/add_product", AddProductPOST).Methods("POST")
	protected.HandleFunc("/admin/edit_product/{id}", EditProductGET).Methods("GET")
	protected.HandleFunc("/admin/edit_product/{id}", EditProductPOST).Methods("POST")
	//protected.HandleFunc("/admin/delete_media/{id}", DeleteMediaGET).Methods("GET")
	//protected.HandleFunc("/admin/delete_product/{id}", DeleteProductGET).Methods("GET")
	//protected.HandleFunc("/admin/delete_product/{id}", DeleteProductPOST).Methods("POST")
	protected.Use(AccessTokenCheck)
	r.HandleFunc("/login", LoginGET).Methods("GET")
	r.HandleFunc("/login", LoginPOST).Methods("POST")
	r.HandleFunc("/logout", LogoutGET).Methods("GET")
	r.HandleFunc("/register", RegisterGET).Methods("GET")
	r.HandleFunc("/register", RegisterPOST).Methods("POST")
	r.HandleFunc("/collection/{id}", CollectionGET).Methods("GET")
	r.HandleFunc("/", IndexGET).Methods("GET")

	if err := http.ListenAndServe(":5000", csrf.Protect([]byte("32-byte-long-auth-key"))(r)); err != nil {
		log.Fatalln(err)
	}
}
