package route

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func OrderPOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	log.Println(vars["id"])
	log.Println(r.PostFormValue("orderquantity"))
	if r.PostFormValue("type") == "cart" {
		log.Println("cart")
	} else {
		log.Println("buy")
	}
}
