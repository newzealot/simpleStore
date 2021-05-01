package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"os"
	. "simpleStore/api"
)

func ValidationMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c1, err := r.Cookie("AccessToken")
		if err != nil {
			log.Println(err)
			//TODO add redirect
			return
		}
		url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
		keyset, err := jwk.Fetch(context.Background(), url)
		src := []byte(c1.Value)
		iss := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
		_, err = jwt.Parse(src,
			jwt.WithKeySet(keyset),
			jwt.WithValidate(true),
			jwt.WithClaimValue("client_id", os.Getenv("AWS_COGNITO_APP_CLIENT_ID")), // replacing aud
			jwt.WithIssuer(iss),
			jwt.WithClaimValue("token_use", "access"),
		)
		if err != nil {
			log.Println(err)
			//TODO if expired, call refresh function
			//TODO add redirect
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/product", AddProduct).Methods("POST")
	r.HandleFunc("/api/v1/register", AddUser).Methods("POST")
	r.HandleFunc("/api/v1/verify", VerifyUser).Methods("POST")
	r.HandleFunc("/api/v1/login", LoginUser).Methods("POST")
	r.PathPrefix("/add_product").Handler(ValidationMiddleWare(http.FileServer(http.Dir("./static/"))))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln(err)
	}
}

//// create a serve mux
//sm := mux.NewRouter()
//
//// register handlers
//postR := sm.Methods(http.MethodPost).Subrouter()
//postR.HandleFunc("/signup", uh.Signup)
//postR.HandleFunc("/login", uh.Login)
//postR.Use(uh.MiddlewareValidateUser)
//
//refToken := sm.PathPrefix("/refresh-token").Subrouter()
//refToken.HandleFunc("", uh.RefreshToken)
//refToken.Use(uh.MiddlewareValidateRefreshToken)
//
//getR := sm.Methods(http.MethodGet).Subrouter()
//getR.HandleFunc("/greet", uh.Greet)
//getR.Use(uh.MiddlewareValidateAccessToken)
