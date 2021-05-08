package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
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

var D DataStore

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

func (D *DataStore) GetAdmin(s string) []Data {
	// Allow only 1 unique row (results in only one filename)
	admin := []Data{}
	temp := map[string]Data{}
	for _, v := range *D {
		v.FileName = fmt.Sprintf("%s%s/%s", os.Getenv("AWS_S3_URL_PREFIX"), v.ProductID, v.FileName)
		temp[v.ProductID] = v
	}
	for _, v := range temp {
		if v.MerchantID == s {
			admin = append(admin, v)
		}
	}
	return admin
}

func (D *DataStore) GetCollection(s string) []Data {
	// Allow only 1 unique row (results in only one filename)
	collection := []Data{}
	temp := map[string]Data{}
	for _, v := range *D {
		v.FileName = fmt.Sprintf("%s%s/%s", os.Getenv("AWS_S3_URL_PREFIX"), v.ProductID, v.FileName)
		temp[v.ProductID] = v
	}
	for _, v := range temp {
		if s == "All" || v.CollectionID == s {
			collection = append(collection, v)
		}
	}
	return collection
}

func (D *DataStore) GetMenu() []string {
	collections := []string{}
	temp := map[string]bool{}
	for _, v := range *D {
		temp[v.CollectionID] = true
	}
	for k, _ := range temp {
		t := strings.Title(k)
		collections = append(collections, t)
	}
	collections = QuickSort(collections, 0, len(collections)-1)
	collections = append([]string{"All"}, collections...)
	return collections
}

func (D *DataStore) GetData() {
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
	if err := json.Unmarshal(b, D); err != nil {
		log.Println(err)
		return
	}
	// sort by newest item first
	sort.Slice(*D, func(i, j int) bool {
		return (*D)[i].Created.After((*D)[j].Created)
	})
}
