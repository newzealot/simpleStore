package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	. "simpleStore/apiserver/db"
)

type Order struct {
	ProductID     string
	Title         string
	OrderQuantity int
	SellingPrice  float64
	SubTotal      float64
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

	//if err != nil {
	//	log.Println(err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var d Data
	//	err := rows.Scan(
	//		&d.ProductID,
	//		&d.Title,
	//		&d.Description,
	//		&d.FileName,
	//		&d.SellingPrice,
	//		&d.DiscountedPrice,
	//		&d.CostPrice,
	//		&d.QuantityAvailable,
	//		&d.CollectionID,
	//		&d.MerchantID,
	//		&d.Created,
	//	)
	//	if err != nil {
	//		log.Println(err)
	//		w.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	ds = append(ds, d)
	//}
	return
}
