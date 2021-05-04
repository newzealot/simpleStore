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

func AddUser(w http.ResponseWriter, r *http.Request) {
	clientId := os.Getenv("AWS_COGNITO_APP_CLIENT_ID")
	r.ParseForm()
	username := r.Form.Get("email")
	password := r.Form.Get("password")
	password2 := r.Form.Get("password2")
	if password != password2 {
		log.Println("Unable to register merchant - Passwords do not match")
		http.Redirect(w, r, "/register?error=Passwords do not match", http.StatusFound)
		return
	}
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Unable to start session - %s\n", err)
		http.Redirect(w, r, "/error", http.StatusFound)
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
	_, err = cognito.New(sess).SignUp(user)
	if err != nil {
		log.Printf("Unable to register merchant - %s\n", err)
		http.Redirect(w, r, "/error", http.StatusFound)
	} else {
		log.Printf("Successfully registered merchant - %s\n", username)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	return
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
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
		if err.Error() == "NotAuthorizedException: Incorrect username or password." {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err.Error() == "UserNotConfirmedException: User is not confirmed." {
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	log.Printf("Login successful - %s\n", username)
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

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	at, err := r.Cookie("AccessToken")
	if err != nil {
		log.Printf("Cookie issue - %s\n", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	password := r.Form.Get("password")
	newPassword := r.Form.Get("newpassword")
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Unable to start session - %s\n", err)
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	changeTry := &cognito.ChangePasswordInput{
		AccessToken:      aws.String(at.Value),
		PreviousPassword: aws.String(password),
		ProposedPassword: aws.String(newPassword),
	}
	_, err = cognito.New(sess).ChangePassword(changeTry)
	if err != nil {
		log.Printf("Unable to change password - %s\n", err)
		if err.Error() == "NotAuthorizedException: Incorrect username or password." {
			http.Redirect(w, r, "/change_password?error=Incorrect username or password.", http.StatusFound)
			return
		} else if err.Error() == "InvalidParameter: 1 validation error(s) found.\n- minimum field size of 6, ChangePasswordInput.ProposedPassword.\n" {
			http.Redirect(w, r, "/change_password?error=Min 6 characters", http.StatusFound)
			return
		} else {
			http.Redirect(w, r, "/error", http.StatusFound)
			return
		}
	}
	log.Println("Change password successful")
	http.Redirect(w, r, "/change_password?error=Change password successful", http.StatusSeeOther)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	clientId := os.Getenv("AWS_COGNITO_APP_CLIENT_ID")
	r.ParseForm()
	email := r.Form.Get("email")
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Unable to start session - %s\n", err)
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	resetTry := &cognito.ForgotPasswordInput{
		ClientId: aws.String(clientId),
		Username: aws.String(email),
	}
	_, err = cognito.New(sess).ForgotPassword(resetTry)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	log.Println("Verification code sent")
	http.Redirect(w, r, "/verification_code?email="+email, http.StatusFound)
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	clientId := os.Getenv("AWS_COGNITO_APP_CLIENT_ID")
	r.ParseForm()
	username := r.URL.Query().Get("email")
	code := r.Form.Get("code")
	password := r.Form.Get("password")
	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Unable to start session - %s\n", err)
		http.Redirect(w, r, "/error", http.StatusFound)
		return
	}
	codeTry := &cognito.ConfirmForgotPasswordInput{
		ClientId:         aws.String(clientId),
		ConfirmationCode: aws.String(code),
		Password:         aws.String(password),
		Username:         aws.String(username),
	}
	_, err = cognito.New(sess).ConfirmForgotPassword(codeTry)
	if err != nil {
		log.Println(err)
		if err.Error() == "CodeMismatchException: Invalid verification code provided, please try again." {
			http.Redirect(w, r, "/verification_code?error=Invalid code provided, please request a code again", http.StatusFound)
			return
		} else if err.Error() == "ExpiredCodeException: Invalid code provided, please request a code again." {
			http.Redirect(w, r, "/verification_code?error=Invalid code provided, please request a code again", http.StatusFound)
			return
		} else if err.Error() == "InvalidParameter: 1 validation error(s) found.\n- minimum field size of 6, ConfirmForgotPasswordInput.Password.\n" {
			http.Redirect(w, r, "/verification_code?error=Min 6 characters", http.StatusFound)
			return
		} else {
			http.Redirect(w, r, "/error", http.StatusFound)
			return
		}
	}
	log.Println("Password reset successful")
	http.Redirect(w, r, "/verification_code?error=Password reset successful", http.StatusFound)
}

func RefreshUser(w http.ResponseWriter, r *http.Request) {
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
	log.Println("Refresh successful")
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
