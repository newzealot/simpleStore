package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func CheckoutPOST(w http.ResponseWriter, r *http.Request) {
	at := r.Header.Get("Authorization")
	at = strings.Replace(at, "Bearer ", "", 1)
	orderList := []Order{}
	cart, err := r.Cookie("Cart")
	if err != nil {
		log.Println("No cookie")
		http.Redirect(w, r, "/cart?error=Empty cart", http.StatusSeeOther)
		return
	} else {
		// these steps are just to verify cookie is valid - able to unmarshal into struct
		j, err := base64.StdEncoding.DecodeString(cart.Value)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/cart?error=Cannot decode base64 string", http.StatusSeeOther)
			return
		}
		if err := json.Unmarshal(j, &orderList); err != nil {
			log.Println(err)
			http.Redirect(w, r, "/cart?error=Cannot unmarshall json", http.StatusSeeOther)
			return
		}
	}
	j, err := json.Marshal(orderList)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/cart?error=Cannot marshall json", http.StatusSeeOther)
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("APISERVER")+"/api/v1/checkout", bytes.NewBuffer(j))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/cart?error=Something went wrong", http.StatusSeeOther)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+at)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/cart?error=Something went wrong", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
		http.Redirect(w, r, "/cart?error=Something went wrong", http.StatusSeeOther)
		return
	}
	js, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/cart?error=Something went wrong", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	D.GetData()
	return
}
