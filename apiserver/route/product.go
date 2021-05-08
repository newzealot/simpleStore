package api

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	. "simpleStore/apiserver/db"
	"strconv"
)

type ProductMedia struct {
	ProductID string
	Filename  string
}

func AddToDB(productid string, merchantid string, values map[string][]string) error {
	singleMap := map[string]string{}
	for k, v := range values {
		if len(v) != 1 {
			return fmt.Errorf("empty value in form")
		}
		singleMap[k] = values[k][0][1 : len(values[k][0])-1]
		log.Println(singleMap[k])
	}
	log.Println(merchantid)
	sellingPrice, err := strconv.ParseFloat(singleMap["sellingprice"], 64)
	if err != nil {
		return err
	}
	discountedPrice, err := strconv.ParseFloat(singleMap["discountedprice"], 64)
	if err != nil {
		return err
	}
	costPrice, err := strconv.ParseFloat(singleMap["costprice"], 64)
	if err != nil {
		return err
	}
	quantityAvailable, err := strconv.Atoi(singleMap["quantityavailable"])
	if err != nil {
		return err
	}
	p := struct {
		ProductID         string
		Title             string
		Description       string
		SellingPrice      float64
		DiscountedPrice   float64
		CostPrice         float64
		ProductCollection string
		QuantityAvailable int
		MerchantID        string
	}{
		productid,
		singleMap["title"],
		singleMap["description"],
		sellingPrice,
		discountedPrice,
		costPrice,
		singleMap["collectionid"],
		quantityAvailable,
		fmt.Sprint(merchantid),
	}
	// call stored procedure
	q := "CALL ADD_PRODUCT(?,?,?,?,?,?,?,?,?)"
	res, err := DB.Exec(
		q, p.ProductID, p.Title, p.Description,
		p.SellingPrice, p.DiscountedPrice, p.CostPrice, p.ProductCollection,
		p.QuantityAvailable, p.MerchantID,
	)
	if err != nil {
		return err
	}
	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		return fmt.Errorf("PRODUCT table 0 rows affected or %s\n", err)
	}
	log.Printf("%s inserted into DB successfully\n", productid)
	return nil
}

func AddToStorage(productid string, files []*multipart.FileHeader) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	var pmArray []ProductMedia
	uploader := s3manager.NewUploader(sess)
	for i := range files {
		filename := files[i].Filename
		storageKey := fmt.Sprintf("%s/%s", productid, filename)
		file, err := files[i].Open()
		if err != nil {
			return err
		}
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
			Key:    aws.String(storageKey),
			Body:   file,
		})
		if err != nil {
			file.Close()
			return err
		}
		log.Printf("Upload Successful - %s", storageKey)
		file.Close()
		pm := ProductMedia{
			ProductID: productid,
			Filename:  filename,
		}
		pmArray = append(pmArray, pm)
	}
	// Add to db
	for _, v := range pmArray {
		q := "INSERT INTO PRODUCT_MEDIA(PRODUCT_ID,FILENAME) VALUES (?,?)"
		res, err := DB.Exec(q, v.ProductID, v.Filename)
		if err != nil {
			return err
		}
		if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
			log.Printf("%s or 0 rows affected\n", err)
			return fmt.Errorf("%s or 0 rows affected\n", err)
		}
		log.Printf("%s/%s inserted into DB successfully\n", v.ProductID, v.Filename)
	}
	return nil
}

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	// Cookie already verified by middleware, so skipping error checks
	merchantID := r.Header.Get("merchantID")
	// normal err checks resume from here on
	productid := uuid.NewString()
	if err := r.ParseMultipartForm(128 << 20); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot parse multipart form"))
		return
	}
	formdata := r.MultipartForm
	if err := AddToDB(productid, fmt.Sprint(merchantID), formdata.Value); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if err := AddToStorage(productid, formdata.File["media"]); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInsufficientStorage)
		w.Write([]byte(err.Error()))
		return
	}
}