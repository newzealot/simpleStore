package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/schema"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Product struct {
	Title             string
	Description       string
	Media             []string `schema:"media[]"`
	SellingPrice      float64
	OriginalPrice     float64
	CostPrice         float64
	QuantityAvailable int
	Collections       string
}

func (p *Product) addToDB() {

}

func addToStorage(user string, files []*multipart.FileHeader) ([]string, error) {
	var filenames []string
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	uploader := s3manager.NewUploader(sess)
	for i := range files {
		filename := files[i].Filename
		file, err := files[i].Open()
		if err != nil {
			return nil, err
		}
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    aws.String(user + "/" + filename),
			Body:   file,
		})
		if err != nil {
			file.Close()
			return nil, err
		}
		file.Close()
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

func postProduct(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(128 << 20); err != nil {
		log.Println(err)
		return
	}
	formdata := r.MultipartForm
	files := formdata.File["media[]"] // grab the filenames
	filenames, err := addToStorage("user", files)
	if err != nil {
		log.Println(err)
	}
	var decoder = schema.NewDecoder()
	var p Product
	if err := decoder.Decode(&p, formdata.Value); err != nil {
		log.Println(err)
	}
	p.Media = filenames
	fmt.Println(p)
}
