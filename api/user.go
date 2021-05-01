package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"log"
	"net/http"
	"os"
	"time"
)

// TODO change password

func ComputeSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(os.Getenv("AWS_COGNITO_APP_CLIENT_SECRET")))
	mac.Write([]byte(username + os.Getenv("AWS_COGNITO_APP_CLIENT_ID")))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("email")
	password := r.Form.Get("password")
	sess, err := session.NewSession()
	if err != nil {
		log.Println(err)
	}
	user := &cognito.SignUpInput{
		Username:   aws.String(username),
		Password:   aws.String(password),
		ClientId:   aws.String(os.Getenv("AWS_COGNITO_APP_CLIENT_ID")),
		SecretHash: aws.String(ComputeSecretHash(username)),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("custom:type"),
				Value: aws.String("merchant"),
			},
		},
	}
	_, err = cognito.New(sess).SignUp(user)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Successfully Registered")
	http.Redirect(w, r, "http://localhost:63342/simpleStore/frontend/verify.html", http.StatusFound)
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	otp := r.Form.Get("verificationcode")
	username := r.Form.Get("email")
	sess, err := session.NewSession()
	if err != nil {
		log.Println(err)
	}
	user := &cognito.ConfirmSignUpInput{
		ConfirmationCode: aws.String(otp),
		Username:         aws.String(username),
		ClientId:         aws.String(os.Getenv("AWS_COGNITO_APP_CLIENT_ID")),
		SecretHash:       aws.String(ComputeSecretHash(username)),
	}
	_, err = cognito.New(sess).ConfirmSignUp(user)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Successfully Verified")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	const flowUsernamePassword = "USER_PASSWORD_AUTH"
	const flowRefreshToken = "REFRESH_TOKEN_AUTH"
	r.ParseForm()
	username := r.Form.Get("email")
	password := r.Form.Get("password")
	refresh := r.Form.Get("refresh")
	refreshToken := r.Form.Get("refresh_token")
	sess, err := session.NewSession()
	if err != nil {
		log.Println(err)
	}
	flow := aws.String(flowUsernamePassword)
	params := map[string]*string{
		"USERNAME": aws.String(username),
		"PASSWORD": aws.String(password),
	}
	if os.Getenv("AWS_COGNITO_APP_CLIENT_SECRET") != "" {
		secretHash := ComputeSecretHash(username)
		params["SECRET_HASH"] = aws.String(secretHash)
	}
	if refresh != "" {
		flow = aws.String(flowRefreshToken)
		params = map[string]*string{
			"REFRESH_TOKEN": aws.String(refreshToken),
		}
	}
	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
		ClientId:       aws.String(os.Getenv("AWS_COGNITO_APP_CLIENT_ID")),
	}
	res, err := cognito.New(sess).InitiateAuth(authTry)
	if err != nil {
		log.Println(err)
		return
	}
	c1 := http.Cookie{
		Name:     "AccessToken",
		Value:    *res.AuthenticationResult.AccessToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Second * time.Duration(*res.AuthenticationResult.ExpiresIn)),
	}
	c2 := http.Cookie{
		Name:     "IdToken",
		Value:    *res.AuthenticationResult.IdToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Second * time.Duration(*res.AuthenticationResult.ExpiresIn)),
	}
	c3 := http.Cookie{
		Name:     "RefreshToken",
		Value:    *res.AuthenticationResult.RefreshToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 30),
	}
	http.SetCookie(w, &c1)
	http.SetCookie(w, &c2)
	http.SetCookie(w, &c3)
}
