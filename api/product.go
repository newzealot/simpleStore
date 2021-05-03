package api

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/schema"
	"github.com/lestrrat-go/jwx/jwt"
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
	Seller            string
}

func (p *Product) AddToDB() {

}

func AddToStorage(username string, productTitle string, files []*multipart.FileHeader) ([]string, error) {
	var filenames []string
	sess, err := session.NewSession()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	uploader := s3manager.NewUploader(sess)
	for i := range files {
		filename := files[i].Filename
		storageKey := fmt.Sprintf("%s/%s/%s", username, productTitle, filename)
		file, err := files[i].Open()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
			Key:    aws.String(storageKey),
			Body:   file,
		})
		if err != nil {
			log.Println(err)
			file.Close()
			return nil, err
		}
		log.Printf("Upload Successful - %s", storageKey)
		file.Close()
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	// Cookie already verified by middleware, so skipping error checks
	at, _ := r.Cookie("AccessToken")
	pat, _ := jwt.Parse([]byte(at.Value))
	merchant, _ := pat.Get("username")
	// normal err checks resume from here on
	if err := r.ParseMultipartForm(128 << 20); err != nil {
		log.Println(err)
		return
	}
	formdata := r.MultipartForm
	var decoder = schema.NewDecoder()
	var p Product
	if err := decoder.Decode(&p, formdata.Value); err != nil {
		log.Println(err)
	}
	p.Seller = fmt.Sprint(merchant)
	files := formdata.File["media[]"] // grab the filenames
	filenames, err := AddToStorage(p.Seller, p.Title, files)
	if err != nil {
		log.Println(err)
	}
	p.Media = filenames
	p.AddToDB()
}
