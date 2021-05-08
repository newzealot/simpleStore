package api

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"log"
	"net/http"
	"os"
	. "simpleStore/apiserver/db"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	clientId := os.Getenv("AWS_COGNITO_APP_CLIENT_ID")
	r.ParseForm()
	username := r.Form.Get("email")
	password := r.Form.Get("password")
	password2 := r.Form.Get("password2")
	address := r.Form.Get("address")
	phone := r.Form.Get("phone")
	name := r.Form.Get("name")
	if password != password2 {
		log.Println("Unable to register merchant - Passwords do not match")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Unable to start session - %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := &cognito.SignUpInput{
		Username: aws.String(username),
		Password: aws.String(password),
		ClientId: aws.String(clientId),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("custom:type"),
				Value: aws.String("merchant"),
			},
		},
	}
	result, err := cognito.New(sess).SignUp(user)
	if err != nil {
		log.Printf("Unable to register merchant - %s\n", err)
		switch err.(awserr.Error).Code() {
		case "InvalidPasswordException":
			w.WriteHeader(http.StatusBadRequest)
			return
		case "UsernameExistsException":
			w.WriteHeader(http.StatusConflict)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// insert into merchant table
	q := "INSERT INTO MERCHANT(MERCHANT_ID,EMAIL,NAME,ADDRESS,PHONE) VALUES (?,?,?,?,?)"
	res, err := DB.Exec(q, result.UserSub, username, name, address, phone)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected, err := res.RowsAffected(); err != nil || rowsAffected == 0 {
		log.Println(fmt.Errorf("MERCHANT table 0 rows affected or %s\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully registered %s in Cognito\n", username)
	return
}
