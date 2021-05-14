package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	D.GetData()
	r := mux.NewRouter()
	// personal preference for not using r.Prefix, as "/admin" will become "" and difficult to understand
	admin := r.Name("admin").Subrouter()
	admin.HandleFunc("/admin", AdminGET).Methods("GET")
	admin.HandleFunc("/admin/add_product", AddProductGET).Methods("GET")
	admin.HandleFunc("/admin/add_product", AddProductPOST).Methods("POST")
	admin.HandleFunc("/admin/edit_product/{id}", EditProductGET).Methods("GET")
	admin.HandleFunc("/admin/edit_product/{id}", EditProductPOST).Methods("POST")
	admin.HandleFunc("/admin/delete_product/{id}", DeleteProductGET).Methods("GET")
	admin.HandleFunc("/admin/delete_product/{id}", DeleteProductPOST).Methods("POST")
	admin.HandleFunc("/admin/delete_media/{id}/{filename}", DeleteMediaGET).Methods("GET")
	admin.Use(AllowOnlyMerchant)
	// keeping use of r.Host consistent
	customer := r.Name("customer").Subrouter()
	customer.HandleFunc("/customer/checkout", CheckoutPOST).Methods("POST")
	customer.Use(AllowOnlyCustomer)
	r.HandleFunc("/checkout-success", CheckoutSuccessGET).Methods("GET")
	r.HandleFunc("/checkout-cancel", CheckoutCancelGET).Methods("GET")
	r.HandleFunc("/add_cart/{id}", AddCartPOST).Methods("POST")
	r.HandleFunc("/about", AboutGET).Methods("GET")
	r.HandleFunc("/cart", CartGET).Methods("GET")
	r.HandleFunc("/login", LoginGET).Methods("GET")
	r.HandleFunc("/login", LoginPOST).Methods("POST")
	r.HandleFunc("/logout", LogoutGET).Methods("GET")
	r.HandleFunc("/register", RegisterGET).Methods("GET")
	r.HandleFunc("/register", RegisterPOST).Methods("POST")
	r.HandleFunc("/product/{id}", ProductGET).Methods("GET")
	r.HandleFunc("/collection/{id}", CollectionGET).Methods("GET")
	r.HandleFunc("/", IndexGET).Methods("GET")
	r.Use(GetUserInfo)

	if err := http.ListenAndServe(":5000", csrf.Protect([]byte("32-byte-long-auth-key"))(r)); err != nil {
		log.Fatalln(err)
	}
}
