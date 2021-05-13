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

type Product struct {
	ProductID         string
	Title             string
	Description       string
	FileName          []string
	FileNamePrefix    string
	SellingPrice      float64
	DiscountedPrice   float64
	CostPrice         float64
	QuantityAvailable int
	CollectionID      string
	MerchantID        string
	Created           time.Time
}

func (D2 *Data) FormattedDate() string {
	loc, _ := time.LoadLocation("Asia/Singapore")
	return D2.Created.In(loc).Format("02-Jan-06 15:04")
}

var D DataStore

func (D *DataStore) GetProduct(s string) Product {
	// Only difference is FileName is []string
	p := Product{}
	for _, v := range *D {
		if v.ProductID == s {
			p.ProductID = v.ProductID
			p.Title = v.Title
			p.Description = v.Description
			p.FileName = append(p.FileName, v.FileName)
			p.FileNamePrefix = fmt.Sprintf("%s%s/", os.Getenv("AWS_S3_URL_PREFIX"), v.ProductID)
			p.SellingPrice = v.SellingPrice
			p.DiscountedPrice = v.DiscountedPrice
			p.CostPrice = v.CostPrice
			p.QuantityAvailable = v.QuantityAvailable
			p.CollectionID = v.CollectionID
			p.MerchantID = v.MerchantID
			p.Created = v.Created
		}
	}
	t := []string{}
	for i := len(p.FileName) - 1; i >= 0; i-- {
		t = append(t, p.FileName[i])
	}
	p.FileName = t
	return p
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
	// Sort by Created as map above disrupts order
	sort.Slice(admin, func(i, j int) bool {
		return admin[i].Created.After(admin[j].Created)
	})
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
	// Sort by Created as map above disrupts order
	sort.Slice(collection, func(i, j int) bool {
		return collection[i].Created.After(collection[j].Created)
	})
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
	sort.Slice(collections, func(i, j int) bool {
		return collections[i] < collections[j]
	})
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
		time.Sleep(3 * time.Second)
		D.GetData()
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
