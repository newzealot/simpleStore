package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Collections []string

func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	c := Collections{}
	q := "SELECT DISTINCT COLLECTION_ID FROM PRODUCT_COLLECTION"
	rows, err := DB.Query(q)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		c = append(c, strings.Title(name))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}
