package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func MediaHandlerDELETE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := fmt.Sprintf("%s/%s", vars["id"], vars["filename"])
	sess, err := session.NewSession()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	deleter := s3.New(sess)
	// Delete from s3
	_, err = deleter.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")), Key: aws.String(filename)})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = deleter.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    aws.String(filename),
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Delete from db
	q := "CALL DELETE_MEDIA(?,?)"
	res, err := DB.Exec(q, vars["id"], vars["filename"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		log.Printf("%s or 0 rows affected\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
