package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type DataStore []Data

type Data struct {
	ProductID         string
	Title             string
	Description       string
	FileName          string
	SellingPrice      float64
	DiscountedPrice   float64
	CostPrice         float64
	QuantityAvailable int
	CollectionID      string
	MerchantID        string
	Created           time.Time
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
	ds := DataStore{}
	q := "CALL GET_PRODUCT()"
	rows, err := DB.Query(q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d Data
		err := rows.Scan(
			&d.ProductID,
			&d.Title,
			&d.Description,
			&d.FileName,
			&d.SellingPrice,
			&d.DiscountedPrice,
			&d.CostPrice,
			&d.QuantityAvailable,
			&d.CollectionID,
			&d.MerchantID,
			&d.Created,
		)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ds = append(ds, d)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ds)
}
