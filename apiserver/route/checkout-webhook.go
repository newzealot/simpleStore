package api

import (
	"encoding/json"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"github.com/stripe/stripe-go/webhook"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func FulfillOrder(returnedSession stripe.CheckoutSession) error {

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")
	s, err := session.Get(returnedSession.ID, params)
	if err != nil {
		return err
	}
	for _, v := range s.LineItems.Data {
		log.Println(v.Description)
		log.Println(v.Quantity)
		log.Println(v.Price.UnitAmount)
		log.Println(v.Price.Metadata["ProductID"])
		log.Println(v.Price.Metadata)
	}
	return nil
}

func CheckoutWebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	endpointSecret := "whsec_KGCel0AgKNkE9CBsKHaEL7V2VPkvIQ3P"
	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endpointSecret)

	if err != nil {
		log.Printf("Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	// Handle the checkout.session.completed event
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Fulfill the purchase...
		FulfillOrder(session)
	}
	w.WriteHeader(http.StatusOK)
}
