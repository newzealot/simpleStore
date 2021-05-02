package middleware

import (
	"context"
	"fmt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"os"
)

//func ValidationMiddleWare(url string) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			c1, err := r.Cookie("AccessToken")
//			if err != nil {
//				log.Println(err)
//				if err.Error() == "http: named cookie not present" {
//					http.Redirect(w, r, "/api/v1/refresh?url="+url, http.StatusSeeOther)
//					return
//				} else {
//					http.Redirect(w, r, "/login.html", http.StatusSeeOther)
//					return
//				}
//			}
//			url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
//			keyset, err := jwk.Fetch(context.Background(), url)
//			src := []byte(c1.Value)
//			iss := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", os.Getenv("AWS_REGION"), os.Getenv("AWS_COGNITO_USER_POOL_ID"))
//			_, err = jwt.Parse(src,
//				jwt.WithKeySet(keyset),
//				jwt.WithValidate(true),
//				jwt.WithClaimValue("client_id", os.Getenv("AWS_COGNITO_APP_CLIENT_ID")), // replacing aud
//				jwt.WithIssuer(iss),
//				jwt.WithClaimValue("token_use", "access"),
//			)
//			if err != nil {
//				log.Println(err)
//				http.Redirect(w, r, "/login.html", http.StatusSeeOther)
//				return
//			}
//			// Call the next handler, which can be another middleware in the chain, or the final handler.
//			next.ServeHTTP(w, r)
//		})
//	}
//}

func ValidationMiddleWare(next http.Handler, url string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c1, err := r.Cookie("AccessToken")
		if err != nil {
			log.Println(err)
			if err.Error() == "http: named cookie not present" {
				http.Redirect(w, r, "/api/v1/refresh?url="+url, http.StatusSeeOther)
				return
			} else {
				http.Redirect(w, r, "/login.html", http.StatusSeeOther)
				return
			}
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
			http.Redirect(w, r, "/login.html", http.StatusSeeOther)
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
