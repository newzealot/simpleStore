package api

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"log"
	"net/http"
	"os"
	"strings"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	rt := r.Header.Get("Authorization")
	rt = strings.Replace(rt, "Bearer ", "", 1)
	sess, err := session.NewSession()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	authTry := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("REFRESH_TOKEN_AUTH"),
		AuthParameters: map[string]*string{
			"REFRESH_TOKEN": aws.String(rt),
		},
		ClientId: aws.String(os.Getenv("AWS_COGNITO_APP_CLIENT_ID")),
	}
	res, err := cognito.New(sess).InitiateAuth(authTry)
	if err != nil {
		log.Printf("RefreshToken error - %s\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	j, err := json.Marshal(*res.AuthenticationResult)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
