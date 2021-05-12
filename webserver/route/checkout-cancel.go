package route

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"log"
	"net/http"
	"os"
	. "simpleStore/webserver/data"
	"strings"
)

func CheckoutCancelGET(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	s, _ := session.Get(r.URL.Query()["session_id"][0], nil)
	orderID := s.Metadata["SimpleStoreOrderID"]
	at := r.Header.Get("Authorization")
	at = strings.Replace(at, "Bearer ", "", 1)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", os.Getenv("APISERVER")+"/api/v1/checkout/"+orderID, nil)
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
		log.Printf("Unable to revert order %s\n", orderID)
		http.Redirect(w, r, "/cart?error=Unable to revert order", http.StatusSeeOther)
		return
	}
	D.GetData()
	http.Redirect(w, r, "/cart?error=Payment not completed", http.StatusSeeOther)
}
