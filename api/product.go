package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func (p *Product) AddToDB() {

}

func AddToStorage(user string, files []*multipart.FileHeader) ([]string, error) {
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
			Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
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

func AddProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("AddProduct")
	//if err := r.ParseMultipartForm(128 << 20); err != nil {
	//	log.Println(err)
	//	return
	//}
	//formdata := r.MultipartForm
	//files := formdata.File["media[]"] // grab the filenames
	//filenames, err := AddToStorage("user", files)
	//if err != nil {
	//	log.Println(err)
	//}
	//var decoder = schema.NewDecoder()
	//var p Product
	//if err := decoder.Decode(&p, formdata.Value); err != nil {
	//	log.Println(err)
	//}
	//p.Media = filenames
	//fmt.Println(p)
}
