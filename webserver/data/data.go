package data

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

type Card struct {
	PictureURL      string
	SellingPrice    float64
	DiscountedPrice float64
}

var d DataStore
var Collections []string

func partition(arr []string, leftIndex int, rightIndex int) int {
	pivotValue := arr[(leftIndex+rightIndex)/2]
	for leftIndex < rightIndex {
		for arr[leftIndex] < pivotValue {
			leftIndex += 1
		}
		for arr[rightIndex] > pivotValue {
			rightIndex -= 1
		}
		if leftIndex < rightIndex {
			arr[leftIndex], arr[rightIndex] = arr[rightIndex], arr[leftIndex]
		}
	}
	return leftIndex
}

func QuickSort(arr []string, leftIndex int, rightIndex int) []string {
	if leftIndex < rightIndex {
		pivotIndex := partition(arr, leftIndex, rightIndex)
		QuickSort(arr, leftIndex, pivotIndex-1)
		QuickSort(arr, pivotIndex+1, rightIndex)
	}
	return arr
}

func GetCollection() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("APISERVER")+"/api/v1/collection", nil)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if err := json.Unmarshal(b, &Collections); err != nil {
		log.Println(err)
		return
	}
	Collections = QuickSort(Collections, 0, len(Collections)-1)
	Collections = append([]string{"All"}, Collections...)
	log.Println("Collections Retrieved", Collections)
}

func GetData() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("APISERVER")+"/api/v1/data", nil)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if err := json.Unmarshal(b, &d); err != nil {
		log.Println(err)
		return
	}
	log.Println("DataStore Retrieved", d)
}
