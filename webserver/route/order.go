package route

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Order struct {
	ProductID     string
	Title         string
	OrderQuantity int
	SellingPrice  float64
	SubTotal      float64
}

func OrderPOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	t := r.PostFormValue("title")
	oq, err := strconv.Atoi(r.PostFormValue("orderquantity"))
	sp, err := strconv.ParseFloat(r.PostFormValue("sellingprice"), 64)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/product/"+vars["id"]+"error?Order quantity not integer", http.StatusSeeOther)
		return
	}
	order := Order{
		vars["id"],
		t,
		oq,
		sp,
		float64(oq) * sp,
	}
	orderList := []Order{}
	cart, err := r.Cookie("Cart")
	if err != nil {
		log.Println("No cookie")
	} else {
		j, err := base64.StdEncoding.DecodeString(cart.Value)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/product/"+vars["id"]+"?error=Cannot decode base64 string", http.StatusSeeOther)
			return
		}
		if err := json.Unmarshal(j, &orderList); err != nil {
			log.Println(err)
			http.Redirect(w, r, "/product/"+vars["id"]+"?error=Cannot unmarshall json", http.StatusSeeOther)
			return
		}
	}
	log.Println(orderList)
	alreadyInList := false
	for i, v := range orderList {
		if v.ProductID == order.ProductID {
			log.Println(v)
			orderList[i].OrderQuantity = v.OrderQuantity + order.OrderQuantity
			alreadyInList = true
			log.Println(v)
		}
	}
	if alreadyInList == false {
		orderList = append(orderList, order)
	}
	log.Println(orderList)
	j, err := json.Marshal(orderList)
	jbase := base64.StdEncoding.EncodeToString(j)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/product/"+vars["id"]+"?error=Cannot marshall json", http.StatusSeeOther)
		return
	}
	c := http.Cookie{
		Name:     "Cart",
		Value:    jbase,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7),
	}
	http.SetCookie(w, &c)
	if r.PostFormValue("type") == "buy" {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	} else {
		http.Redirect(w, r, "/product/"+vars["id"]+"?success=Added product to cart", http.StatusSeeOther)
		return
	}
}
