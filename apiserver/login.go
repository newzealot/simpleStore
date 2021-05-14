package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"log"
	"net/http"
	"os"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("email")
	password := r.Form.Get("password")
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Unable to start session - %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	authTry := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		ClientId: aws.String(os.Getenv("AWS_COGNITO_APP_CLIENT_ID")),
	}
	res, err := cognito.New(sess).InitiateAuth(authTry)
	if err != nil {
		log.Printf("Unable to login - %s\n", err)
		switch err.(awserr.Error).Code() {
		case "NotAuthorizedException":
			log.Println("NotAuthorizedException")
			w.WriteHeader(http.StatusUnauthorized)
			return
		case "UserNotConfirmedException":
			w.WriteHeader(http.StatusForbidden)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	j, err := json.Marshal(*res.AuthenticationResult)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
	return
}
