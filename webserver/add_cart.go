package main

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
	ProductID       string
	Title           string
	OrderQuantity   int
	DiscountedPrice float64
	SubTotal        float64
}

func AddCartPOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	r.ParseForm()
	t := r.PostFormValue("title")
	oq, err := strconv.Atoi(r.PostFormValue("orderquantity"))
	dp, err := strconv.ParseFloat(r.PostFormValue("discountedprice"), 64)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/product/"+vars["id"]+"error?Order quantity not integer", http.StatusSeeOther)
		return
	}
	order := Order{
		vars["id"],
		t,
		oq,
		dp,
		float64(oq) * dp,
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
	alreadyInList := false
	for i, v := range orderList {
		if v.ProductID == order.ProductID {
			orderList[i].OrderQuantity = v.OrderQuantity + order.OrderQuantity
			alreadyInList = true
		}
	}
	if alreadyInList == false {
		orderList = append(orderList, order)
	}
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
