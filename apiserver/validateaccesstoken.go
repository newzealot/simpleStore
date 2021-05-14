package main

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"os"
	"strings"
)

func VerifyToken(at string) (string, error) {
	url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
	keyset, err := jwk.Fetch(context.Background(), url)
	if err != nil {
		return "", err
	}
	src := []byte(at)
	iss := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
	p, err := jwt.Parse(src,
		jwt.WithKeySet(keyset),
		jwt.WithValidate(true),
		jwt.WithClaimValue("client_id", os.Getenv("AWS_COGNITO_APP_CLIENT_ID")), // replacing aud
		jwt.WithIssuer(iss),
		jwt.WithClaimValue("token_use", "access"),
	)
	if err != nil {
		return "", err
	}
	m, ok := p.Get("username")
	if !ok {
		log.Println("No username field")
	}
	return fmt.Sprint(m), nil
}

func ValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		at := r.Header.Get("Authorization")
		at = strings.Replace(at, "Bearer ", "", 1)
		userID, err := VerifyToken(at)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// setting username in header
		r.Header.Add("userID", fmt.Sprint(userID))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
