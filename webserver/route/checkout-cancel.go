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
	// to delete line below
	at = "eyJraWQiOiJlaUhVK0lISTgyRFlMczNIOWZEVDRcL1BMSWtidHBoQklmSVh1M09LOFwveWM9IiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiI4MWE4NDE4Mi0wMGQ1LTQ2MTEtOTQwNC1hODhiNDQ1MTYyNzgiLCJldmVudF9pZCI6ImI4YzFmZTIzLTQ0MjMtNDRjNi05NjRiLWY5OWM5ZjVlYjQ5MSIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE2MjA3NTAxMzAsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tXC9hcC1zb3V0aGVhc3QtMV9tQWVpSFZ6SlciLCJleHAiOjE2MjA3ODE1OTMsImlhdCI6MTYyMDc3OTc5MywianRpIjoiMjk2NjAxNjUtOGY5NC00MjI0LWI3OWQtMGI5NTM2ODRiYWMzIiwiY2xpZW50X2lkIjoiN3R1bXA0bG5iMms0OGQ4Zmo5N2tiZWE3YTMiLCJ1c2VybmFtZSI6IjgxYTg0MTgyLTAwZDUtNDYxMS05NDA0LWE4OGI0NDUxNjI3OCJ9.iZJsxxHCOh7pqU45Vpv91c7P0OpJt0bcBsVy3Dy18AxTzijuLv5-Y9PuqxzsjN9GMfBUkgK67OqSICJ6dcSJf3PuCLD7KfJ-COHGoeEzkwpNH3aXzrWiaN4vUqnq7CrtqVuKQY8x2gULIxxYtu7wocUQs5waxwNiKw2niQC24zkg_HXgt9J92I5jBfYTsUPSNNaFksbV7sZmO_DxLKKXQrEsUg7D4GBm33sDDJjruU6Aiq8vR3CmsxtmcS7IIBLjMMxHV93g1wFD22svJuql2CBuXai2RlAEze-NrjbGI5b1v36B30LryLWrnZzeTbRx6xcjiJpyweXJ1ND0G9m7bg"
	// end delete
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
