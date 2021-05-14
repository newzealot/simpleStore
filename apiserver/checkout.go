package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Order struct {
	ProductID       string
	Title           string
	OrderQuantity   int
	DiscountedPrice float64
	SubTotal        float64
}

type ConfirmedOrder struct {
	ProductID         string
	Title             string
	OrderQuantity     int
	QuantityAvailable int
	DiscountedPrice   float64
	FileName          string
}

type createCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
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
	// verify orderList SellingPrice same as in DB and OrderQuantity less than in DB
	confirmedList := []ConfirmedOrder{}
	confirmed := ConfirmedOrder{}
	for _, v := range orderList {
		q := "CALL GET_PRODUCT_FOR_PAYMENT(?)"
		row := DB.QueryRow(q, v.ProductID)
		if err := row.Scan(&confirmed.ProductID, &confirmed.Title, &confirmed.FileName, &confirmed.DiscountedPrice, &confirmed.QuantityAvailable); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		confirmed.FileName = fmt.Sprintf("%s%s/%s", os.Getenv("AWS_S3_URL_PREFIX"), confirmed.ProductID, url.QueryEscape(confirmed.FileName))
		if confirmed.DiscountedPrice != v.DiscountedPrice {
			log.Println("Price not equal")
			w.WriteHeader(http.StatusConflict)
			return
		}
		if confirmed.QuantityAvailable < v.OrderQuantity {
			log.Println("Quantity not enough")
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		confirmed.OrderQuantity = v.OrderQuantity
		confirmedList = append(confirmedList, confirmed)
	}
	// Write order to ORDERITEM
	orderID := uuid.NewString()
	userID := r.Header.Get("userID")
	totalPrice := 0.0
	for _, v := range confirmedList {
		totalPrice += v.DiscountedPrice * float64(v.OrderQuantity)
	}
	q := "INSERT INTO ORDERITEM(ORDERITEM_ID,CUSTOMER_ID,TOTALPRICE,STATUS) VALUES (?,?,?,?);"
	res, err := DB.Exec(q, orderID, userID, totalPrice, "Awaiting Payment")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		log.Printf("%s or 0 rows affected\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// write order to PRODUCT_ORDERITEM
	for _, v := range confirmedList {
		q := "INSERT INTO PRODUCT_ORDERITEM(PRODUCT_ID,ORDERITEM_ID,QUANTITY,DISCOUNTEDPRICE) VALUES (?,?,?,?);"
		res, err := DB.Exec(q, v.ProductID, orderID, v.OrderQuantity, v.DiscountedPrice)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
			log.Printf("%s or 0 rows affected\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// write order to PRODUCT
		q = "UPDATE PRODUCT_QUANTITY SET QUANTITYAVAILABLE = QUANTITYAVAILABLE - (?) WHERE PRODUCT_ID = (?) AND QUANTITYAVAILABLE>=(?);"
		res, err = DB.Exec(q, v.OrderQuantity, v.ProductID, v.OrderQuantity)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
			log.Printf("%s or 0 rows affected\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// stripe integration
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	LineItems := []*stripe.CheckoutSessionLineItemParams{}
	for _, v := range confirmedList {
		item := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(string(stripe.CurrencySGD)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:   stripe.String(v.Title),
					Images: stripe.StringSlice([]string{v.FileName}),
				},
				UnitAmount: stripe.Int64(int64(v.DiscountedPrice * 100)),
			},
			Quantity: stripe.Int64(int64(v.OrderQuantity)),
		}
		LineItems = append(LineItems, item)
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: LineItems,
		Mode:      stripe.String(string(stripe.CheckoutSessionModePayment)),
		// hardcoded localhost is needed to test as it would not work otherwise
		SuccessURL: stripe.String("http://localhost:5000/checkout-success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("http://localhost:5000/checkout-cancel?session_id={CHECKOUT_SESSION_ID}"),
	}
	params.AddMetadata("SimpleStoreOrderID", orderID)
	session, err := session.New(params)
	if err != nil {
		log.Printf("session.New: %v", err)
	}
	data := createCheckoutSessionResponse{
		SessionID: session.ID,
	}
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

func CheckoutCancelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	// find out affected product_id and quantity
	q := "SELECT PRODUCT_ID, QUANTITY FROM PRODUCT_ORDERITEM WHERE ORDERITEM_ID = (?)"
	rows, err := DB.Query(q, orderID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var productID string
		var quantityOrdered int
		if err := rows.Scan(&productID, &quantityOrdered); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// reverting
		q = "UPDATE PRODUCT_QUANTITY SET QUANTITYAVAILABLE = QUANTITYAVAILABLE + (?) WHERE PRODUCT_ID = (?)"
		res, err := DB.Exec(q, quantityOrdered, productID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
			log.Printf("product_quantity %s or 0 rows affected\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	q = "UPDATE ORDERITEM SET STATUS = (?) WHERE ORDERITEM_ID = (?);"
	res, err := DB.Exec(q, "Cancelled", orderID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		log.Printf("orderitem%s or 0 rows affected\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
