package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type ProductMedia struct {
	ProductID string
	Filename  string
}

func QueryDB(productid string, userID string, q string, values map[string][]string) error {
	singleMap := map[string]string{}
	for k, v := range values {
		if len(v) != 1 {
			return fmt.Errorf("empty value in form")
		}
		singleMap[k] = values[k][0][1 : len(values[k][0])-1]
	}
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
		fmt.Sprint(userID),
	}
	// call stored procedure
	_, err = DB.Exec(
		q, p.ProductID, p.Title, p.Description,
		p.SellingPrice, p.DiscountedPrice, p.CostPrice, p.ProductCollection,
		p.QuantityAvailable, p.MerchantID,
	)
	if err != nil {
		return err
	}
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
		contentType := ""
		switch filepath.Ext(filename) {
		case ".gif":
			contentType = "image/gif"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".webp":
			contentType = "image/webp"
		}
		storageKey := fmt.Sprintf("%s/%s", productid, filename)
		file, err := files[i].Open()
		if err != nil {
			return err
		}
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
			Key:         aws.String(storageKey),
			Body:        file,
			ContentType: aws.String(contentType),
		})
		if err != nil {
			file.Close()
			return err
		}
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
	}
	return nil
}

func ProductHandlerPOST(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userID")
	productID := uuid.NewString()
	if err := r.ParseMultipartForm(128 << 20); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot parse multipart form"))
		return
	}
	formdata := r.MultipartForm
	q := "CALL ADD_PRODUCT(?,?,?,?,?,?,?,?,?)"
	if err := QueryDB(productID, fmt.Sprint(userID), q, formdata.Value); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if len(formdata.File["media"]) != 0 {
		if err := AddToStorage(productID, formdata.File["media"]); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInsufficientStorage)
			w.Write([]byte(err.Error()))
			return
		}
	}
}

func ProductHandlerPUT(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := r.Header.Get("userID")
	productID := vars["id"]
	if err := r.ParseMultipartForm(128 << 20); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot parse multipart form"))
		return
	}
	formdata := r.MultipartForm
	q := "CALL EDIT_PRODUCT(?,?,?,?,?,?,?,?,?)"
	if err := QueryDB(productID, fmt.Sprint(userID), q, formdata.Value); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if len(formdata.File["media"]) != 0 {
		if err := AddToStorage(productID, formdata.File["media"]); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInsufficientStorage)
			w.Write([]byte(err.Error()))
			return
		}
	}
}

func ProductHandlerDELETE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]
	if err := DeleteFromStorage(productID); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	q := "DELETE FROM PRODUCT WHERE PRODUCT_ID = (?)"
	_, err := DB.Exec(q, productID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func DeleteFromStorage(productID string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Prefix: aws.String(productID),
	})
	if err != nil {
		return err
	}
	key := []string{}
	for _, v := range resp.Contents {
		key = append(key, *v.Key)
	}
	for _, v := range key {
		_, err = svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
			Key:    aws.String(v),
		})
		if err != nil {
			return err
		}
		err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
			Key:    aws.String(v),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
