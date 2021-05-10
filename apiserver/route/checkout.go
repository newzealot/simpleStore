package api

import (
	"encoding/json"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	. "simpleStore/apiserver/db"
)

type Order struct {
	ProductID     string
	Title         string
	OrderQuantity int
	SellingPrice  float64
	SubTotal      float64
}

type createCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("in here")
	orderList := []Order{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(b, &orderList); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(orderList)
	// verify orderList SellingPrice same as in DB and OrderQuantity less than in DB
	var sellingPrice float64
	var quantityAvailable int
	for _, v := range orderList {
		q := "SELECT PRODUCT.SELLINGPRICE, PRODUCT_QUANTITY.QUANTITYAVAILABLE FROM PRODUCT INNER JOIN PRODUCT_QUANTITY ON PRODUCT.PRODUCT_ID = PRODUCT_QUANTITY.PRODUCT_ID WHERE PRODUCT.PRODUCT_ID = (?)"
		row := DB.QueryRow(q, v.ProductID)
		if err := row.Scan(&sellingPrice, &quantityAvailable); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if sellingPrice != v.SellingPrice {
			log.Println(err)
			w.WriteHeader(http.StatusConflict)
			return
		}
		if quantityAvailable < v.OrderQuantity {
			log.Println(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	}
	// stripe integration
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	LineItems := []*stripe.CheckoutSessionLineItemParams{}
	for _, v := range orderList {
		item := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(string(stripe.CurrencySGD)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(v.Title),
				},
				UnitAmount: stripe.Int64(int64(v.SellingPrice * 100)),
			},
			Quantity: stripe.Int64(int64(v.OrderQuantity)),
		}
		LineItems = append(LineItems, item)
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems:  LineItems,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://www.google.com/"),
		CancelURL:  stripe.String("https://www.bing.com/"),
	}
	session, err := session.New(params)
	if err != nil {
		log.Printf("session.New: %v", err)
	}
	data := createCheckoutSessionResponse{
		SessionID: session.ID,
	}
	log.Println("session: ", session)
	log.Println("data: ", data)
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}
